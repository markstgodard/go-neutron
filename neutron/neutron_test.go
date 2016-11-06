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
      "description": "",
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

	Describe("Networks", func() {

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, networks)
			}))

			client = neutron.NewClient(server.URL, "some-token")
		})

		AfterEach(func() {
			server.Close()
		})

		It("can list networks", func() {
			resp, err := client.Networks()
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Networks).To(HaveLen(1))
		})

	})

})
