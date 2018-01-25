package force

import (
	. "github.com/onsi/ginkgo"
	"net/http"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("has access", func() {

		oauth := &forceOauth{
			clientId:      testClientId,
			clientSecret:  testClientSecret,
			userName:      testUserName,
			password:      testPassword,
			securityToken: testSecurityToken,
			environment:   testEnvironment,
		}

		forceApi := &ForceApi{
			apiResources:           make(map[string]string),
			apiSObjects:            make(map[string]*SObjectMetaData),
			apiSObjectDescriptions: make(map[string]*SObjectDescription),
			apiVersion:             version,
			oauth:                  oauth,
			httpClient:             http.DefaultClient,
		}

		apiName := "AbadCompliance___c"

		forceApi.apiSObjects[apiName] = &SObjectMetaData{}
		validObjects := []string{apiName}
		if value := forceApi.HasAccess(validObjects); !value {
			GinkgoT().Error("expected to return true, but got false")
		}
		invalidObjects := []string{"Alien"}
		if value := forceApi.HasAccess(invalidObjects); value {
			GinkgoT().Error("expected to return false, but got true")
		}
	})
})
