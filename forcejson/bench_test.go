package forcejson

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
})

type codeResponse struct {
	Tree     *codeNode `force:"tree"`
	Username string    `force:"username"`
}

type codeNode struct {
	Name     string      `force:"name"`
	Kids     []*codeNode `force:"kids"`
	CLWeight float64     `force:"cl_weight"`
	Touches  int         `force:"touches"`
	MinT     int64       `force:"min_t"`
	MaxT     int64       `force:"max_t"`
	MeanT    int64       `force:"mean_t"`
}

var codeJSON []byte
var codeStruct codeResponse

func codeInit() {
	f, err := os.Open("testdata/code.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	codeJSON = data

	if err := Unmarshal(codeJSON, &codeStruct); err != nil {
		panic("unmarshal code.json: " + err.Error())
	}

	if data, err = Marshal(&codeStruct); err != nil {
		panic("marshal code.json: " + err.Error())
	}

	if !bytes.Equal(data, codeJSON) {
		println("different lengths", len(data), len(codeJSON))
		for i := 0; i < len(data) && i < len(codeJSON); i++ {
			if data[i] != codeJSON[i] {
				println("re-marshal: changed at byte", i)
				println("orig: ", string(codeJSON[i-10:i+10]))
				println("new: ", string(data[i-10:i+10]))
				break
			}
		}
		panic("re-marshal code.json: different result")
	}
}

func BenchmarkCodeEncoder(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	enc := NewEncoder(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		if err := enc.Encode(&codeStruct); err != nil {
			b.Fatal("Encode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeMarshal(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(&codeStruct); err != nil {
			b.Fatal("Marshal:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeDecoder(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var buf bytes.Buffer
	dec := NewDecoder(&buf)
	var r codeResponse
	for i := 0; i < b.N; i++ {
		buf.Write(codeJSON)

		buf.WriteByte('\n')
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		if err := dec.Decode(&r); err != nil {
			b.Fatal("Decode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeUnmarshal(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	for i := 0; i < b.N; i++ {
		var r codeResponse
		if err := Unmarshal(codeJSON, &r); err != nil {
			b.Fatal("Unmmarshal:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeUnmarshalReuse(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var r codeResponse
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(codeJSON, &r); err != nil {
			b.Fatal("Unmmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalString(b *testing.B) {
	data := []byte(`"hello, world"`)
	var s string

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &s); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalFloat64(b *testing.B) {
	var f float64
	data := []byte(`3.14`)

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &f); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalInt64(b *testing.B) {
	var x int64
	data := []byte(`3`)

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &x); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}
