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
				fmt.Fprintln(w, networks)
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
	})

})
