package parser_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opendoor-labs/go-force/force/parser"
)

type Key struct {
	ID  *int    `json:"id"`
	Key *string `json:"key"`
}

type Keys []Key

var _ = Describe("Parser", func() {
	var (
		jsonArray = []byte(`[
      {"id": 1, "key": "secret-1"},
      {"id": 2, "key": "secret-2"},
      {"id": 3, "key": "secret-3"}
    ]`)

		jsonObj     = []byte(`{"id": 1, "key": "secret-1"}`)
		jsonInvalid = []byte(`☣`)

		slice []Key
		key   Key
		keys  Keys
	)

	BeforeEach(func() {
		slice = []Key{}
		keys = Keys{}
	})

	Describe("IsSlicePtr", func() {
		It("should be true for a pointer to a slice type", func() {
			Expect(parser.IsSlicePtr(&keys)).To(BeTrue())
		})

		It("should be true for a pointer to a slice", func() {
			Expect(parser.IsSlicePtr(&slice)).To(BeTrue())
		})

		It("should be false otherwise", func() {
			Expect(parser.IsSlicePtr(slice)).To(BeFalse())
			Expect(parser.IsSlicePtr("")).To(BeFalse())
			Expect(parser.IsSlicePtr(234)).To(BeFalse())
			Expect(parser.IsSlicePtr(nil)).To(BeFalse())
		})
	})

	Describe("ParseSFJSON", func() {
		Context("with a pointer to a slice", func() {
			It("should parse single object into the slice of length 1", func() {
				slice := make([]Key, 0)
				err := parser.ParseSFJSON(jsonObj, &slice)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(slice)).To(Equal(1))

				obj := slice[0]
				Expect(*obj.ID).To(Equal(1))
				Expect(*obj.Key).To(Equal("secret-1"))
			})

			It("should parse an array into the slice", func() {
				err := parser.ParseSFJSON(jsonArray, &slice)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(slice)).To(Equal(3))

				for i := 0; i < 3; i++ {
					obj := slice[i]
					Expect(*obj.ID).To(Equal(i + 1))
					Expect(*obj.Key).To(Equal(fmt.Sprintf("secret-%d", i+1)))
				}
			})

			It("should return an error if invalid JSON", func() {
				err := parser.ParseSFJSON(jsonInvalid, &slice)
				Expect(err).To(HaveOccurred())
				Expect(len(slice)).To(Equal(0))
			})
		})

		Context("with a pointer to a struct", func() {
			It("should parse an object", func() {
				err := parser.ParseSFJSON(jsonObj, &key)
				Expect(err).ToNot(HaveOccurred())
				Expect(*key.ID).To(Equal(1))
				Expect(*key.Key).To(Equal("secret-1"))
			})

			It("should return an error if json is an array", func() {
				err := parser.ParseSFJSON(jsonArray, &key)
				Expect(err).To(HaveOccurred())
			})

			It("should return an error if invalid JSON", func() {
				err := parser.ParseSFJSON(jsonInvalid, &key)
				Expect(err).To(HaveOccurred())
				Expect(len(slice)).To(Equal(0))
			})
		})
	})
})
