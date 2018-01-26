package force_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opendoor-labs/go-force/force"
	"github.com/opendoor-labs/go-force/force/forcefakes"
)

const FakeOauthRespBody = `{"access_token": "at", "instance_url": "iu", "id": "anid", "issued_at": "ia", "signature": "sig"}`

func NewFakeResponse(body string, statusCode int) (resp *http.Response) {
	return &http.Response{
		StatusCode: statusCode,
		Proto:      "HTTP/1.0",
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}

func createForceApi(httpClient *forcefakes.FakeHttpClient) (*force.ForceApi, error) {
	// The first 3 http Do() calls are made by Create to setup Auth +
	// resources.  Users of the the returned ForceApi instance with the
	// mocked httpClient that's returned should mock return calls starting
	// at 3.

	authResp := NewFakeResponse(FakeOauthRespBody, 200)
	httpClient.DoReturnsOnCall(0, authResp, nil)

	apiResourceResp := NewFakeResponse(`{"sobjects": "sobjects-resources"}`, 200)
	httpClient.DoReturnsOnCall(1, apiResourceResp, nil)

	apiSObjectsResp := NewFakeResponse(`{"sobjects": [{"name": "APIName", "urls": {"sobject": "the/url"}}]}`, 200)
	httpClient.DoReturnsOnCall(2, apiSObjectsResp, nil)

	return force.Create("", "", "", "", "", "", "", httpClient)
}

func TestForce(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Force Suite")
}
