package force_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opendoor-labs/go-force/force"
	"github.com/opendoor-labs/go-force/force/forcefakes"
	"github.com/opendoor-labs/go-force/sobjects"
)

const oauthRespBody = `{"access_token": "at", "instance_url": "iu", "id": "anid", "issued_at": "ia", "signature": "sig"}`

func newResponse(body string) (resp *http.Response) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}

func createForceApi(httpClient *forcefakes.FakeHttpClient) (*force.ForceApi, error) {
	authResp := newResponse(oauthRespBody)
	httpClient.DoReturnsOnCall(0, authResp, nil)

	apiResourceResp := newResponse(`{"aresourceKey": "aResourceValue"}`)
	httpClient.DoReturnsOnCall(1, apiResourceResp, nil)

	apiSObjectsResp := newResponse(`{"ASObject": {"name": "TheName"}}`)
	httpClient.DoReturnsOnCall(2, apiSObjectsResp, nil)

	return force.Create("", "", "", "", "", "", "", httpClient)
}

var _ = Describe("Client", func() {
	Describe("Get", func() {
		It("should unmarshal a resource", func() {
			httpClient := forcefakes.FakeHttpClient{}

			forceApi, err := createForceApi(&httpClient)
			Expect(err).NotTo(HaveOccurred())

			apiSObjectsResp := newResponse(`{"Id": "SFID-123"}`)
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
