package force_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opendoor-labs/go-force/force"
	"github.com/opendoor-labs/go-force/force/forcefakes"
)

const sobjId_1 = "00Q29000004RbVPEA0"
const sobjId_2 = "00Q29000004RbVUEA0"
const baseLeadURI = "/services/data/v40.0/sobjects/Lead"

var multiURISObjectResponse = fmt.Sprintf(`["%s/%s","%s/%s"]`, baseLeadURI, sobjId_1, baseLeadURI, sobjId_2)

var _ = Describe("Sobjects", func() {
	var httpClient forcefakes.FakeHttpClient
	var forceApi *force.ForceApi

	BeforeEach(func() {
		httpClient = forcefakes.FakeHttpClient{}

		var err error
		forceApi, err = createForceApi(&httpClient)
		Expect(err).NotTo(HaveOccurred())
		Expect(forceApi).NotTo(BeNil())
	})

	Describe("GetSFIDsByExternalId", func() {
		Context("a single object is returned", func() {
			It("should return a slice with one ID", func() {
				apiSObjectsResp := NewFakeResponse(`{"Id": "SFID-123"}`, 200)
				httpClient.DoReturnsOnCall(3, apiSObjectsResp, nil)

				actualIds, status, err := forceApi.GetSFIDsByExternalId("APIName", "ExKey", "ExId-123")
				Expect(err).NotTo(HaveOccurred())
				Expect(actualIds).To(Equal([]string{"SFID-123"}))
				Expect(status).To(Equal(200))
			})
		})

		Context("multiple objects are returned", func() {
			It("should return a slice with the IDs", func() {
				// When there are multiple results, the same call is made
				// twice, the first call attempts unmarshalling into the SObject,
				// when a 300 is detected, the call is made again passing into it
				// a slice for unmarshalling multiple URIs.
				apiSObjectsResp := NewFakeResponse(multiURISObjectResponse, 300)
				httpClient.DoReturnsOnCall(3, apiSObjectsResp, nil)
				apiSObjectsResp = NewFakeResponse(multiURISObjectResponse, 300)
				httpClient.DoReturnsOnCall(4, apiSObjectsResp, nil)

				actualIds, status, err := forceApi.GetSFIDsByExternalId("APIName", "ExKey", "ExId-123")
				Expect(err).NotTo(HaveOccurred())
				Expect(actualIds).To(Equal([]string{sobjId_1, sobjId_2}))
				Expect(status).To(Equal(300))
			})
		})

		Context("request failed", func() {
			It("should return the statusCode and the error", func() {
				apiSObjectsResp := NewFakeResponse(multiURISObjectResponse, 400)
				httpClient.DoReturnsOnCall(3, apiSObjectsResp, nil)
				apiSObjectsResp = NewFakeResponse(multiURISObjectResponse, 400)
				httpClient.DoReturnsOnCall(4, apiSObjectsResp, nil)

				actualIds, status, err := forceApi.GetSFIDsByExternalId("APIName", "ExKey", "ExId-123")
				Expect(err).To(HaveOccurred())
				Expect(status).To(Equal(400))
				Expect(actualIds).To(BeNil())
			})
		})
	})
})
