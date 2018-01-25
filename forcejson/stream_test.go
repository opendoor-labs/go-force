package forcejson

import (
	"bytes"
	"io/ioutil"
	"net"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("encoder", func() {

		for i := 0; i <= len(streamTest); i++ {
			var buf bytes.Buffer
			enc := NewEncoder(&buf)
			for j, v := range streamTest[0:i] {
				if err := enc.Encode(v); err != nil {
					GinkgoT().Fatalf("encode #%d: %v", j, err)
				}
			}
			if have, want := buf.String(), nlines(streamEncoded, i); have != want {
				GinkgoT().Errorf("encoding %d items: mismatch", i)
				diff(GinkgoT(), []byte(have), []byte(want))
				break
			}
		}
	})
	It("decoder", func() {

		for i := 0; i <= len(streamTest); i++ {

			var buf bytes.Buffer
			for _, c := range nlines(streamEncoded, i) {
				if c != '\n' {
					buf.WriteRune(c)
				}
			}
			out := make([]interface{}, i)
			dec := NewDecoder(&buf)
			for j := range out {
				if err := dec.Decode(&out[j]); err != nil {
					GinkgoT().Fatalf("decode #%d/%d: %v", j, i, err)
				}
			}
			if !reflect.DeepEqual(out, streamTest[0:i]) {
				GinkgoT().Errorf("decoding %d items: mismatch", i)
				for j := range out {
					if !reflect.DeepEqual(out[j], streamTest[j]) {
						GinkgoT().Errorf("#%d: have %v want %v", j, out[j], streamTest[j])
					}
				}
				break
			}
		}
	})
	It("decoder buffered", func() {

		r := strings.NewReader(`{"Name": "Gopher"} extra `)
		var m struct {
			Name string
		}
		d := NewDecoder(r)
		err := d.Decode(&m)
		if err != nil {
			GinkgoT().Fatal(err)
		}
		if m.Name != "Gopher" {
			GinkgoT().Errorf("Name = %q; want Gopher", m.Name)
		}
		rest, err := ioutil.ReadAll(d.Buffered())
		if err != nil {
			GinkgoT().Fatal(err)
		}
		if g, w := string(rest), " extra "; g != w {
			GinkgoT().Errorf("Remaining = %q; want %q", g, w)
		}
	})
	It("raw message", func() {

		var data struct {
			X  float64
			Id *RawMessage
			Y  float32
		}
		const raw = `["\u0056",null]`
		const msg = `{"X":0.1,"Id":["\u0056",null],"Y":0.2}`
		err := Unmarshal([]byte(msg), &data)
		if err != nil {
			GinkgoT().Fatalf("Unmarshal: %v", err)
		}
		if string([]byte(*data.Id)) != raw {
			GinkgoT().Fatalf("Raw mismatch: have %#q want %#q", []byte(*data.Id), raw)
		}
		b, err := Marshal(&data)
		if err != nil {
			GinkgoT().Fatalf("Marshal: %v", err)
		}
		if string(b) != msg {
			GinkgoT().Fatalf("Marshal: have %#q want %#q", b, msg)
		}
	})
	It("null raw message", func() {

		var data struct {
			X  float64
			Id *RawMessage
			Y  float32
		}
		data.Id = new(RawMessage)
		const msg = `{"X":0.1,"Id":null,"Y":0.2}`
		err := Unmarshal([]byte(msg), &data)
		if err != nil {
			GinkgoT().Fatalf("Unmarshal: %v", err)
		}
		if data.Id != nil {
			GinkgoT().Fatalf("Raw mismatch: have non-nil, want nil")
		}
		b, err := Marshal(&data)
		if err != nil {
			GinkgoT().Fatalf("Marshal: %v", err)
		}
		if string(b) != msg {
			GinkgoT().Fatalf("Marshal: have %#q want %#q", b, msg)
		}
	})
	It("blocking", func() {

		for _, enc := range blockingTests {
			r, w := net.Pipe()
			go w.Write([]byte(enc))
			var val interface{}

			if err := NewDecoder(r).Decode(&val); err != nil {
				GinkgoT().Errorf("decoding %s: %v", enc, err)
			}
			r.Close()
			w.Close()
		}
	})
})
var streamTest = []interface{}{
	0.1,
	"hello",
	nil,
	true,
	false,
	[]interface{}{"a", "b", "c"},
	map[string]interface{}{"K": "Kelvin", "ß": "long s"},
	3.14,
}

var streamEncoded = `0.1
"hello"
null
true
false
["a","b","c"]
{"ß":"long s","K":"Kelvin"}
3.14
`

func nlines(s string, n int) string {
	if n <= 0 {
		return ""
	}
	for i, c := range s {
		if c == '\n' {
			if n--; n == 0 {
				return s[0 : i+1]
			}
		}
	}
	return s
}

var blockingTests = []string{
	`{"x": 1}`,
	`[1, 2, 3]`,
}

func BenchmarkEncoderEncode(b *testing.B) {
	b.ReportAllocs()
	type T struct {
		X, Y string
	}
	v := &T{"foo", "bar"}
	for i := 0; i < b.N; i++ {
		if err := NewEncoder(ioutil.Discard).Encode(v); err != nil {
			b.Fatal(err)
		}
	}
}
