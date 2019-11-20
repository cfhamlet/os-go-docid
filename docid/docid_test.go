package docid

import "testing"

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
	}

	for _, test := range tests {
		if inDomainMap(test.lmap, Bytes(test.input)) != test.want {
			t.Errorf("inDomainMap(%q, Bytes(%q)) != %v", test.lmap, test.input, test.want)
		}
	}
}

func TestFromURLBytes(t *testing.T) {

}

func TestFromBytes(t *testing.T) {

}
