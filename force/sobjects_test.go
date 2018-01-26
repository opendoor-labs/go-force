package force

import (
	"math/rand"

	"time"

	. "github.com/onsi/ginkgo"
	"github.com/opendoor-labs/go-force/sobjects"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("describe sobjects", func() {

		forceAPI := createTest()
		objects, err := forceAPI.DescribeSObjects()
		if err != nil {
			GinkgoT().Fatal("Failed to retrieve SObjects", err)
		}
		GinkgoT().Logf("SObjects for Account Retrieved: %+v", objects)
	})
	It("describe s object", func() {

		forceApi := createTest()
		acc := &sobjects.Account{}

		desc, err := forceApi.DescribeSObject(acc)
		if err != nil {
			GinkgoT().Fatalf("Cannot retrieve SObject Description for Account SObject: %v", err)
		}
		GinkgoT().Logf("SObject Description for Account Retrieved: %+v", desc)
	})
	It("get s object", func() {

		forceApi := createTest()

		acc := &sobjects.Account{}

		err := forceApi.GetSObject(AccountId, nil, acc)
		if err != nil {
			GinkgoT().Fatalf("Cannot retrieve SObject Account: %v", err)
		}
		GinkgoT().Logf("SObject Account Retrieved: %+v", acc)

		customObject := &CustomSObject{}

		err = forceApi.GetSObject(CustomObjectId, nil, customObject)
		if err != nil {
			GinkgoT().Fatalf("Cannot retrieve SObject CustomObject: %v", err)
		}
		GinkgoT().Logf("SObject CustomObject Retrieved: %+v", customObject)

		fields := []string{"Name", "Id"}

		accFields := &sobjects.Account{}

		err = forceApi.GetSObject(AccountId, fields, accFields)
		if err != nil {
			GinkgoT().Fatalf("Cannot retrieve SObject Account fields: %v", err)
		}
		GinkgoT().Logf("SObject Account Name and Id Retrieved: %+v", accFields)
	})
	It("update s object", func() {

		forceApi := createTest()

		rand.Seed(time.Now().UTC().UnixNano())
		someText := randomString(10)

		acc := &sobjects.Account{}
		acc.Name = someText

		err := forceApi.UpdateSObject(AccountId, acc)
		if err != nil {
			GinkgoT().Fatalf("Cannot update SObject Account: %v", err)
		}

		err = forceApi.GetSObject(AccountId, nil, acc)
		if err != nil {
			GinkgoT().Fatalf("Cannot retrieve SObject Account: %v", err)
		}

		if acc.Name != someText {
			GinkgoT().Fatalf("Update SObject Account failed. Failed to persist.")
		}
		GinkgoT().Logf("Updated SObject Account: %+v", acc)
	})
	It("insert delete s object", func() {

		forceApi := createTest()
		objectId := insertSObject(forceApi, GinkgoT())
		deleteSObject(forceApi, GinkgoT(), objectId)
	})
})

const (
	AccountId      = "001i000000RxW18"
	CustomObjectId = "a00i0000009SPer"
)

type CustomSObject struct {
	sobjects.BaseSObject
	Active    bool   `force:"Active__c"`
	AccountId string `force:"Account__c"`
}

func (t *CustomSObject) APIName() string {
	return "CustomObject__c"
}

func insertSObject(forceApi *ForceApi, t GinkgoTInterface) string {

	rand.Seed(time.Now().UTC().UnixNano())
	someText := randomString(10)

	acc := &sobjects.Account{}
	acc.Name = someText

	resp, err := forceApi.InsertSObject(acc)
	if err != nil {
		t.Fatalf("Insert SObject Account failed: %v", err)
	}

	if len(resp.Id) == 0 {
		t.Fatalf("Insert SObject Account failed to return Id: %+v", resp)
	}

	return resp.Id
}

func deleteSObject(forceApi *ForceApi, t GinkgoTInterface, id string) {

	acc := &sobjects.Account{}

	err := forceApi.DeleteSObject(id, acc)
	if err != nil {
		t.Fatalf("Delete SObject Account failed: %v", err)
	}

	err = forceApi.GetSObject(id, nil, acc)
	if err == nil {
		t.Fatalf("Delete SObject Account failed, was able to retrieve deleted object: %+v", acc)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
