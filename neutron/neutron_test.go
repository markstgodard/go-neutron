package neutron_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/markstgodard/go-neutron/neutron"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const networks = `{
  "networks": [
    {
      "id": "e53a3b67-0074-404c-90b5-52ae217c3587",
      "name": "public",
      "description": "public network",
      "status": "ACTIVE",
      "subnets": [
        "3cc622c0-1ed8-470b-8d2b-081305df63b5",
        "9081fc4f-2415-4d99-ae99-0abb262ada90"
      ],
      "tenant_id": "1f77bad08b454898803a3d9f9e3799ec",
      "mtu": 1500,
      "project_id": "1f77bad08b454898803a3d9f9e3799ec"
    }
  ]
}`

const networksByName = `{
  "networks": [
    {
      "provider:physical_network": null,
      "ipv6_address_scope": null,
      "revision_number": 5,
      "port_security_enabled": true,
      "mtu": 1450,
      "id": "bd62af4c-bbe7-43fb-af21-29f3082fd734",
      "router:external": false,
      "availability_zone_hints": [],
      "availability_zones": [],
      "ipv4_address_scope": null,
      "shared": false,
      "project_id": "1f77bad08b454898803a3d9f9e3799ec",
      "status": "ACTIVE",
      "subnets": [
        "d087782e-3779-4982-b7ca-a4bde71b5aa5"
      ],
      "description": "",
      "tags": [],
      "updated_at": "2016-11-07T03:24:33Z",
      "provider:segmentation_id": 16,
      "name": "network1",
      "admin_state_up": true,
      "tenant_id": "1f77bad08b454898803a3d9f9e3799ec",
      "created_at": "2016-11-07T03:24:33Z",
      "provider:network_type": "vxlan"
    }
  ]
}`

const createNetworkResp = `{
  "network": {
    "provider:physical_network": null,
    "updated_at": "2016-11-07T05:58:51Z",
    "revision_number": 3,
    "port_security_enabled": true,
    "mtu": 1450,
    "id": "cc6c1929-6b26-4a1a-8680-3ea3dd09bfc6",
    "router:external": false,
    "availability_zone_hints": [],
    "availability_zones": [],
    "ipv4_address_scope": null,
    "shared": false,
    "project_id": "1f77bad08b454898803a3d9f9e3799ec",
    "status": "ACTIVE",
    "subnets": [],
    "description": "a sample network",
    "tags": [],
    "ipv6_address_scope": null,
    "provider:segmentation_id": 56,
    "name": "sample_network",
    "admin_state_up": true,
    "tenant_id": "1f77bad08b454898803a3d9f9e3799ec",
    "created_at": "2016-11-07T05:58:51Z",
    "provider:network_type": "vxlan"
  }
}
`
const createSubnetResp = `{
  "subnet": {
    "name": "",
    "network_id": "ed2e3c10-2e43-4297-9006-2863a2d1abbc",
    "tenant_id": "c1210485b2424d48804aad5d39c61b8f",
    "allocation_pools": [{"start": "10.0.3.20", "end": "10.0.3.150"}],
    "gateway_ip": "10.0.3.1",
    "ip_version": 4,
    "cidr": "10.0.3.0/24",
    "id": "9436e561-47bf-436a-b1f1-fe23a926e031",
    "enable_dhcp": true
  }
}`

const createPortResp = `{
    "port": {
        "admin_state_up": true,
        "device_id": "d6b4d3a5-c700-476f-b609-1493dd9dadc0",
        "device_owner": "",
        "fixed_ips": [
            {
                "ip_address": "192.168.111.4",
                "subnet_id": "22b44fc2-4ffb-4de4-b0f9-69d58b37ae27"
            }
        ],
        "id": "ebe69f1e-bc26-4db5-bed0-c0afb4afe3db",
        "mac_address": "fa:16:3e:a6:50:c1",
        "name": "port1",
        "network_id": "6aeaf34a-c482-4bd3-9dc3-7faf36412f12",
        "status": "ACTIVE",
        "tenant_id": "cf1a5775e766426cb1968766d0191908"
    }
}`

const networksEmpty = `{
  "networks": []
}`

const subnets = `{
  "subnets": [
    {
      "service_types": [],
      "description": "",
      "enable_dhcp": true,
      "network_id": "bd62af4c-bbe7-43fb-af21-29f3082fd734",
      "tenant_id": "1f77bad08b454898803a3d9f9e3799ec",
      "created_at": "2016-11-07T03:24:33Z",
      "dns_nameservers": [],
      "updated_at": "2016-11-07T03:24:33Z",
      "gateway_ip": "10.0.1.1",
      "ipv6_ra_mode": null,
      "allocation_pools": [
        {
          "start": "10.0.1.2",
          "end": "10.0.1.254"
        }
      ],
      "host_routes": [],
      "revision_number": 2,
      "ip_version": 4,
      "ipv6_address_mode": null,
      "cidr": "10.0.1.0/24",
      "project_id": "1f77bad08b454898803a3d9f9e3799ec",
      "id": "d087782e-3779-4982-b7ca-a4bde71b5aa5",
      "subnetpool_id": "759f65e4-9f27-4370-b329-1b3fb6ca529e",
      "name": "cf-subnet1"
    }
  ]
}`

var _ = Describe("Neutron API", func() {
	var (
		client *neutron.Client
		server *httptest.Server
	)

	Describe("NewClient", func() {
		var err error

		It("requires a URL and token", func() {
			client, err = neutron.NewClient("http://192.168.56.101:9696", "some-token")
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when URL is missing", func() {
			It("returns an error", func() {
				client, err = neutron.NewClient("", "some-token")
				Expect(err).To(MatchError("missing URL"))
			})
		})

		Context("when token is missing", func() {
			It("returns an error", func() {
				client, err = neutron.NewClient("http://192.168.56.101:9696", "")
				Expect(err).To(MatchError("missing token"))
			})
		})

	})

	Describe("Networks", func() {
		Describe("CreateNetwork", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(createNetworkResp))
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when network does not exist", func() {
				It("creates a new network", func() {
					net := neutron.Network{
						Name:         "sample_network",
						Description:  "a sample network",
						AdminStateUp: true,
					}
					network, err := client.CreateNetwork(net)
					Expect(err).ToNot(HaveOccurred())
					Expect(network.Name).To(Equal("sample_network"))
					Expect(network.ID).To(Equal("cc6c1929-6b26-4a1a-8680-3ea3dd09bfc6"))
					Expect(network.Description).To(Equal("a sample network"))
					Expect(network.Status).To(Equal("ACTIVE"))
					Expect(network.Subnets).To(HaveLen(0))
					Expect(network.MTU).To(Equal(1450))
					Expect(network.TenantID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
					Expect(network.ProjectID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				})
			})
		})

		Describe("Networks", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, networks)
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			It("lists networks", func() {
				networks, err := client.Networks()
				Expect(err).ToNot(HaveOccurred())
				Expect(networks).To(HaveLen(1))
				Expect(networks[0].Name).To(Equal("public"))
				Expect(networks[0].ID).To(Equal("e53a3b67-0074-404c-90b5-52ae217c3587"))
				Expect(networks[0].Description).To(Equal("public network"))
				Expect(networks[0].Status).To(Equal("ACTIVE"))
				Expect(networks[0].Subnets).To(HaveLen(2))
				Expect(networks[0].TenantID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				Expect(networks[0].MTU).To(Equal(1500))
				Expect(networks[0].ProjectID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
			})
		})

		Describe("NetworksByName", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Query().Get("name") {
					case "network1":
						fmt.Fprintln(w, networksByName)
					default:
						fmt.Fprintln(w, networksEmpty)
					}
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			It("list networks by name", func() {
				networks, err := client.NetworksByName("network1")
				Expect(err).ToNot(HaveOccurred())
				Expect(networks).To(HaveLen(1))
				Expect(networks[0].Name).To(Equal("network1"))
				Expect(networks[0].ID).To(Equal("bd62af4c-bbe7-43fb-af21-29f3082fd734"))
				Expect(networks[0].Description).To(Equal(""))
				Expect(networks[0].Status).To(Equal("ACTIVE"))
				Expect(networks[0].Subnets).To(HaveLen(1))
				Expect(networks[0].TenantID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				Expect(networks[0].MTU).To(Equal(1450))
				Expect(networks[0].ProjectID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
			})

			Context("when network name is invalid", func() {
				It("returns an error", func() {
					_, err := client.NetworksByName("")
					Expect(err).To(MatchError("empty 'name' parameter"))
				})
			})

			Context("when network does not exist", func() {
				It("returns empty when not found by name", func() {
					networks, err := client.NetworksByName("does-not-exist")
					Expect(err).ToNot(HaveOccurred())
					Expect(networks).To(HaveLen(0))
				})
			})
		})

		Describe("DeleteNetwork", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when network exists", func() {
				It("deletes a network", func() {
					err := client.DeleteNetwork("network-id-1")
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

	})

	Describe("Subnets", func() {
		Describe("CreateSubnet", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(createSubnetResp))
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when subnet does not exist", func() {
				It("creates a new subnet", func() {
					subnet := neutron.Subnet{
						NetworkID: "network1",
						IPVersion: 4,
						CIDR:      "10.0.3.0/24",
						AllocationPools: []neutron.AllocationPool{
							{
								Start: "10.0.3.20",
								End:   "10.0.3.150",
							},
						},
					}

					sn, err := client.CreateSubnet(subnet)
					Expect(err).ToNot(HaveOccurred())
					Expect(sn.ID).To(Equal("9436e561-47bf-436a-b1f1-fe23a926e031"))
					Expect(sn.NetworkID).To(Equal("ed2e3c10-2e43-4297-9006-2863a2d1abbc"))
					Expect(sn.CIDR).To(Equal("10.0.3.0/24"))
					Expect(sn.IPVersion).To(Equal(4))
				})
			})
		})

		Describe("Subnets", func() {
			BeforeEach(func() {
				var err error
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, subnets)
				}))

				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			It("lists all subnets owned by a project", func() {
				subnets, err := client.Subnets()
				Expect(err).ToNot(HaveOccurred())
				Expect(subnets).To(HaveLen(1))
				Expect(subnets[0].Name).To(Equal("cf-subnet1"))
				Expect(subnets[0].ID).To(Equal("d087782e-3779-4982-b7ca-a4bde71b5aa5"))
				Expect(subnets[0].SubnetPoolID).To(Equal("759f65e4-9f27-4370-b329-1b3fb6ca529e"))
				Expect(subnets[0].ProjectID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				Expect(subnets[0].TenantID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				Expect(subnets[0].NetworkID).To(Equal("bd62af4c-bbe7-43fb-af21-29f3082fd734"))
				Expect(subnets[0].IPVersion).To(Equal(4))
				Expect(subnets[0].CIDR).To(Equal("10.0.1.0/24"))
				Expect(subnets[0].AllocationPools).To(Equal(
					[]neutron.AllocationPool{{Start: "10.0.1.2", End: "10.0.1.254"}},
				))
			})

			It("lists subnets owned by a project by name", func() {
				subnets, err := client.SubnetsByName("cf-subnet1")
				Expect(err).ToNot(HaveOccurred())
				Expect(subnets).To(HaveLen(1))
				Expect(subnets[0].Name).To(Equal("cf-subnet1"))
				Expect(subnets[0].ID).To(Equal("d087782e-3779-4982-b7ca-a4bde71b5aa5"))
				Expect(subnets[0].SubnetPoolID).To(Equal("759f65e4-9f27-4370-b329-1b3fb6ca529e"))
				Expect(subnets[0].ProjectID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				Expect(subnets[0].TenantID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
				Expect(subnets[0].NetworkID).To(Equal("bd62af4c-bbe7-43fb-af21-29f3082fd734"))
				Expect(subnets[0].IPVersion).To(Equal(4))
				Expect(subnets[0].CIDR).To(Equal("10.0.1.0/24"))
				Expect(subnets[0].AllocationPools).To(Equal(
					[]neutron.AllocationPool{{Start: "10.0.1.2", End: "10.0.1.254"}},
				))
			})

			Context("when subnet name is invalid", func() {
				It("returns an error", func() {
					_, err := client.SubnetsByName("")
					Expect(err).To(MatchError("empty 'name' parameter"))
				})
			})
		})
	})

	Describe("Ports", func() {
		Describe("CreatePort", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(createPortResp))
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when port does not exist", func() {
				It("creates a new port", func() {
					port := neutron.Port{
						NetworkID:    "6aeaf34a-c482-4bd3-9dc3-7faf36412f12",
						Name:         "port1",
						DeviceID:     "d6b4d3a5-c700-476f-b609-1493dd9dadc0",
						AdminStateUp: true,
					}

					p, err := client.CreatePort(port)
					Expect(err).ToNot(HaveOccurred())
					Expect(p.ID).To(Equal("ebe69f1e-bc26-4db5-bed0-c0afb4afe3db"))
					Expect(p.Name).To(Equal("port1"))
					Expect(p.NetworkID).To(Equal("6aeaf34a-c482-4bd3-9dc3-7faf36412f12"))
					Expect(p.TenantID).To(Equal("cf1a5775e766426cb1968766d0191908"))
					Expect(p.Status).To(Equal("ACTIVE"))
					Expect(p.AdminStateUp).To(BeTrue())
					Expect(p.MacAddress).To(Equal("fa:16:3e:a6:50:c1"))
					Expect(p.DeviceOwner).To(Equal(""))
					Expect(p.DeviceID).To(Equal("d6b4d3a5-c700-476f-b609-1493dd9dadc0"))
					Expect(p.FixedIPs).To(Equal(
						[]neutron.FixedIP{
							{
								IPAddress: "192.168.111.4",
								SubnetID:  "22b44fc2-4ffb-4de4-b0f9-69d58b37ae27",
							},
						},
					))
				})
			})
		})

		Describe("DeletePort", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}))
				var err error
				client, err = neutron.NewClient(server.URL, "some-token")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when port exists", func() {
				It("deletes the port", func() {
					err := client.DeletePort("port1")
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
