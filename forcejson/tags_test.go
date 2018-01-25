package forcejson

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("tag parsing", func() {

		name, opts := parseTag("field,foobar,foo")
		if name != "field" {
			GinkgoT().Fatalf("name = %q, want field", name)
		}
		for _, tt := range []struct {
			opt  string
			want bool
		}{
			{"foobar", true},
			{"foo", true},
			{"bar", false},
		} {
			if opts.Contains(tt.opt) != tt.want {
				GinkgoT().Errorf("Contains(%q) = %v", tt.opt, !tt.want)
			}
		}
	})
})
