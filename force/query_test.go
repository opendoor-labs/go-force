package force

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	"github.com/opendoor-labs/go-force/sobjects"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("query", func() {

		forceApi := createTest()
		desc, err := forceApi.DescribeSObject(&sobjects.Account{})
		if err != nil {
			GinkgoT().Fatalf("Failed to retrieve description of sobject: %v", err)
		}

		list := &AccountQueryResponse{}
		err = forceApi.Query(BuildQuery(desc.AllFields, desc.Name, nil), list)
		if err != nil {
			GinkgoT().Fatalf("Failed to query: %v", err)
		}
		GinkgoT().Logf("%#v", list)
	})
	It("query all", func() {

		forceApi := createTest()

		newId := insertSObject(forceApi, GinkgoT())
		deleteSObject(forceApi, GinkgoT(), newId)

		desc, err := forceApi.DescribeSObject(&sobjects.Account{})
		if err != nil {
			GinkgoT().Fatalf("Failed to retrieve description of sobject: %v", err)
		}

		list := &AccountQueryResponse{}
		err = forceApi.QueryAll(fmt.Sprintf(queryAll, desc.AllFields, newId), list)
		if err != nil {
			GinkgoT().Fatalf("Failed to queryAll: %v", err)
		}

		if len(list.Records) == 0 {
			GinkgoT().Fatal("Failed to retrieve deleted record using queryAll")
		}
		GinkgoT().Logf("%#v", list)
	})
	It("query next", func() {
	})
})

const (
	queryAll = "SELECT %v FROM Account WHERE Id = '%v'"
)

type AccountQueryResponse struct {
	sobjects.BaseQuery
	Records []sobjects.Account `json:"Records" force:"records"`
}
