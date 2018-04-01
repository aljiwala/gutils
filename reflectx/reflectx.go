package gutils

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

// IsExportedName should check if it's an exported - upper case - name or not.
func IsExportedName(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// IsExportedOrBuiltinType should check if this type is exported or a built-in?
func IsExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return IsExportedName(t.Name()) || t.PkgPath() == ""
}
