package docid

import (
	"fmt"
	"testing"
)

func TestInDomainMap(t *testing.T) {

	for _, k := range topDomains {
		if !inDomainMap(topDomainMap, Bytes(k)) {
			t.Errorf("Bytes(%q) not in topDomainMap", k)
		}
	}
	for _, k := range secondDomains {
		if !inDomainMap(secondDomainMap, Bytes(k)) {
			t.Errorf("Bytes(%q) not in secondDomainMap", k)
		}
	}

	var tests = []struct {
		lmap  lengthBytesSliceMap
		input string
		want  bool
	}{
		{topDomainMap, "com", true},
		{topDomainMap, "1234567", false},
		{topDomainMap, "123", false},
	}

	for _, test := range tests {
		if inDomainMap(test.lmap, Bytes(test.input)) != test.want {
			t.Errorf("inDomainMap(%q, Bytes(%q)) != %v", test.lmap, test.input, test.want)
		}
	}
}

func BenchmarkInDomainMap(b *testing.B) {
	t := Bytes("com")
	for i := 0; i < b.N; i++ {
		inDomainMap(topDomainMap, t)
	}
}

type parseTest struct {
	data     string
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

type parseFunc func(Bytes) (*DocID, error)

func testParse(t *testing.T, f parseFunc, tests *[]parseTest) {
	for _, test := range *tests {
		docid, err := f(Bytes(test.data))
		if test.isErr {
			if err == nil {
				t.Errorf("parse fail: %v expect err != nil", test)
			}
		} else {
			if err != nil {
				t.Errorf("parse fail: %v expect err == nil, err=%v", test, err)
			} else if fmt.Sprintf("%s", docid) != test.docidStr {
				t.Errorf("parse fail: %v expect %s, result is %s", test, test.docidStr, docid.String())
			}
		}
	}

}

func TestFromURLBytes(t *testing.T) {
	testParse(t, FromURLBytes, &testsFromURLBytes)
}

func TestFromDocIDHexBytes(t *testing.T) {
	testParse(t, FromDocIDHexBytes, &testsFromDocIDHexBytes)
}

func TestFromDocIDHexReadableBytes(t *testing.T) {
	testParse(t, FromDocIDHexReadableBytes, &testsFromDocIDHexReadableBytes)
}
func TestFromBytes(t *testing.T) {
	for _, tests := range [](*[]parseTest){&testsFromURLBytes,
		&testsFromDocIDHexBytes,
		&testsFromDocIDHexReadableBytes} {

		testParse(t, FromBytes, tests)
	}

}

func BenchmarkFromBytes(b *testing.B) {
	t := Bytes("http://www.google.com.hk/abc")
	for i := 0; i < b.N; i++ {
		FromBytes(t)
	}
}

func BenchmarkFromURLBytes(b *testing.B) {
	t := Bytes("http://www.google.com.hk/abc")
	for i := 0; i < b.N; i++ {
		FromURLBytes(t)
	}
}
func BenchmarkFromDocIDHexBytes(b *testing.B) {
	t := Bytes("1d5920f4b44b27a8ed646a3334ca891fed646a3334ca891fd3467db131372140")
	for i := 0; i < b.N; i++ {
		FromURLBytes(t)
	}
}

func BenchmarkFromDocIDHexReadableBytes(b *testing.B) {
	t := Bytes("1d5920f4b44b27a8-ed646a3334ca891f-ed646a3334ca891fd3467db131372140")
	for i := 0; i < b.N; i++ {
		FromURLBytes(t)
	}
}
