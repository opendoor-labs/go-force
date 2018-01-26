package forcejson

import (
	"bytes"
	"math"
	"math/rand"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("compact", func() {

		var buf bytes.Buffer
		for _, tt := range examples {
			buf.Reset()
			if err := Compact(&buf, []byte(tt.compact)); err != nil {
				GinkgoT().Errorf("Compact(%#q): %v", tt.compact, err)
			} else if s := buf.String(); s != tt.compact {
				GinkgoT().Errorf("Compact(%#q) = %#q, want original", tt.compact, s)
			}

			buf.Reset()
			if err := Compact(&buf, []byte(tt.indent)); err != nil {
				GinkgoT().Errorf("Compact(%#q): %v", tt.indent, err)
				continue
			} else if s := buf.String(); s != tt.compact {
				GinkgoT().Errorf("Compact(%#q) = %#q, want %#q", tt.indent, s, tt.compact)
			}
		}
	})
	It("compact separators", func() {

		tests := []struct {
			in, compact string
		}{
			{"{\"\u2028\": 1}", `{"\u2028":1}`},
			{"{\"\u2029\" :2}", `{"\u2029":2}`},
		}
		for _, tt := range tests {
			var buf bytes.Buffer
			if err := Compact(&buf, []byte(tt.in)); err != nil {
				GinkgoT().Errorf("Compact(%q): %v", tt.in, err)
			} else if s := buf.String(); s != tt.compact {
				GinkgoT().Errorf("Compact(%q) = %q, want %q", tt.in, s, tt.compact)
			}
		}
	})
	It("indent", func() {

		var buf bytes.Buffer
		for _, tt := range examples {
			buf.Reset()
			if err := Indent(&buf, []byte(tt.indent), "", "\t"); err != nil {
				GinkgoT().Errorf("Indent(%#q): %v", tt.indent, err)
			} else if s := buf.String(); s != tt.indent {
				GinkgoT().Errorf("Indent(%#q) = %#q, want original", tt.indent, s)
			}

			buf.Reset()
			if err := Indent(&buf, []byte(tt.compact), "", "\t"); err != nil {
				GinkgoT().Errorf("Indent(%#q): %v", tt.compact, err)
				continue
			} else if s := buf.String(); s != tt.indent {
				GinkgoT().Errorf("Indent(%#q) = %#q, want %#q", tt.compact, s, tt.indent)
			}
		}
	})
	It("compact big", func() {

		initBig()
		var buf bytes.Buffer
		if err := Compact(&buf, jsonBig); err != nil {
			GinkgoT().Fatalf("Compact: %v", err)
		}
		b := buf.Bytes()
		if !bytes.Equal(b, jsonBig) {
			GinkgoT().Error("Compact(jsonBig) != jsonBig")
			diff(GinkgoT(), b, jsonBig)
			return
		}
	})
	It("indent big", func() {

		initBig()
		var buf bytes.Buffer
		if err := Indent(&buf, jsonBig, "", "\t"); err != nil {
			GinkgoT().Fatalf("Indent1: %v", err)
		}
		b := buf.Bytes()
		if len(b) == len(jsonBig) {
			GinkgoT().Fatalf("Indent(jsonBig) did not get bigger")
		}

		var buf1 bytes.Buffer
		if err := Indent(&buf1, b, "", "\t"); err != nil {
			GinkgoT().Fatalf("Indent2: %v", err)
		}
		b1 := buf1.Bytes()
		if !bytes.Equal(b1, b) {
			GinkgoT().Error("Indent(Indent(jsonBig)) != Indent(jsonBig)")
			diff(GinkgoT(), b1, b)
			return
		}

		buf1.Reset()
		if err := Compact(&buf1, b); err != nil {
			GinkgoT().Fatalf("Compact: %v", err)
		}
		b1 = buf1.Bytes()
		if !bytes.Equal(b1, jsonBig) {
			GinkgoT().Error("Compact(Indent(jsonBig)) != jsonBig")
			diff(GinkgoT(), b1, jsonBig)
			return
		}
	})
	It("indent errors", func() {

		for i, tt := range indentErrorTests {
			slice := make([]uint8, 0)
			buf := bytes.NewBuffer(slice)
			if err := Indent(buf, []uint8(tt.in), "", ""); err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					GinkgoT().Errorf("#%d: Indent: %#v", i, err)
					continue
				}
			}
		}
	})
	It("next value big", func() {

		initBig()
		var scan scanner
		item, rest, err := nextValue(jsonBig, &scan)
		if err != nil {
			GinkgoT().Fatalf("nextValue: %s", err)
		}
		if len(item) != len(jsonBig) || &item[0] != &jsonBig[0] {
			GinkgoT().Errorf("invalid item: %d %d", len(item), len(jsonBig))
		}
		if len(rest) != 0 {
			GinkgoT().Errorf("invalid rest: %d", len(rest))
		}

		item, rest, err = nextValue(append(jsonBig, "HELLO WORLD"...), &scan)
		if err != nil {
			GinkgoT().Fatalf("nextValue extra: %s", err)
		}
		if len(item) != len(jsonBig) {
			GinkgoT().Errorf("invalid item: %d %d", len(item), len(jsonBig))
		}
		if string(rest) != "HELLO WORLD" {
			GinkgoT().Errorf("invalid rest: %d", len(rest))
		}
	})
})

type example struct {
	compact string
	indent  string
}

var examples = []example{
	{`1`, `1`},
	{`{}`, `{}`},
	{`[]`, `[]`},
	{`{"":2}`, "{\n\t\"\": 2\n}"},
	{`[3]`, "[\n\t3\n]"},
	{`[1,2,3]`, "[\n\t1,\n\t2,\n\t3\n]"},
	{`{"x":1}`, "{\n\t\"x\": 1\n}"},
	{ex1, ex1i},
}

var ex1 = `[true,false,null,"x",1,1.5,0,-5e+2]`

var ex1i = `[
	true,
	false,
	null,
	"x",
	1,
	1.5,
	0,
	-5e+2
]`

type indentErrorTest struct {
	in  string
	err error
}

var indentErrorTests = []indentErrorTest{
	{`{"X": "foo", "Y"}`, &SyntaxError{"invalid character '}' after object key", 17}},
	{`{"X": "foo" "Y": "bar"}`, &SyntaxError{"invalid character '\"' after object key:value pair", 13}},
}

var benchScan scanner

func BenchmarkSkipValue(b *testing.B) {
	initBig()
	for i := 0; i < b.N; i++ {
		nextValue(jsonBig, &benchScan)
	}
	b.SetBytes(int64(len(jsonBig)))
}

func diff(t GinkgoTInterface, a, b []byte) {
	for i := 0; ; i++ {
		if i >= len(a) || i >= len(b) || a[i] != b[i] {
			j := i - 10
			if j < 0 {
				j = 0
			}
			t.Errorf("diverge at %d: «%s» vs «%s»", i, trim(a[j:]), trim(b[j:]))
			return
		}
	}
}

func trim(b []byte) []byte {
	if len(b) > 20 {
		return b[0:20]
	}
	return b
}

var jsonBig []byte

const (
	big   = 10000
	small = 100
)

func initBig() {
	n := big
	if testing.Short() {
		n = small
	}
	if len(jsonBig) != n {
		b, err := Marshal(genValue(n))
		if err != nil {
			panic(err)
		}
		jsonBig = b
	}
}

func genValue(n int) interface{} {
	if n > 1 {
		switch rand.Intn(2) {
		case 0:
			return genArray(n)
		case 1:
			return genMap(n)
		}
	}
	switch rand.Intn(3) {
	case 0:
		return rand.Intn(2) == 0
	case 1:
		return rand.NormFloat64()
	case 2:
		return genString(30)
	}
	panic("unreachable")
}

func genString(stddev float64) string {
	n := int(math.Abs(rand.NormFloat64()*stddev + stddev/2))
	c := make([]rune, n)
	for i := range c {
		f := math.Abs(rand.NormFloat64()*64 + 32)
		if f > 0x10ffff {
			f = 0x10ffff
		}
		c[i] = rune(f)
	}
	return string(c)
}

func genArray(n int) []interface{} {
	f := int(math.Abs(rand.NormFloat64()) * math.Min(10, float64(n/2)))
	if f > n {
		f = n
	}
	x := make([]interface{}, f)
	for i := range x {
		x[i] = genValue(((i+1)*n)/f - (i*n)/f)
	}
	return x
}

func genMap(n int) map[string]interface{} {
	f := int(math.Abs(rand.NormFloat64()) * math.Min(10, float64(n/2)))
	if f > n {
		f = n
	}
	if n > 0 && f == 0 {
		f = 1
	}
	x := make(map[string]interface{})
	for i := 0; i < f; i++ {
		x[genString(10)] = genValue(((i+1)*n)/f - (i*n)/f)
	}
	return x
}
