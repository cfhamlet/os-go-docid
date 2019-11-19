package docid

import "encoding/hex"
import "fmt"

const (
	domainIDLength = 16
	siteIDLength   = 16
	urlIDLength    = 32
	docIDLength    = domainIDLength + siteIDLength + urlIDLength
)

const (
	domainIDHexLength   = domainIDLength * 2
	domainSiteHexSepPos = domainIDHexLength
	siteIDHexLength     = siteIDLength * 2
	siteIDHexStart      = domainSiteHexSepPos + 1
	siteURLHexSepPos    = siteIDHexStart + siteIDHexLength
	urlIDHexLength      = urlIDLength * 2
	urlIDHexStart       = siteURLHexSepPos + 1
	docIDHexLength      = docIDLength*2 + 2
)

// byte symbols
const (
	SymbolMinus byte = '-'
	SymbolSlash byte = '/'
	SymbolDot   byte = '.'
	SymbolColon byte = ':'
)

type byteKey2L = [2]byte
type byteKey3L = [3]byte
type byteKey4L = [4]byte
type byteKey5L = [5]byte

var secondDomainMap2L map[byteKey2L]bool
var secondDomainMap3L map[byteKey3L]bool
var secondDomainMap4L map[byteKey4L]bool
var secondDomainMap5L map[byteKey5L]bool

var topDomainMap2L map[byteKey2L]bool
var topDomainMap3L map[byteKey3L]bool
var topDomainMap4L map[byteKey4L]bool
var topDomainMap5L map[byteKey5L]bool

var topDomains []string = []string{
	"ac", "co",
	"cat", "edu", "net", "biz", "mil", "int", "com", "gov", "org", "pro",
	"name", "aero", "info", "coop", "jobs", "mobi", "arpa",
	"travel", "museum",
}

var secondDomains []string = []string{
	"ha", "hb", "ac", "sc", "gd", "sd", "he", "ah", "qh", "sh", "hi",
	"bj", "fj", "tj", "xj", "zj", "hk", "hl", "jl", "nm", "hn", "ln",
	"sn", "yn", "co", "mo", "cq", "gs", "js", "tw", "gx", "jx", "nx",
	"sx", "gz", "xz",
	"cat", "edu", "net", "biz", "mil", "int", "com", "gov", "org", "pro",
	"name", "aero", "info", "coop", "jobs", "mobi", "arpa",
	"travel", "museum",
}

func init() {
	for _, s := range topDomains {
		switch len(s) {
		case 2:
			var k byteKey2L
			copy(k[:], []byte(s))
			topDomainMap2L[k] = true
		case 3:
			var k byteKey3L
			copy(k[:], []byte(s))
			topDomainMap3L[k] = true
		case 4:
			var k byteKey4L
			copy(k[:], []byte(s))
			topDomainMap4L[k] = true
		case 5:
			var k byteKey5L
			copy(k[:], []byte(s))
			topDomainMap5L[k] = true
		}
	}
}

func init() {
	for _, s := range secondDomains {
		switch len(s) {
		case 2:
			var k byteKey2L
			copy(k[:], []byte(s))
			secondDomainMap2L[k] = true
		case 3:
			var k byteKey3L
			copy(k[:], []byte(s))
			secondDomainMap3L[k] = true
		case 4:
			var k byteKey4L
			copy(k[:], []byte(s))
			secondDomainMap4L[k] = true
		case 5:
			var k byteKey5L
			copy(k[:], []byte(s))
			secondDomainMap5L[k] = true
		}
	}
}

func inSecondDomain(part []byte) bool {
	return false
}

func inTopDomain(part []byte) bool {
	return false
}

// DomainID is the first 16 bytes of DocID
type DomainID [domainIDLength]byte

// SiteID is the middle 16 bytes of DocID
type SiteID [siteIDLength]byte

// URLID is the last 32 bytes of DocID
type URLID [urlIDLength]byte

// DocID is 64 bytes array
type DocID [docIDLength]byte

func (id *DomainID) String() string {
	var buf [domainIDHexLength]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id *SiteID) String() string {
	var buf [siteIDHexLength]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id *URLID) String() string {
	var buf [urlIDHexLength]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id *DocID) String() string {
	var buf [docIDHexLength]byte
	hex.Encode(buf[:], id[0:domainIDLength])
	buf[domainSiteHexSepPos] = SymbolMinus
	hex.Encode(buf[siteIDHexStart:], id[domainIDLength:urlIDLength])
	buf[siteURLHexSepPos] = SymbolMinus
	hex.Encode(buf[urlIDHexStart:], id[urlIDLength:])
	return string(buf[:])
}

// DomainID get the domain ID
func (id *DocID) DomainID() DomainID {
	var d DomainID
	copy(d[:], id[0:domainIDLength])
	return d
}

// SiteID get the Site ID
func (id *DocID) SiteID() SiteID {
	var d SiteID
	copy(d[:], id[domainIDLength:urlIDLength])
	return d
}

// URLID get the URL ID
func (id *DocID) URLID() URLID {
	var d URLID
	copy(d[:], id[urlIDLength:])
	return d
}

func splitDomainSite(urlBytes []byte) (domain []byte, site []byte) {
	urlLength := len(urlBytes)
	var hostHead, hostTail int

	domainHead, domainTail, domainPreHead, domainPostHead := -1, -1, -1, -1
	findDomain, dealDomain := false, false

	i := 0
	for i < urlLength {
		c := urlBytes[i]
		if c == SymbolDot {
			dealDomain = true

		} else if c == SymbolSlash {
			break
		} else if c == SymbolColon {
			if (i+2 < urlLength) && (urlBytes[i+1] == SymbolSlash) && (urlBytes[i+2] == SymbolSlash) {
				i += 3
				domainHead, domainPostHead, domainPostHead, domainTail = i, i, i, i
				continue
			} else if !findDomain {
				dealDomain, findDomain = true, true
			}

		}
		if dealDomain {
			domainPreHead, domainHead = domainHead, domainPostHead
			domainPostHead, domainTail = domainTail, i
		}
		i++
	}

	hostTail = i
	if !findDomain {
		domainPreHead, domainHead = domainHead, domainPostHead
		domainPostHead, domainTail = domainTail, i
	}

	if inSecondDomain(urlBytes[domainHead+1:domainPostHead]) && !inTopDomain(urlBytes[domainPostHead+1:domainTail]) {
		domainHead = domainPreHead
	}

	domainHead++

	return urlBytes[domainHead:domainTail], urlBytes[hostHead:hostTail]
}

// FromBytes get DocID from URL bytes
func FromBytes(urlBytes []byte) DocID {
	fmt.Println(topDomainMap3L)
	return DocID{}
}
