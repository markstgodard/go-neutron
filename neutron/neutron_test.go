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
