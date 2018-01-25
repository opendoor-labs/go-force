package forcejson

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("struct tag object key", func() {

		for _, tt := range structTagObjectKeyTests {
			b, err := Marshal(tt.raw)
			if err != nil {
				GinkgoT().Fatalf("Marshal(%#q) failed: %v", tt.raw, err)
			}
			var f interface{}
			err = Unmarshal(b, &f)
			if err != nil {
				GinkgoT().Fatalf("Unmarshal(%#q) failed: %v", b, err)
			}
			for i, v := range f.(map[string]interface{}) {
				switch i {
				case tt.key:
					if s, ok := v.(string); !ok || s != tt.value {
						GinkgoT().Fatalf("Unexpected value: %#q, want %v", s, tt.value)
					}
				default:
					GinkgoT().Fatalf("Unexpected key: %#q, from %#q", i, b)
				}
			}
		}
	})
})

type basicLatin2xTag struct {
	V string `force:"$%-/"`
}

type basicLatin3xTag struct {
	V string `force:"0123456789"`
}

type basicLatin4xTag struct {
	V string `force:"ABCDEFGHIJKLMO"`
}

type basicLatin5xTag struct {
	V string `force:"PQRSTUVWXYZ_"`
}

type basicLatin6xTag struct {
	V string `force:"abcdefghijklmno"`
}

type basicLatin7xTag struct {
	V string `force:"pqrstuvwxyz"`
}

type miscPlaneTag struct {
	V string `force:"色は匂へど"`
}

type percentSlashTag struct {
	V string `force:"text/html%"`
}

type punctuationTag struct {
	V string `force:"!#$%&()*+-./:<=>?@[]^_{|}~"`
}

type emptyTag struct {
	W string
}

type misnamedTag struct {
	X string `jsom:"Misnamed"`
}

type badFormatTag struct {
	Y string `:"BadFormat"`
}

type badCodeTag struct {
	Z string `force:" !\"#&'()*+,."`
}

type spaceTag struct {
	Q string `force:"With space"`
}

type unicodeTag struct {
	W string `force:"Ελλάδα"`
}

var structTagObjectKeyTests = []struct {
	raw   interface{}
	value string
	key   string
}{
	{basicLatin2xTag{"2x"}, "2x", "$%-/"},
	{basicLatin3xTag{"3x"}, "3x", "0123456789"},
	{basicLatin4xTag{"4x"}, "4x", "ABCDEFGHIJKLMO"},
	{basicLatin5xTag{"5x"}, "5x", "PQRSTUVWXYZ_"},
	{basicLatin6xTag{"6x"}, "6x", "abcdefghijklmno"},
	{basicLatin7xTag{"7x"}, "7x", "pqrstuvwxyz"},
	{miscPlaneTag{"いろはにほへと"}, "いろはにほへと", "色は匂へど"},
	{emptyTag{"Pour Moi"}, "Pour Moi", "W"},
	{misnamedTag{"Animal Kingdom"}, "Animal Kingdom", "X"},
	{badFormatTag{"Orfevre"}, "Orfevre", "Y"},
	{badCodeTag{"Reliable Man"}, "Reliable Man", "Z"},
	{percentSlashTag{"brut"}, "brut", "text/html%"},
	{punctuationTag{"Union Rags"}, "Union Rags", "!#$%&()*+-./:<=>?@[]^_{|}~"},
	{spaceTag{"Perreddu"}, "Perreddu", "With space"},
	{unicodeTag{"Loukanikos"}, "Loukanikos", "Ελλάδα"},
}
