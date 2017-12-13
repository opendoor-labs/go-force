package force

import (
	"net/http"
	"testing"

	"github.com/opendoor-labs/go-force/sobjects"
)

func TestCreateWithAccessToken(t *testing.T) {

	// Manually grab an OAuth token, so that we can pass it into CreateWithAccessToken
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
		t.Fatalf("Unable to authenticate: %#v", err)
	}
	if err := forceApi.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// We shouldn't hit any errors creating a new force instance and manually passing in these oauth details now.
	newForceApi, err := CreateWithAccessToken(testVersion, testClientId, forceApi.oauth.AccessToken, forceApi.oauth.InstanceUrl, http.DefaultClient)
	if err != nil {
		t.Fatalf("Unable to create new force api instance using pre-defined oauth details: %#v", err)
	}
	if err := newForceApi.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// We should be able to make a basic query now with the newly created object (i.e. the oauth details should be correctly usable).
	_, err = newForceApi.DescribeSObject(&sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}
}

type TestRoundTripper struct {
	CallCount int
	proxied   http.RoundTripper
}

func (trt *TestRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	trt.CallCount++
	return trt.proxied.RoundTrip(req)
}

func TestPassedHttpClientIsUsed(t *testing.T) {
	// Manually grab an OAuth token, so that we can pass it into CreateWithAccessToken
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
		t.Fatalf("Unable to authenticate: %#v", err)
	}
	if err := forceApi.oauth.Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}

	// need to proxy the real client b/c that's how goforce tests roll
	trt := TestRoundTripper{CallCount: 0, proxied: http.DefaultTransport}
	httpClient := &http.Client{Transport: &trt}
	newForceApi, _ := CreateWithAccessToken(testVersion, testClientId, forceApi.oauth.AccessToken, forceApi.oauth.InstanceUrl, httpClient)

	newForceApi.DescribeSObject(&sobjects.Account{})

	if trt.CallCount < 1 {
		t.Fatalf("Passed in http client was not called")
	}
}
