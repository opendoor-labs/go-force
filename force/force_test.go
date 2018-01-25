package force

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	"github.com/opendoor-labs/go-force/sobjects"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("create with access token", func() {

		oauth := &forceOauth{
			clientId:      testClientId,
			clientSecret:  testClientSecret,
			userName:      testUserName,
			password:      testPassword,
			securityToken: testSecurityToken,
			environment:   testEnvironment,
			httpClient:    http.DefaultClient,
		}

		forceApi := &ForceApi{
			apiResources:           make(map[string]string),
			apiSObjects:            make(map[string]*SObjectMetaData),
			apiSObjectDescriptions: make(map[string]*SObjectDescription),
			apiVersion:             version,
			oauth:                  oauth,
			httpClient:             http.DefaultClient,
		}

		err := forceApi.oauth.Authenticate()
		if err != nil {
			GinkgoT().Fatalf("Unable to authenticate: %#v", err)
		}
		if err := forceApi.oauth.Validate(); err != nil {
			GinkgoT().Fatalf("Oauth object is invlaid: %#v", err)
		}

		newForceApi, err := CreateWithAccessToken(testVersion, testClientId, forceApi.oauth.AccessToken, forceApi.oauth.InstanceUrl, http.DefaultClient)
		if err != nil {
			GinkgoT().Fatalf("Unable to create new force api instance using pre-defined oauth details: %#v", err)
		}
		if err := newForceApi.oauth.Validate(); err != nil {
			GinkgoT().Fatalf("Oauth object is invlaid: %#v", err)
		}

		_, err = newForceApi.DescribeSObject(&sobjects.Account{})
		if err != nil {
			GinkgoT().Fatalf("Failed to retrieve description of sobject: %v", err)
		}
	})
	It("passed http client is used", func() {

		oauth := &forceOauth{
			clientId:      testClientId,
			clientSecret:  testClientSecret,
			userName:      testUserName,
			password:      testPassword,
			securityToken: testSecurityToken,
			environment:   testEnvironment,
			httpClient:    http.DefaultClient,
		}

		forceApi := &ForceApi{
			apiResources:           make(map[string]string),
			apiSObjects:            make(map[string]*SObjectMetaData),
			apiSObjectDescriptions: make(map[string]*SObjectDescription),
			apiVersion:             version,
			oauth:                  oauth,
			httpClient:             http.DefaultClient,
		}

		err := forceApi.oauth.Authenticate()
		if err != nil {
			GinkgoT().Fatalf("Unable to authenticate: %#v", err)
		}
		if err := forceApi.oauth.Validate(); err != nil {
			GinkgoT().Fatalf("Oauth object is invlaid: %#v", err)
		}

		trt := TestRoundTripper{CallCount: 0, proxied: http.DefaultTransport}
		httpClient := &http.Client{Transport: &trt}
		newForceApi, _ := CreateWithAccessToken(testVersion, testClientId, forceApi.oauth.AccessToken, forceApi.oauth.InstanceUrl, httpClient)

		newForceApi.DescribeSObject(&sobjects.Account{})

		if trt.CallCount < 1 {
			GinkgoT().Fatalf("Passed in http client was not called")
		}
	})
})

type TestRoundTripper struct {
	CallCount int
	proxied   http.RoundTripper
}

func (trt *TestRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	trt.CallCount++
	return trt.proxied.RoundTrip(req)
}
