package force_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opendoor-labs/go-force/force/forcefakes"
	"github.com/opendoor-labs/go-force/sobjects"
)

const oauthRespBody = `{"access_token": "at", "instance_url": "iu", "id": "anid", "issued_at": "ia", "signature": "sig"}`

var _ = Describe("Client", func() {
	Describe("Get", func() {
		It("should unmarshal a resource", func() {
			httpClient := forcefakes.FakeHttpClient{}

			forceApi, err := createForceApi(&httpClient)
			Expect(err).NotTo(HaveOccurred())

			apiSObjectsResp := NewFakeResponse(`{"Id": "SFID-123"}`, 200)
			httpClient.DoReturnsOnCall(3, apiSObjectsResp, nil)

			params := url.Values{"fields": []string{"Id"}}
			sobj := sobjects.BaseSObject{}

			status, err := forceApi.Get("/sobjects/path", params, &sobj)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(200))
			Expect(sobj.Id).To(Equal("SFID-123"))
		})
	})
})
