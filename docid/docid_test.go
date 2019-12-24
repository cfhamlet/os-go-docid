package docid

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInDomainMap(t *testing.T) {
	assert := assert.New(t)

	for _, k := range topDomains {
		assert.Equal(inDomainMap(topDomainMap, Bytes(k)), true)
	}
	for _, k := range secondDomains {
		assert.Equal(inDomainMap(secondDomainMap, Bytes(k)), true)
	}

	var tests = []struct {
		lmap     lengthBytesSliceMap
		input    string
		expected bool
	}{
		{topDomainMap, "com", true},
		{topDomainMap, "1234567", false},
		{topDomainMap, "123", false},
	}

	for _, test := range tests {
		assert.Equal(inDomainMap(test.lmap, Bytes(test.input)), test.expected)
	}
}

func BenchmarkInDomainMap(b *testing.B) {
	t := Bytes("com")
	for i := 0; i < b.N; i++ {
		inDomainMap(topDomainMap, t)
	}
}

type parseTest struct {
	data     interface{}
	docidStr string
	isErr    bool
}

var testsFromURLBytes = []parseTest{
	{
		"http://www.google.com/",
		"1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd",
		false,
	},
	{
		"http://www.google.com",
		"1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140",
		false,
	},
	{
		"http://www.google.com.hk/abc",
		"da90da7194dbc779-a735b82241adc4d2-d896d112b3ee45903c11a2cf67d4059f",
		false,
	},
	{
		"1",
		"",
		true,
	},
}

var testsFromDocIDHexBytes = []parseTest{
	{
		"1d5920f4b44b27a8ed646a3334ca891fff90821feeb2b02a33a6f9fc8e5f3fcd",
		"1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd",
		false,
	},
	{
		"1d5920f4b44b27a8ed646a3334ca891fed646a3334ca891fd3467db131372140",
		"1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140",
		false,
	},
	{
		"da90da7194dbc779a735b82241adc4d2d896d112b3ee45903c11a2cf67d4059f",
		"da90da7194dbc779-a735b82241adc4d2-d896d112b3ee45903c11a2cf67d4059f",
		false,
	},
	{
		"1",
		"",
		true,
	},
}
var testsFromDocIDHexReadableBytes = []parseTest{
	{
		"1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd",
		"1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd",
		false,
	},
	{
		"1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140",
		"1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140",
		false,
	},
	{
		"da90da7194dbc779-a735b82241adc4d2-d896d112b3ee45903c11a2cf67d4059f",
		"da90da7194dbc779-a735b82241adc4d2-d896d112b3ee45903c11a2cf67d4059f",
		false,
	},
	{
		"1",
		"",
		true,
	},
}

var testsNew = []parseTest{
	{
		Bytes("http://www.google.com/"),
		"1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd",
		false,
	},
	{
		[]byte("http://www.google.com"),
		"1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140",
		false,
	},
	{
		"http://www.google.com.hk/abc",
		"da90da7194dbc779-a735b82241adc4d2-d896d112b3ee45903c11a2cf67d4059f",
		false,
	},
	{
		"1",
		"",
		true,
	},
	{
		1,
		"",
		true,
	},
}

func testParse(t *testing.T, parseFunc interface{}, tests *[]parseTest, toBytes bool) {
	pfunc := reflect.ValueOf(parseFunc)
	for _, test := range *tests {
		var d interface{} = test.data
		if toBytes {
			switch v := test.data.(type) {
			case string:
				d = Bytes(v)
			case []byte:
				d = Bytes(v)
			}
		}

		results := pfunc.Call([]reflect.Value{reflect.ValueOf(d)})
		docid := results[0].Interface().(*DocID)
		err := results[1].Interface()

		if test.isErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, docid.String(), test.docidStr)
		}
	}
}

func TestFromURLBytes(t *testing.T) {
	testParse(t, FromURLBytes, &testsFromURLBytes, true)
}

func TestFromDocIDHexBytes(t *testing.T) {
	testParse(t, FromDocIDHexBytes, &testsFromDocIDHexBytes, true)
}

func TestFromDocIDHexReadableBytes(t *testing.T) {
	testParse(t, FromDocIDHexReadableBytes, &testsFromDocIDHexReadableBytes, true)
}
func TestFromBytes(t *testing.T) {
	for _, tests := range [](*[]parseTest){
		&testsFromURLBytes,
		&testsFromDocIDHexBytes,
		&testsFromDocIDHexReadableBytes,
	} {
		testParse(t, FromBytes, tests, true)
	}
}

func TestNew(t *testing.T) {
	for _, tests := range [](*[]parseTest){
		&testsFromURLBytes,
		&testsFromDocIDHexBytes,
		&testsFromDocIDHexReadableBytes,
		&testsNew,
	} {
		testParse(t, New, tests, false)
	}
}

func BenchmarkDocIDCreate(b *testing.B) {
	tests := []struct {
		name  string
		input string
		f     func(Bytes) (*DocID, error)
	}{
		{"URLBytes", "http://www.google.com.hk/abc", FromURLBytes},
		{"URL", "http://www.google.com.hk/abc", FromBytes},
		{"DocIDHexBytes", "1d5920f4b44b27a8ed646a3334ca891fed646a3334ca891fd3467db131372140", FromBytes},
		{"DocIDHexReadableBytes", "1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140", FromBytes},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			t := Bytes(test.input)
			for i := 0; i < b.N; i++ {
				_, _ = test.f(t)
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	raw := "1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140"
	tests := []struct {
		name  string
		input interface{}
	}{
		{"string", raw},
		{"Bytes", Bytes(raw)},
		{"[]byte", []byte(raw)},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			switch input := test.input.(type) {
			case string, Bytes, []byte:
				for i := 0; i < b.N; i++ {
					_, _ = New(input)
				}
			}
		})
	}
}

func ExampleNew() {
	docid, _ := New("http://www.google.com/")
	fmt.Println(docid)
	docid, _ = New("1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd")
	fmt.Println(docid)
	docid, _ = New("1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd")
	fmt.Println(docid)
	// Output:
	// 1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd
	// 1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd
	// 1d5920f4b44b27a8-ed646a3334ca891f-ff90821feeb2b02a33a6f9fc8e5f3fcd
}

func TestXID(t *testing.T) {
	tests := []struct {
		funcName string
		input    string
		expected string
	}{
		{"SiteID", "http://www.google.com/", "ed646a3334ca891f"},
		{"DomainID", "http://www.google.com/", "1d5920f4b44b27a8"},
		{"URLID", "http://www.google.com/", "ff90821feeb2b02a33a6f9fc8e5f3fcd"},
	}

	for _, test := range tests {
		d, _ := New(test.input)
		result := reflect.ValueOf(d).MethodByName(test.funcName).Call([]reflect.Value{})
		v := result[0].Interface().(fmt.Stringer).String()
		assert.Equal(t, test.expected, v)
	}
}
