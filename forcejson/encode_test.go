package forcejson

import (
	"bytes"
	"math"
	"reflect"

	. "github.com/onsi/ginkgo"
	"unicode"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("omit empty", func() {

		var o Optionals
		o.Sw = "something"
		o.Mr = map[string]interface{}{}
		o.Mo = map[string]interface{}{}

		got, err := MarshalIndent(&o, "", " ")
		if err != nil {
			GinkgoT().Fatal(err)
		}
		if got := string(got); got != optionalsExpected {
			GinkgoT().Errorf(" got: %s\nwant: %s\n", got, optionalsExpected)
		}
	})
	It("string tag", func() {

		var s StringTag
		s.BoolStr = true
		s.IntStr = 42
		s.StrStr = "xzbit"
		got, err := MarshalIndent(&s, "", " ")
		if err != nil {
			GinkgoT().Fatal(err)
		}
		if got := string(got); got != stringTagExpected {
			GinkgoT().Fatalf(" got: %s\nwant: %s\n", got, stringTagExpected)
		}

		var s2 StringTag
		err = NewDecoder(bytes.NewBuffer(got)).Decode(&s2)
		if err != nil {
			GinkgoT().Fatalf("Decode: %v", err)
		}
		if !reflect.DeepEqual(s, s2) {
			GinkgoT().Fatalf("decode didn't match.\nsource: %#v\nEncoded as:\n%s\ndecode: %#v", s, string(got), s2)
		}
	})
	It("encode renamed byte slice", func() {

		s := renamedByteSlice("abc")
		result, err := Marshal(s)
		if err != nil {
			GinkgoT().Fatal(err)
		}
		expect := `"YWJj"`
		if string(result) != expect {
			GinkgoT().Errorf(" got %s want %s", result, expect)
		}
		r := renamedRenamedByteSlice("abc")
		result, err = Marshal(r)
		if err != nil {
			GinkgoT().Fatal(err)
		}
		if string(result) != expect {
			GinkgoT().Errorf(" got %s want %s", result, expect)
		}
	})
	It("unsupported values", func() {

		for _, v := range unsupportedValues {
			if _, err := Marshal(v); err != nil {
				if _, ok := err.(*UnsupportedValueError); !ok {
					GinkgoT().Errorf("for %v, got %T want UnsupportedValueError", v, err)
				}
			} else {
				GinkgoT().Errorf("for %v, expected error", v)
			}
		}
	})
	It("ref val marshal", func() {

		var s = struct {
			R0 Ref
			R1 *Ref
			R2 RefText
			R3 *RefText
			V0 Val
			V1 *Val
			V2 ValText
			V3 *ValText
		}{
			R0: 12,
			R1: new(Ref),
			R2: 14,
			R3: new(RefText),
			V0: 13,
			V1: new(Val),
			V2: 15,
			V3: new(ValText),
		}
		const want = `{"R0":"ref","R1":"ref","R2":"\"ref\"","R3":"\"ref\"","V0":"val","V1":"val","V2":"\"val\"","V3":"\"val\""}`
		b, err := Marshal(&s)
		if err != nil {
			GinkgoT().Fatalf("Marshal: %v", err)
		}
		if got := string(b); got != want {
			GinkgoT().Errorf("got %q, want %q", got, want)
		}
	})
	It("marshaler escaping", func() {

		var c C
		want := `"\u003c\u0026\u003e"`
		b, err := Marshal(c)
		if err != nil {
			GinkgoT().Fatalf("Marshal(c): %v", err)
		}
		if got := string(b); got != want {
			GinkgoT().Errorf("Marshal(c) = %#q, want %#q", got, want)
		}

		var ct CText
		want = `"\"\u003c\u0026\u003e\""`
		b, err = Marshal(ct)
		if err != nil {
			GinkgoT().Fatalf("Marshal(ct): %v", err)
		}
		if got := string(b); got != want {
			GinkgoT().Errorf("Marshal(ct) = %#q, want %#q", got, want)
		}
	})
	It("anonymous nonstruct", func() {

		var i IntType = 11
		a := MyStruct{i}
		const want = `{"IntType":11}`

		b, err := Marshal(a)
		if err != nil {
			GinkgoT().Fatalf("Marshal: %v", err)
		}
		if got := string(b); got != want {
			GinkgoT().Errorf("got %q, want %q", got, want)
		}
	})
	It("embedded bug", func() {

		v := BugB{
			BugA{"A"},
			"B",
		}
		b, err := Marshal(v)
		if err != nil {
			GinkgoT().Fatal("Marshal:", err)
		}
		want := `{"S":"B"}`
		got := string(b)
		if got != want {
			GinkgoT().Fatalf("Marshal: got %s want %s", got, want)
		}

		x := BugX{
			A: 23,
		}
		b, err = Marshal(x)
		if err != nil {
			GinkgoT().Fatal("Marshal:", err)
		}
		want = `{"A":23}`
		got = string(b)
		if got != want {
			GinkgoT().Fatalf("Marshal: got %s want %s", got, want)
		}
	})
	It("tagged field dominates", func() {

		v := BugY{
			BugA{"BugA"},
			BugD{"BugD"},
		}
		b, err := Marshal(v)
		if err != nil {
			GinkgoT().Fatal("Marshal:", err)
		}
		want := `{"S":"BugD"}`
		got := string(b)
		if got != want {
			GinkgoT().Fatalf("Marshal: got %s want %s", got, want)
		}
	})
	It("duplicated field disappears", func() {

		v := BugZ{
			BugA{"BugA"},
			BugC{"BugC"},
			BugY{
				BugA{"nested BugA"},
				BugD{"nested BugD"},
			},
		}
		b, err := Marshal(v)
		if err != nil {
			GinkgoT().Fatal("Marshal:", err)
		}
		want := `{}`
		got := string(b)
		if got != want {
			GinkgoT().Fatalf("Marshal: got %s want %s", got, want)
		}
	})
	It("string bytes", func() {

		es := &encodeState{}
		var r []rune
		for i := '\u0000'; i <= unicode.MaxRune; i++ {
			r = append(r, i)
		}
		s := string(r) + "\xff\xff\xffhello"
		_, err := es.string(s)
		if err != nil {
			GinkgoT().Fatal(err)
		}

		esBytes := &encodeState{}
		_, err = esBytes.stringBytes([]byte(s))
		if err != nil {
			GinkgoT().Fatal(err)
		}

		enc := es.Buffer.String()
		encBytes := esBytes.Buffer.String()
		if enc != encBytes {
			i := 0
			for i < len(enc) && i < len(encBytes) && enc[i] == encBytes[i] {
				i++
			}
			enc = enc[i:]
			encBytes = encBytes[i:]
			i = 0
			for i < len(enc) && i < len(encBytes) && enc[len(enc)-i-1] == encBytes[len(encBytes)-i-1] {
				i++
			}
			enc = enc[:len(enc)-i]
			encBytes = encBytes[:len(encBytes)-i]

			if len(enc) > 20 {
				enc = enc[:20] + "..."
			}
			if len(encBytes) > 20 {
				encBytes = encBytes[:20] + "..."
			}
			GinkgoT().Errorf("encodings differ at %#q vs %#q", enc, encBytes)
		}
	})
	It("issue6458", func() {

		type Foo struct {
			M RawMessage
		}
		x := Foo{RawMessage(`"foo"`)}

		b, err := Marshal(&x)
		if err != nil {
			GinkgoT().Fatal(err)
		}
		if want := `{"M":"foo"}`; string(b) != want {
			GinkgoT().Errorf("Marshal(&x) = %#q; want %#q", b, want)
		}

		b, err = Marshal(x)
		if err != nil {
			GinkgoT().Fatal(err)
		}

		if want := `{"M":"ImZvbyI="}`; string(b) != want {
			GinkgoT().Errorf("Marshal(x) = %#q; want %#q", b, want)
		}
	})
})

type Optionals struct {
	Sr string `force:"sr"`
	So string `force:"so,omitempty"`
	Sw string `force:"-"`

	Ir int `force:"omitempty"`
	Io int `force:"io,omitempty"`

	Slr []string `force:"slr,random"`
	Slo []string `force:"slo,omitempty"`

	Mr map[string]interface{} `force:"mr"`
	Mo map[string]interface{} `force:",omitempty"`
}

var optionalsExpected = `{
 "sr": "",
 "omitempty": 0,
 "slr": null,
 "mr": {}
}`

type StringTag struct {
	BoolStr bool   `force:",string"`
	IntStr  int64  `force:",string"`
	StrStr  string `force:",string"`
}

var stringTagExpected = `{
 "BoolStr": "true",
 "IntStr": "42",
 "StrStr": "\"xzbit\""
}`

type renamedByte byte
type renamedByteSlice []byte
type renamedRenamedByteSlice []renamedByte

var unsupportedValues = []interface{}{
	math.NaN(),
	math.Inf(-1),
	math.Inf(1),
}

type Ref int

func (*Ref) MarshalJSON() ([]byte, error) {
	return []byte(`"ref"`), nil
}

func (r *Ref) UnmarshalJSON([]byte) error {
	*r = 12
	return nil
}

type Val int

func (Val) MarshalJSON() ([]byte, error) {
	return []byte(`"val"`), nil
}

type RefText int

func (*RefText) MarshalText() ([]byte, error) {
	return []byte(`"ref"`), nil
}

func (r *RefText) UnmarshalText([]byte) error {
	*r = 13
	return nil
}

type ValText int

func (ValText) MarshalText() ([]byte, error) {
	return []byte(`"val"`), nil
}

type C int

func (C) MarshalJSON() ([]byte, error) {
	return []byte(`"<&>"`), nil
}

type CText int

func (CText) MarshalText() ([]byte, error) {
	return []byte(`"<&>"`), nil
}

type IntType int

type MyStruct struct {
	IntType
}

type BugA struct {
	S string
}

type BugB struct {
	BugA
	S string
}

type BugC struct {
	S string
}

type BugX struct {
	A int
	BugA
	BugB
}

type BugD struct {
	XXX string `force:"S"`
}

type BugY struct {
	BugA
	BugD
}

type BugZ struct {
	BugA
	BugC
	BugY
}
