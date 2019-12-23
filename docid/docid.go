package docid

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// InvalidBytesError value describe invalid bytes
type InvalidBytesError Bytes

func (e InvalidBytesError) Error() string {
	return fmt.Sprintf("docid: invalid Bytes %q", e)
}

const (
	domainIDLength     = 8
	siteIDLength       = 8
	urlIDLength        = 16
	domainSiteIDLength = domainIDLength + siteIDLength
	docIDLength        = domainSiteIDLength + urlIDLength
)

const (
	domainIDHexLength = domainIDLength * 2
	siteIDHexLength   = siteIDLength * 2
	urlIDHexLength    = urlIDLength * 2
	docIDHexLength    = docIDLength * 2
)

const (
	domainSiteHexReadableSepPos = domainIDHexLength
	siteIDHexReadableStart      = domainSiteHexReadableSepPos + 1
	siteURLHexReadableSepPos    = siteIDHexReadableStart + siteIDHexLength
	urlIDHexReadableStart       = siteURLHexReadableSepPos + 1
	docIDHexReadableLength      = docIDLength*2 + 2
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

	return i == l
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

func (docid *DocID) String() string {
	var buf [docIDHexReadableLength]byte
	hex.Encode(buf[:], docid[0:domainIDLength])
	buf[domainSiteHexReadableSepPos] = SymbolMinus
	hex.Encode(buf[siteIDHexReadableStart:], docid[domainIDLength:urlIDLength])
	buf[siteURLHexReadableSepPos] = SymbolMinus
	hex.Encode(buf[urlIDHexReadableStart:], docid[urlIDLength:])
	return string(buf[:])
}

func digest(data Bytes) [md5.Size]byte {
	return md5.Sum(data)
}

// DomainID get the domain ID
func (docid *DocID) DomainID() *DomainID {
	var d DomainID
	copy(d[:], docid[0:domainIDLength])
	return &d
}

// SiteID get the Site ID
func (docid *DocID) SiteID() *SiteID {
	var d SiteID
	copy(d[:], docid[domainIDLength:urlIDLength])
	return &d
}

// URLID get the URL ID
func (docid *DocID) URLID() *URLID {
	var d URLID
	copy(d[:], docid[urlIDLength:])
	return &d
}

func splitDomainSite(urlBytes Bytes) (Bytes, Bytes) {
	urlLength := len(urlBytes)
	var siteHead, siteTail int

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

	siteTail = i
	if !findDomain {
		domainPreHead, domainHead = domainHead, domainPostHead
		domainPostHead, domainTail = domainTail, i
	}

	if inSecondDomain(urlBytes[domainHead+1:domainPostHead]) && !inTopDomain(urlBytes[domainPostHead+1:domainTail]) {
		domainHead = domainPreHead
	}

	domainHead++

	return urlBytes[domainHead:domainTail], urlBytes[siteHead:siteTail]
}

// FromURLBytes parse DocID from URL bytes
func FromURLBytes(urlBytes Bytes) (docid *DocID, err error) {
	defer func() {
		if e := recover(); e != nil {
			docid, err = nil, fmt.Errorf("docid: %w", e.(error))
		}
	}()
	docid = new(DocID)
	domain, site := splitDomainSite(urlBytes)
	domainDigest := digest(domain)
	copy(docid[:], domainDigest[:domainIDLength])
	siteDigest := digest(site)
	copy(docid[domainIDLength:], siteDigest[:siteIDLength])
	urlDigest := digest(urlBytes)
	copy(docid[urlIDLength:], urlDigest[:])
	return docid, nil
}

// FromDocIDHexBytes parse DocID from hex bytes
func FromDocIDHexBytes(data Bytes) (docid *DocID, err error) {
	if len(data) != docIDHexLength {
		return nil, InvalidBytesError(data)
	}

	docid = new(DocID)
	_, err = hex.Decode(docid[:], data)
	return docid, err
}

// FromDocIDHexReadableBytes parse DocID from readable hex bytes
func FromDocIDHexReadableBytes(data Bytes) (docid *DocID, err error) {
	if len(data) != docIDHexReadableLength || data[domainSiteHexReadableSepPos] != SymbolMinus || data[siteURLHexReadableSepPos] != SymbolMinus {
		return nil, InvalidBytesError(data)
	}
	docid = new(DocID)

	_, err = hex.Decode(docid[:], data[0:domainIDHexLength])
	if err == nil {
		_, err = hex.Decode(docid[domainIDLength:], data[siteIDHexReadableStart:siteURLHexReadableSepPos])
		if err == nil {
			_, err = hex.Decode(docid[domainSiteIDLength:], data[urlIDHexReadableStart:])
		}
	}
	return docid, err
}

// FromBytes parse DocID from bytes
func FromBytes(data Bytes) (docid *DocID, err error) {
	docid, err = nil, nil
	switch len(data) {
	case docIDHexLength:
		docid, err = FromDocIDHexBytes(data)
	case docIDHexReadableLength:
		docid, err = FromDocIDHexReadableBytes(data)
	}
	if (docid == nil && err == nil) || err != nil {
		return FromURLBytes(data)
	}
	return docid, err
}

// New create new DocID
func New(data interface{}) (docid *DocID, err error) {
	docid, err = nil, nil
	switch v := data.(type) {
	case []byte:
		return FromBytes(Bytes(v))
	case string:
		return FromBytes(Bytes(v))
	case Bytes:
		return FromBytes(v)
	default:
		err = errors.Errorf("docid: not support type %T", v)
	}
	return docid, err
}
