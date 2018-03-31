// Package utilsx provides tools and generic functions to use as a util to make
// use if go even easier.
package utilsx

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/aljiwala/gutils/strx"
)

var (
	hostname    string
	getHostname sync.Once
)

// WordCount returns counts of each word from given string.
//
// Example:
//   - utilsx.WordCount("Australia Canada Germany Australia Japan Canada") should
//     return map[Japan:1 Australia:2 Canada:2 Germany:1].
func WordCount(str string) map[string]int {
	return wordCount(str)
}

// EscapeQuotes should escape quotes from given string str.
func EscapeQuotes(str string) string {
	return escapeQuotes(str)
}

// IsURL should check if the string is an URL.
func IsURL(str string) bool {
	return isURL(str)
}

// GetLastElemOfURL should return last segment/element from given URL string.
func GetLastElemOfURL(urlStr string) (string, error) {
	return getLastElemOfURL(urlStr)
}

// MakeHash should convert/make hash of given string.
func MakeHash(str string) (hash string) {
	return makeHashCRC32(str)
}

// HashString ...
func HashString(encoded string) uint64 {
	return hashString(encoded)
}

// MakeUnique should return eigenvalues' string.
func MakeUnique(obj interface{}) string {
	return makeUnique(obj)
}

// MakeMd5 should return encoded string using MD5 checksum.
func MakeMd5(obj interface{}, length int) string {
	return makeMD5(obj, length)
}

// DayOrdinalSuffix should return suffix based on `day` value provided.
func DayOrdinalSuffix(day int) string {
	return dayOrdinalSuffix(day)
}

// GetTimeByZone should return time.Time value as per provided timezone.
// It will return nil time value and error if `In` method panics, otherwise
// original values.
func GetTimeByZone(t time.Time, zoneStr string) (time.Time, error) {
	return getTimeByZone(t, zoneStr)
}

// MakeMsgID creates a new, globally unique message ID, useable as a Message-ID
// as per RFC822/RFC2822.
func MakeMsgID() string {
	return makeMsgID()
}

// BuildFileRequest should make a request with `Content-Type` as a
// multipart-formdata.
func BuildFileRequest(
	urlStr, fieldname, path string, params, headers map[string]string) (
	*http.Request, error) {
	return buildFileRequest(urlStr, fieldname, path, params, headers)
}

// GetEnvWithDefault should return the value of $env from the OS
// and if it's empty, returns default one.
func GetEnvWithDefault(env, def string) (value string) {
	return getEnvWithDefault(env, def)
}

// GetEnvWithDefaultInt return the int value of $env from the OS and
// if it's empty, returns def.
func GetEnvWithDefaultInt(env string, def int) (int, error) {
	return getEnvWithDefaultInt(env, def)
}

// GetEnvWithDefaultBool should return the bool value of $env from the OS
// and if it's empty, returns def.
func GetEnvWithDefaultBool(env string, def bool) (bool, error) {
	return getEnvWithDefaultBool(env, def)
}

// GetEnvWithDefaultDuration return the time duration value of $env from the OS and
// if it's empty, returns def.
func GetEnvWithDefaultDuration(env, def string) (time.Duration, error) {
	return getEnvWithDefaultDuration(env, def)
}

// GetEnvWithDefaultStrings should return a slice of sorted strings from
// the environment or default split on, So "foo,bar" returns ["bar","foo"].
func GetEnvWithDefaultStrings(env, def string) (v []string) {
	return getEnvWithDefaultStrings(env, def)
}

////////////////////////////////////////////////////////////////////////////////

func wordCount(str string) map[string]int {
	counts := make(map[string]int)
	wordList := strings.Fields(str)
	for _, word := range wordList {
		_, ok := counts[word]
		if ok {
			counts[word]++
		} else {
			counts[word] = 1
		}
	}
	return counts
}

func escapeQuotes(str string) string {
	return QuoteEscaper.Replace(str)
}

func isURL(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" || strings.HasPrefix(str, ".") ||
		len(str) <= minURLRuneCount ||
		utf8.RuneCountInString(str) >= maxURLRuneCount {
		return false
	}

	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	if strings.HasPrefix(u.Host, ".") {
		return false
	}

	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return strx.RXURL.MatchString(str)
}

func getLastElemOfURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return path.Base(parsedURL.Path), nil
}

func makeHashCRC32(str string) (hash string) {
	const IEEE = 0xedb88320
	var IEEETable = crc32.MakeTable(IEEE)
	hash = fmt.Sprintf("%x", crc32.Checksum([]byte(str), IEEETable))
	return
}

func hashString(encoded string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(encoded))
	return hash.Sum64()
}

func makeUnique(obj interface{}) string {
	baseString, _ := json.Marshal(obj)
	return strconv.FormatUint(hashString(string(baseString)), 10)
}

func makeMD5(obj interface{}, length int) string {
	if length > 32 {
		length = 32
	}
	h := md5.New()
	baseString, _ := json.Marshal(obj)
	h.Write([]byte(baseString))
	s := hex.EncodeToString(h.Sum(nil))
	return s[:length]
}

func dayOrdinalSuffix(day int) string {
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func getTimeByZone(t time.Time, zoneStr string) (time.Time, error) {
	loc, err := time.LoadLocation(zoneStr)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

func makeMsgID() string {
	getHostname.Do(func() {
		var err error
		if hostname, err = os.Hostname(); err != nil {
			log.Printf("ERROR: Get hostname: %v", err)
			hostname = "localhost"
		}
	})

	now := time.Now()
	return fmt.Sprintf(
		"<%d.%d.%d@%s>", now.Unix(), now.UnixNano(), rand.Int63(), hostname,
	)
}

func buildFileRequest(
	urlStr, fieldname, path string, params, headers map[string]string) (
	*http.Request, error) {
	// Opens the named file for reading.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	// Write to destination file from source file.
	if _, cpErr := io.Copy(part, file); cpErr != nil {
		return nil, cpErr
	}

	// Range over given params and write the given value.
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	// Write the trailing boundary to end line to output and close it.
	if cErr := writer.Close(); cErr != nil {
		return nil, cErr
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	// Set provided header values.
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, err
}

func getEnvWithDefault(env, def string) string {
	v := os.Getenv(env)
	if v == "" {
		return def
	}

	return v
}

func getEnvWithDefaultInt(env string, def int) (int, error) {
	v := os.Getenv(env)
	if v == "" {
		return def, nil
	}

	return strconv.Atoi(v)
}

func getEnvWithDefaultBool(env string, def bool) (bool, error) {
	v := os.Getenv(env)
	if v == "" {
		return def, nil
	}

	return strconv.ParseBool(v)
}

func getEnvWithDefaultDuration(env, def string) (time.Duration, error) {
	v := os.Getenv(env)
	if v == "" {
		v = def
	}

	return time.ParseDuration(v)
}

func getEnvWithDefaultStrings(env, def string) (v []string) {
	env = GetEnvWithDefault(env, def)
	if env == "" {
		return make([]string, 0)
	}

	v = strings.Split(env, ",")
	if !sort.StringsAreSorted(v) {
		sort.Strings(v)
	}

	return v
}

// // SaveToFile should copy the contents of given file to destination file.
// func SaveToFile(r *http.Request, fromFile string, path string) error {
// 	file, _, err := r.FormFile(fromFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
//
// 	_, cErr := io.Copy(f, file)
// 	if cErr != nil {
// 		return cErr
// 	}
//
// 	return nil
// }

// // GetContentFromURL should get the response from specified URL.
// func GetContentFromURL(url string) (*http.Response, error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return response, nil
// }

////////////////////////////////////////////////////////////////////////////////

// TimeTrack will print the execution time of the function.
// Possible Usage(s):
//   - Call `TimeTrack` function using defer statement.
//
// Ref: https://stackoverflow.com/a/45773638/4039768
func TimeTrack(start time.Time) {
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)
	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Printf("%s took %s", funcName, time.Since(start))
}

////////////////////////////////////////////////////////////////////////////////
