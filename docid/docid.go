package docid

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	domainIDLength = 8
	siteIDLength   = 8
	urlIDLength    = 16
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

// Bytes is byte slice
type Bytes []byte

// BytesSlice is slice of Bytes
type BytesSlice []Bytes

type lengthBytesSliceMap map[int]BytesSlice

var topDomainMap = lengthBytesSliceMap{}
var secondDomainMap = lengthBytesSliceMap{}

// topDomains contain top domain parts
// order is important
var topDomains []string = []string{
	"ac", "co",
	"cat", "edu", "net", "biz", "mil", "int", "com", "gov", "org", "pro",
	"name", "aero", "info", "coop", "jobs", "mobi", "arpa",
	"travel", "museum",
}

// secondDomains contain second domain parts
// order is important
var secondDomains []string = []string{
	"ha", "hb", "ac", "sc", "gd", "sd", "he", "ah", "qh", "sh", "hi",
	"bj", "fj", "tj", "xj", "zj", "hk", "hl", "jl", "nm", "hn", "ln",
	"sn", "yn", "co", "mo", "cq", "gs", "js", "tw", "gx", "jx", "nx",
	"sx", "gz", "xz",
	"cat", "edu", "net", "biz", "mil", "int", "com", "gov", "org", "pro",
	"name", "aero", "info", "coop", "jobs", "mobi", "arpa",
	"travel", "museum",
}

func initDomainMap(domains []string, domainMap lengthBytesSliceMap) {
	for _, s := range domains {
		l := len(s)
		if _, ok := domainMap[l]; !ok {
			domainMap[l] = BytesSlice{}
		}
		domainMap[l] = append(domainMap[l], Bytes(s))
	}
}

func init() {
	initDomainMap(topDomains, topDomainMap)
	initDomainMap(secondDomains, secondDomainMap)
}

func inDomainMap(domainMap lengthBytesSliceMap, s Bytes) bool {
	l := len(s)
	if _, ok := domainMap[l]; !ok {
		return false
	}

	bytesSlice := domainMap[l]
	var begin int = 0
	var end int = len(bytesSlice) - 1
	var mid int = -1

	for begin <= end {
		mid = (begin + end) / 2
		b := bytesSlice[mid]
		if s[1] > b[1] {
			begin = mid + 1
		} else if s[1] < b[1] {
			end = mid - 1
		} else {
			if s[0] > b[0] {
				begin = mid + 1
			} else if s[0] < b[0] {
				end = mid - 1
			} else {
				break
			}
		}
	}

	if begin > end {
		return false
	}

	var i int = 2

	b := bytesSlice[mid]
	for i < l && s[i] == b[i] {
		i++
	}

	if i == l {
		return true
	}

	return false
}

func inSecondDomain(s Bytes) bool {
	return inDomainMap(secondDomainMap, s)
}

func inTopDomain(s []byte) bool {
	return inDomainMap(topDomainMap, s)
}

// DomainID is the first 8 bytes of DocID
type DomainID [domainIDLength]byte

// SiteID is the middle 8 bytes of DocID
type SiteID [siteIDLength]byte

// URLID is the last 16 bytes of DocID
type URLID [urlIDLength]byte

// DocID is 32 bytes array
type DocID [docIDLength]byte

func (id DomainID) String() string {
	var buf [domainIDHexLength]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id SiteID) String() string {
	var buf [siteIDHexLength]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id URLID) String() string {
	var buf [urlIDHexLength]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id DocID) String() string {
	var buf [docIDHexLength]byte
	hex.Encode(buf[:], id[0:domainIDLength])
	buf[domainSiteHexSepPos] = SymbolMinus
	hex.Encode(buf[siteIDHexStart:], id[domainIDLength:urlIDLength])
	buf[siteURLHexSepPos] = SymbolMinus
	hex.Encode(buf[urlIDHexStart:], id[urlIDLength:])
	return string(buf[:])
}

func digest(data Bytes) [md5.Size]byte {
	return md5.Sum(data)
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

func splitDomainSite(urlBytes Bytes) (Bytes, Bytes) {
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
				domainHead, domainPreHead, domainPostHead, domainTail = i, i, i, i
				continue
			} else if !findDomain {
				dealDomain, findDomain = true, true
			}

		}
		if dealDomain {
			domainPreHead, domainHead = domainHead, domainPostHead
			domainPostHead, domainTail = domainTail, i
			dealDomain = false
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
func FromURLBytes(urlBytes Bytes) (DocID, error) {
	domain, site := splitDomainSite(urlBytes)
	d := DocID{}
	domainDigest := digest(domain)
	copy(d[:], domainDigest[:domainIDLength])
	siteDigest := digest(site)
	copy(d[domainIDLength:], siteDigest[:siteIDLength])
	urlDigest := digest(urlBytes)
	copy(d[urlIDLength:], urlDigest[:])
	return d, nil
}

func FromBytes(data Bytes) (DocID, error) {
	return FromURLBytes(data)
}
