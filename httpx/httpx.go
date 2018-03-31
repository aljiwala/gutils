package httpx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aljiwala/gutils/ioutilx"
	"github.com/aljiwala/gutils/utilsx"
	"github.com/pkg/errors"
)

// ErrNotOK is used when the status code is not 200 OK.
type ErrNotOK struct {
	URL string
	Err string
}

func (e ErrNotOK) Error() string {
	return fmt.Sprintf("code %v while downloading %v", e.Err, e.URL)
}

// EnsureHTTPS wraps a HTTP handler and ensures that it was requested over HTTPS.
// If "DISABLED_ENSURE_HTTPS" is in the environment and set to either "1" or "true",
// then EnsureHTTPS should always pass.
func EnsureHTTPS(handler http.HandlerFunc) http.HandlerFunc {
	v := os.Getenv("DISABLE_ENSURE_HTTPS")
	disabled := v == "1" || v == "true"
	return func(w http.ResponseWriter, r *http.Request) {
		if !disabled && (r.URL.Scheme != "https" && r.Header.Get("X-Forwarded-Proto") != "https") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handler(w, r)
	}
}

// CreateFormFile is like multipart.Writer.CreateFormFile, but allows the
// setting of Content-Type.
func CreateFormFile(w *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	eq := utilsx.EscapeQuotes
	h := make(textproto.MIMEHeader)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h.Set("Content-Type", contentType)
	h.Set("Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			eq(fieldname), eq(filename)),
	)

	return w.CreatePart(h)
}

// ReadRequestOneFile reads the first file from the request (if multipart/),
// or returns the body if not
func ReadRequestOneFile(r *http.Request) (body io.ReadCloser, contentType string, status int, err error) {
	body = r.Body
	contentType = r.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "multipart/") {
		// not multipart-form
		status = http.StatusOK
		return
	}
	defer r.Body.Close()

	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		status, err = http.StatusMethodNotAllowed, errors.New("error parsing request as multipart-form: "+err.Error())
		return
	}

	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		status, err = http.StatusMethodNotAllowed, errors.New("no files?")
		return
	}

Outer:
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			if body, err = fileHeader.Open(); err != nil {
				status, err =
					http.StatusMethodNotAllowed,
					fmt.Errorf(
						"error opening part %q: %s", fileHeader.Filename, err,
					)
				return
			}
			contentType = fileHeader.Header.Get("Content-Type")
			break Outer
		}
	}

	status = http.StatusOK
	return
}

// // ReadRequestFiles reads the files from the request, and calls ReaderToFile on them
// func ReadRequestFiles(r *http.Request) (filenames []string, status int, err error) {
// 	defer r.Body.Close()
// 	err = r.ParseMultipartForm(1 << 20)
// 	if err != nil {
// 		status, err =
// 			http.StatusMethodNotAllowed,
// 			errors.New("cannot parse request as multipart-form: "+err.Error())
// 		return
// 	}
// 	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
// 		status, err = http.StatusMethodNotAllowed, errors.New("no files?")
// 		return
// 	}

// 	filenames = make([]string, 0, len(r.MultipartForm.File))
// 	var f multipart.File
// 	var fn string
// 	for _, fileHeaders := range r.MultipartForm.File {
// 		for _, fh := range fileHeaders {
// 			if f, err = fh.Open(); err != nil {
// 				status, err =
// 					http.StatusMethodNotAllowed,
// 					fmt.Errorf("error reading part %q: %s", fh.Filename, err)
// 				return
// 			}

// 			if fn, err = temp.ReaderToFile(f, fh.Filename, ""); err != nil {
// 				f.Close()
// 				status, err =
// 					http.StatusInternalServerError,
// 					fmt.Errorf("error saving %q: %s", fh.Filename, err)
// 				return
// 			}
// 			f.Close()
// 			filenames = append(filenames, fn)
// 		}
// 	}
// 	if len(filenames) == 0 {
// 		status, err = http.StatusMethodNotAllowed, errors.New("no files??")
// 		return
// 	}

// 	status = http.StatusOK
// 	return
// }

// SendFile sends the given file as response
func SendFile(w http.ResponseWriter, filename, contentType string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	fi, err := fh.Stat()
	if err != nil {
		return err
	}

	size := fi.Size()
	if _, err = fh.Seek(0, 0); err != nil {
		err = fmt.Errorf("error seeking in %v: %s", fh, err)
		http.Error(w, err.Error(), 500)
		return err
	}
	if contentType != "" {
		w.Header().Add("Content-Type", contentType)
	}

	w.Header().Add("Content-Length", fmt.Sprintf("%d", size))
	w.WriteHeader(200)
	fh.Seek(0, 0)

	if _, err = io.CopyN(w, fh, size); err != nil {
		err = fmt.Errorf("error sending file %q: %s", filename, err)
	}

	return err
}

// Fetch the contents of an HTTP URL.
//
// This is not intended to cover all possible use cases  for fetching files,
// only the most common ones. Use the net/http package for more advanced usage.
func Fetch(url string) ([]byte, error) {
	client := http.Client{Timeout: 60 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot download %v", url)
	}
	defer response.Body.Close() // nolint: errcheck

	// TODO: Maybe add sanity check to bail out of the Content-Length is very
	// large?
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of %v", url)
	}

	if response.StatusCode != http.StatusOK {
		return data, ErrNotOK{
			URL: url,
			Err: fmt.Sprintf("%v %v", response.StatusCode, response.Status),
		}
	}

	return data, nil
}

// Save an HTTP URL to the directory dir with the filename. The filename can be
// generated from the URL if empty.
//
// It will return the full path to the save file. Note that it may create both a
// file *and* return an error (e.g. in cases of non-200 status codes).
//
// This is not intended to cover all possible use cases  for fetching files,
// only the most common ones. Use the net/http package for more advanced usage.
func Save(url string, dir string, filename string) (string, error) {
	// Use last path of url if filename is empty
	if filename == "" {
		tokens := strings.Split(url, "/")
		filename = tokens[len(tokens)-1]
	}
	path := filepath.FromSlash(dir + "/" + filename)

	client := http.Client{Timeout: 60 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "cannot download %v", url)
	}
	defer response.Body.Close() // nolint: errcheck

	output, err := os.Create(path)
	if err != nil {
		return "", errors.Wrapf(err, "cannot create %v", path)
	}
	defer output.Close() // nolint: errcheck

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return path, errors.Wrapf(err, "cannot read body of %v in to %v", url, path)
	}

	if response.StatusCode != http.StatusOK {
		return path, ErrNotOK{
			URL: url,
			Err: fmt.Sprintf("%v %v", response.StatusCode, response.Status),
		}
	}

	return path, nil
}

// DumpBody reads the body of a HTTP request without consuming it, so it can be
// read again later.
// It will read at most maxSize of bytes. Use -1 to read everything.
//
// It's based on httputil.DumpRequest.
func DumpBody(r *http.Request, maxSize int64) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	save, body, err := ioutilx.DumpReader(r.Body)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	var dest io.Writer = &b

	chunked := len(r.TransferEncoding) > 0 && r.TransferEncoding[0] == "chunked"
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}

	if maxSize < 0 {
		_, err = io.Copy(dest, body)
	} else {
		_, err = io.CopyN(dest, body, maxSize)
		if err == io.EOF {
			err = nil
		}
	}
	if err != nil {
		return nil, err
	}
	if chunked {
		_ = dest.(io.Closer).Close()
		_, _ = io.WriteString(&b, "\r\n")
	}

	r.Body = save
	return b.Bytes(), nil
}

// Header Utils ----------------------------------------------------------------

// Octet types from RFC 2616.
var octetTypes [256]octetType

type octetType byte

const (
	isToken octetType = 1 << iota
	isSpace
)

func init() {
	// OCTET      = <any 8-bit sequence of data>
	// CHAR       = <any US-ASCII character (octets 0 - 127)>
	// CTL        = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
	// CR         = <US-ASCII CR, carriage return (13)>
	// LF         = <US-ASCII LF, linefeed (10)>
	// SP         = <US-ASCII SP, space (32)>
	// HT         = <US-ASCII HT, horizontal-tab (9)>
	// <">        = <US-ASCII double-quote mark (34)>
	// CRLF       = CR LF
	// LWS        = [CRLF] 1*( SP | HT )
	// TEXT       = <any OCTET except CTLs, but including LWS>
	// separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <">
	//              | "/" | "[" | "]" | "?" | "=" | "{" | "}" | SP | HT
	// token      = 1*<any CHAR except CTLs or separators>
	// qdtext     = <any TEXT except <">>

	for c := 0; c < 256; c++ {
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.ContainsRune(" \t\"(),/:;<=>?@[]\\{}", rune(c))
		if strings.ContainsRune(" \t\r\n", rune(c)) {
			t |= isSpace
		}
		if isChar && !isCtl && !isSeparator {
			t |= isToken
		}
		octetTypes[c] = t
	}
}

// Copy returns a shallow copy of the header.
func Copy(header http.Header) http.Header {
	h := make(http.Header)
	for k, vs := range header {
		h[k] = vs
	}
	return h
}

var timeLayouts = []string{"Mon, 02 Jan 2006 15:04:05 GMT", time.RFC850, time.ANSIC}

// ParseTime parses the header as time. The zero value is returned if the
// header is not present or there is an error parsing the
// header.
func ParseTime(header http.Header, key string) time.Time {
	if s := header.Get(key); s != "" {
		for _, layout := range timeLayouts {
			if t, err := time.Parse(layout, s); err == nil {
				return t.UTC()
			}
		}
	}
	return time.Time{}
}

// ParseList parses a comma separated list of values. Commas are ignored in
// quoted strings. Quoted values are not unescaped or unquoted. Whitespace is
// trimmed.
func ParseList(header http.Header, key string) []string {
	var result []string
	for _, s := range header[http.CanonicalHeaderKey(key)] {
		begin := 0
		end := 0
		escape := false
		quote := false
		for i := 0; i < len(s); i++ {
			b := s[i]
			switch {
			case escape:
				escape = false
				end = i + 1
			case quote:
				switch b {
				case '\\':
					escape = true
				case '"':
					quote = false
				}
				end = i + 1
			case b == '"':
				quote = true
				end = i + 1
			case octetTypes[b]&isSpace != 0:
				if begin == end {
					begin = i + 1
					end = begin
				}
			case b == ',':
				if begin < end {
					result = append(result, s[begin:end])
				}
				begin = i + 1
				end = begin
			default:
				end = i + 1
			}
		}
		if begin < end {
			result = append(result, s[begin:end])
		}
	}
	return result
}

// ParseValueAndParams parses a comma separated list of values with optional
// semicolon separated name-value pairs. Content-Type and Content-Disposition
// headers are in this format.
func ParseValueAndParams(header http.Header, key string) (value string, params map[string]string) {
	params = make(map[string]string)
	s := header.Get(key)
	value, s = expectTokenSlash(s)
	if value == "" {
		return
	}
	value = strings.ToLower(value)
	s = skipSpace(s)
	for strings.HasPrefix(s, ";") {
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
		if pkey == "" {
			return
		}
		if !strings.HasPrefix(s, "=") {
			return
		}
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" {
			return
		}
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
	}
	return
}

// AcceptSpec describes an Accept* header.
type AcceptSpec struct {
	Value string
	Q     float64
}

// ParseAccept parses Accept* headers.
func ParseAccept(header http.Header, key string) (specs []AcceptSpec) {
loop:
	for _, s := range header[key] {
		for {
			var spec AcceptSpec
			spec.Value, s = expectTokenSlash(s)
			if spec.Value == "" {
				continue loop
			}
			spec.Q = 1.0
			s = skipSpace(s)
			if strings.HasPrefix(s, ";") {
				s = skipSpace(s[1:])
				if !strings.HasPrefix(s, "q=") {
					continue loop
				}
				spec.Q, s = expectQuality(s[2:])
				if spec.Q < 0.0 {
					continue loop
				}
			}
			specs = append(specs, spec)
			s = skipSpace(s)
			if !strings.HasPrefix(s, ",") {
				continue loop
			}
			s = skipSpace(s[1:])
		}
	}
	return
}

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

func expectToken(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isToken == 0 {
			break
		}
	}
	return s[:i], s[i:]
}

func expectTokenSlash(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		if (octetTypes[b]&isToken == 0) && b != '/' {
			break
		}
	}
	return s[:i], s[i:]
}

func expectQuality(s string) (q float64, rest string) {
	switch {
	case len(s) == 0:
		return -1, ""
	case s[0] == '0':
		q = 0
	case s[0] == '1':
		q = 1
	default:
		return -1, ""
	}
	s = s[1:]
	if !strings.HasPrefix(s, ".") {
		return q, s
	}
	s = s[1:]
	i := 0
	n := 0
	d := 1
	for ; i < len(s); i++ {
		b := s[i]
		if b < '0' || b > '9' {
			break
		}
		n = n*10 + int(b) - '0'
		d *= 10
	}
	return q + float64(n)/float64(d), s[i:]
}

func expectTokenOrQuoted(s string) (value, rest string) {
	pkey, s := expectToken(s)
	if pkey == "" {
		return
	}
	if !strings.HasPrefix(s, "\"") {
		return "", s
	}

	s = s[1:]
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ {
				b := s[i]
				switch {
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				}
			}
			return "", ""
		}
	}
	return "", ""
}

// Set Utils -------------------------------------------------------------------

// Constants for DispositionArgs.
const (
	TypeInline     = "inline"
	TypeAttachment = "attachment"
)

// DispositionArgs are arguments for SetContentDisposition().
type DispositionArgs struct {
	Type     string // disposition-type
	Filename string // filename-parm
	//CreationDate     time.Time // creation-date-parm
	//ModificationDate time.Time // modification-date-parm
	//ReadDate         time.Time // read-date-parm
	//Size             int       // size-parm
}

// SetContentDisposition sets the Content-Disposition header. Any previous value
// will be overwritten.
//
// https://tools.ietf.org/html/rfc2183
// https://tools.ietf.org/html/rfc6266
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
func SetContentDisposition(header http.Header, args DispositionArgs) error {
	if header == nil {
		return errors.New("header is nil map")
	}

	if args.Type == "" {
		return errors.New("the Type field is mandatory")
	}
	if args.Type != TypeInline && args.Type != TypeAttachment {
		return fmt.Errorf("the Type field must be %#v or %#v", TypeInline, TypeAttachment)
	}
	v := args.Type

	if args.Filename != "" {
		// Format filename= according to <quoted-string> as defined in RFC822.
		// We don't don't allow \, and % though. Replacing \ is a slightly lazy
		// way to prevent certain injections in case of user-provided strings
		// (ending the quoting and injecting their own values or even headers).
		// % because some user agents interpret percent-encodings, and others do
		// not (according to the RFC anyway). Finally escape " with \".
		r := strings.NewReplacer("\\", "", "%", "", `"`, `\"`)
		args.Filename = r.Replace(args.Filename)

		// Don't allow unicode.
		ascii, hasUni := hasUnicode(args.Filename)
		v += fmt.Sprintf(`; filename="%v"`, ascii)

		// Add filename* for unicode, encoded according to
		// https://tools.ietf.org/html/rfc5987
		if hasUni {
			v += fmt.Sprintf("; filename*=UTF-8''%v",
				url.QueryEscape(args.Filename))
		}
	}

	header.Set("Content-Disposition", v)
	return nil
}

func hasUnicode(s string) (string, bool) {
	i := 0
	has := false
	deuni := make([]rune, len(s))
	for _, c := range s {
		// TODO: maybe also disallow any escape chars?
		switch {
		case c > 255:
			has = true
		default:
			deuni[i] = c
			i++
		}
	}

	return strings.TrimRight(string(deuni), "\x00"), has
}

// -----------------------------------------------------------------------------
