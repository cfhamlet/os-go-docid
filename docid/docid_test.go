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
			if err == nil {
				t.Errorf("parse fail: %v expect err != nil", test)
			}
		} else {
			if err != nil {
				t.Errorf("parse fail: %v expect err == nil, err=%v", test, err)
			} else if docid.String() != test.docidStr {
				t.Errorf("parse fail: %v expect %s, result is %s", test, test.docidStr, docid.String())
			}
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

func BenchmarkFromBytes(b *testing.B) {
	t := Bytes("http://www.google.com.hk/abc")
	for i := 0; i < b.N; i++ {
		_, _ = FromBytes(t)
	}
}

func BenchmarkFromURLBytes(b *testing.B) {
	t := Bytes("http://www.google.com.hk/abc")
	for i := 0; i < b.N; i++ {
		_, _ = FromURLBytes(t)
	}
}
func BenchmarkFromDocIDHexBytes(b *testing.B) {
	t := Bytes("1d5920f4b44b27a8ed646a3334ca891fed646a3334ca891fd3467db131372140")
	for i := 0; i < b.N; i++ {
		_, _ = FromURLBytes(t)
	}
}

func BenchmarkFromDocIDHexReadableBytes(b *testing.B) {
	t := Bytes("1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140")
	for i := 0; i < b.N; i++ {
		_, _ = FromURLBytes(t)
	}
}

func BenchmarkNewString(b *testing.B) {
	t := "1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140"
	for i := 0; i < b.N; i++ {
		_, _ = New(t)
	}
}

func BenchmarkNewBytes(b *testing.B) {
	t := Bytes("1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140")
	for i := 0; i < b.N; i++ {
		_, _ = New(t)
	}
}
func BenchmarkNewByteSlice(b *testing.B) {
	t := []byte("1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140")
	for i := 0; i < b.N; i++ {
		_, _ = New(t)
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

func TestSiteID(t *testing.T) {
	url := "http://www.google.com/"
	expect := "ed646a3334ca891f"
	d, _ := New(url)
	r := d.SiteID()
	if r.String() != expect {
		t.Errorf("expect: %s result: %s", expect, r.String())
	}
}
func TestDomainID(t *testing.T) {
	url := "http://www.google.com/"
	expected := "1d5920f4b44b27a8"
	d, _ := New(url)
	r := d.DomainID()
	assert.Equal(t, r.String(), expected)
}
func TestURLID(t *testing.T) {
	url := "http://www.google.com/"
	expected := "ff90821feeb2b02a33a6f9fc8e5f3fcd"
	d, _ := New(url)
	r := d.URLID()
	assert.Equal(t, r.String(), expected)
}
