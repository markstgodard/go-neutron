package neutron_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/markstgodard/go-neutron/neutron"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Neutron API", func() {
	var (
		client *neutron.Client
		server *httptest.Server
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

		BeforeEach(func() {
			var err error
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("name") != "" {
					fmt.Fprintln(w, networksByName)
				} else {
					fmt.Fprintln(w, networks)
				}
			}))

			client, err = neutron.NewClient(server.URL, "some-token")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			server.Close()
		})

		It("can list networks", func() {
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

		It("can list networks by name", func() {
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
	})

	// Describe("Subnets", func() {
	// 	BeforeEach(func() {
	// 		var err error
	// 		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			fmt.Fprintln(w, subnets)
	// 		}))

	// 		client, err = neutron.NewClient(server.URL, "some-token")
	// 		Expect(err).ToNot(HaveOccurred())
	// 	})

	// 	AfterEach(func() {
	// 		server.Close()
	// 	})

	// 	It("can list subnets", func() {
	// 		subnets, err := client.Subnets()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(subnets).To(HaveLen(1))
	// 		Expect(subnets[0].Name).To(Equal("public"))
	// 		Expect(subnets[0].ID).To(Equal("e53a3b67-0074-404c-90b5-52ae217c3587"))
	// 		Expect(subnets[0].Description).To(Equal("public network"))
	// 		Expect(networks[0].ProjectID).To(Equal("1f77bad08b454898803a3d9f9e3799ec"))
	// 	})
	// })

})
