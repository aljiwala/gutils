package strx

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/aljiwala/gutils/slicex"
)

// CamelCaseSplit splits the camelcase word and returns a list of words. It also
// supports digits. Both lower camel case and upper camel case are supported.
//
// Examples:
//   "" =>                     [""]
//   "lowercase" =>            ["lowercase"]
//   "Class" =>                ["Class"]
//   "MyClass" =>              ["My", "Class"]
//   "MyC" =>                  ["My", "C"]
//   "HTML" =>                 ["HTML"]
//   "PDFLoader" =>            ["PDF", "Loader"]
//   "AString" =>              ["A", "String"]
//   "SimpleXMLParser" =>      ["Simple", "XML", "Parser"]
//   "vimRPCPlugin" =>         ["vim", "RPC", "Plugin"]
//   "GL11Version" =>          ["GL", "11", "Version"]
//   "99Bottles" =>            ["99", "Bottles"]
//   "May5" =>                 ["May", "5"]
//   "BFG9000" =>              ["BFG", "9000"]
//   "BöseÜberraschung" =>     ["Böse", "Überraschung"]
//   "Two  spaces" =>          ["Two", "  ", "spaces"]
//   "BadUTF8\xe2\xe2\xa1" =>  ["BadUTF8\xe2\xe2\xa1"]
//
// Splitting rules
//
//  1) If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2) Assign all unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3) Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4) Iterate through array of split strings, and if a given string
//     is upper case:
//       if subsequent string is lower case:
//         move last character of upper case string to beginning of
//         lower case string
func CamelCaseSplit(str string) (splitted []string) {
	return camelCaseSplit(str)
}

////////////////////////////////////////////////////////////////////////////////

func camelCaseSplit(str string) (container []string) {
	// don't split invalid utf8
	if !utf8.ValidString(str) {
		return []string{str}
	}

	lastClass := 0
	container = []string{}
	var (
		class int
		runes [][]rune
	)

	// split into fields based on class of unicode character
	for _, r := range str {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			container = append(container, string(s))
		}
	}

	return
}

// IsNull should check if the string is null.
func IsNull(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsEmail should check if the string is an email.
func IsEmail(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
}

// TrimAndLowercase should trim space and lower the string.
func TrimAndLowercase(str string) string {
	return strings.ToLower(strings.Replace(str, " ", "", -1))
}

// TrimRightSpace should return trimmed string after removing newline, tabline,
// et cetera.
func TrimRightSpace(s string) string {
	return strings.TrimRight(string(s), "\r\n\t ")
}

// BytesToString convert []byte type to string type.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes convert string type to []byte type.
// NOTE: Panics; if modify the member value of the []byte.
func StringToBytes(s string) []byte {
	sp := *(*[2]uintptr)(unsafe.Pointer(&s))
	bp := [3]uintptr{sp[0], sp[1], sp[1]}
	return *(*[]byte)(unsafe.Pointer(&bp))
}

// Md5 should return hash string with MD5 checksum.
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HashStr should return hexadecimal encoding string of new `sha1` checksum.
func HashStr(s string) string {
	h := sha1.New()
	if _, err := h.Write([]byte(s)); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// HashPassword should return hash string of password by writing salt to it.
func HashPassword(pwd, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pwd)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Left returns the "n" left characters of the string.
//
// If the string is shorter than "n" it will return the first "n" characters of
// the string with "…" appended. Otherwise the entire string is returned as-is.
func Left(s string, n int) string {
	if n < 0 {
		n = 0
	}
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// SplitStr should return container of strings splitted by comma value (,).
func SplitStr(str string) (result []string) {
	s := strings.Split(str, ",")
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
		if s[i] != "" {
			result = append(result, s[i])
		}
	}
	return
}

// SplitAndRemoveDups should split and return unique set of string values.
func SplitAndRemoveDups(str string) (splitted []string) {
	splitted = SplitStr(str)
	slicex.RemoveDuplicates(&splitted, false)
	return
}

// GenerateRandStr generates the 64-bit long random unique string. Additionaly,
// takes n and strGen to create custom string for the same.
func GenerateRandStr(n int, strGen string) string {
	b := make([]byte, n)
	randNewSource := rand.NewSource(time.Now().UnixNano())

	// A randNewSource.Int63() generates 63 random bits, enough for
	// letterIdxMax characters!
	for i, cache, remain := n-1, randNewSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randNewSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(strGen) {
			b[i] = strGen[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b) + strconv.FormatInt(time.Now().Unix(), 10)
}

// SnakeStr converts the accepted string to a snake string (XxYy to xx_yy).
func SnakeStr(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)

	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}

	return strings.ToLower(string(data[:]))
}

// CamelStr converts the accepted string to a camel string (xx_yy to XxYy).
func CamelStr(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1

	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}

	return string(data[:])
}

// ObjectName gets the type name of the object.
func ObjectName(obj interface{}) string {
	v := reflect.ValueOf(obj)
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}

	return t.String()
}

// JsQueryEscape escapes the string in javascript standard so it can be safely placed
// inside a URL query.
func JsQueryEscape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

// JsQueryUnescape does the inverse transformation of JsQueryEscape, converting
// %AB into the byte 0xAB and '+' into ' ' (space). It returns an error if
// any % is not followed by two hexadecimal digits.
func JsQueryUnescape(s string) (string, error) {
	return url.QueryUnescape(strings.Replace(s, "%20", "+", -1))
}
