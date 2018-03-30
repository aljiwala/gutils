package gutils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// JSONString should convert objects to JSON strings.
func JSONString(obj interface{}) (r string) {
	b, _ := json.Marshal(obj)
	s := fmt.Sprintf("%+v", string(b))
	r = strings.Replace(s, `\u003c`, "<", -1)
	r = strings.Replace(r, `\u003e`, ">", -1)
	return
}

// JsonpToJSON should modify JSONP string to json string.
// Usecase: JsonpToJson({a:1,b:2}) -> {"a":1,"b":2}
// Ref: https://stackoverflow.com/a/3840118/4039768 (What is JSONP?)
func JsonpToJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	nextStart := strings.Index(s, "[")
	if nextStart > 0 && start > nextStart {
		start = nextStart
		end = strings.LastIndex(s, "]")
	}
	if end > start && end != -1 && start != -1 {
		s = s[start : end+1]
	}
	s = strings.Replace(s, "\\'", "", -1)
	regexp, _ := regexp.Compile(RegexJsonpToJSON)
	return regexp.ReplaceAllString(s, "\"$1\":")
}
