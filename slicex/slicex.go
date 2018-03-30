package slicex

import (
	"reflect"
	"strings"
)

// takeArg should take arg as interface value and check if it's kind is matched
// with provided kind value.
func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

// SliceMerge merges interface slices to one slice.
func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

// SliceMergeInt merges int slices to one slice.
func SliceMergeInt(slice1, slice2 []int) (c []int) {
	c = append(slice1, slice2...)
	return
}

// SliceMergeInt64 merges int64 slices to one slice.
func SliceMergeInt64(slice1, slice2 []int64) (c []int64) {
	c = append(slice1, slice2...)
	return
}

// SliceMergeString merges string slices to one slice.
func SliceMergeString(slice1, slice2 []string) (c []string) {
	c = append(slice1, slice2...)
	return
}

// SliceContains would return true if v is present in container.
func SliceContains(container []interface{}, v interface{}) bool {
	for _, vv := range container {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceContainsInt should return true if int v is present in int container.
func SliceContainsInt(container []int, v int) bool {
	for _, vv := range container {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceContainsUint should return true if uint v is present in uint container.
func SliceContainsUint(container []uint, v uint) bool {
	for _, vv := range container {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceContainsString should return true if string v is present in string
// container.
func SliceContainsString(container []string, v string) bool {
	for _, vv := range container {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceUniqueInt should return int container with unique values.
func SliceUniqueInt(s []int) []int {
	size := len(s)
	if size == 0 {
		return []int{}
	}

	m := make(map[int]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}

	realLen := len(m)
	ret := make([]int, realLen)

	idx := 0
	for key := range m {
		ret[idx] = key
		idx++
	}

	return ret
}

// SliceUniqueInt64 should return int64 container with unique values.
func SliceUniqueInt64(s []int64) []int64 {
	size := len(s)
	if size == 0 {
		return []int64{}
	}

	m := make(map[int64]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}

	realLen := len(m)
	ret := make([]int64, realLen)

	idx := 0
	for key := range m {
		ret[idx] = key
		idx++
	}

	return ret
}

// SliceUniqueString should return string container with unique values.
func SliceUniqueString(s []string) []string {
	size := len(s)
	if size == 0 {
		return []string{}
	}

	m := make(map[string]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}

	realLen := len(m)
	ret := make([]string, realLen)

	idx := 0
	for key := range m {
		ret[idx] = key
		idx++
	}

	return ret
}

// SliceSumInt should return sum of given int container values.
func SliceSumInt(container []int) (sum int) {
	for _, v := range container {
		sum += v
	}
	return
}

// SliceSumInt64 should return sum of given int64 container values.
func SliceSumInt64(container []int64) (sum int64) {
	for _, v := range container {
		sum += v
	}
	return
}

// SliceSumFloat64 should return sum of given float64 container values.
func SliceSumFloat64(container []float64) (sum float64) {
	for _, v := range container {
		sum += v
	}
	return
}

// Dedup should remove duplicate uint values from slice.
func Dedup(uintSlice []uint) (dedupSlice []uint) {
	for _, value := range uintSlice {
		if !SliceContainsUint(dedupSlice, value) {
			dedupSlice = append(dedupSlice, value)
		}
	}
	return
}

// RemoveDuplicates should remove the duplicate values from slice of string.
// If caseSensitive is true, then, it'll distinguish string "Abc" from "abc".
//
// Input: RemoveDuplicates(&[]string{"abc", "Abc"}, false) -> Output: &[abc]
// Input: RemoveDuplicates(&[]string{"abc", "Abc"}, true)  -> Output: &[abc Abc]
func RemoveDuplicates(list *[]string, caseSensitive bool) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *list {
		if !caseSensitive {
			x = strings.ToLower(x)
		}
		if !found[x] {
			found[x] = true
			(*list)[j] = (*list)[i]
			j++
		}
	}
	*list = (*list)[:j]
}

// Distinct returns the unique vals of a slice.
func Distinct(arr interface{}) (reflect.Value, bool) {
	// Create a slice from our input interface.
	slice, ok := takeArg(arr, reflect.Slice)
	if !ok {
		return reflect.Value{}, ok
	}

	// Put the values of our slice into a map the key's of the map will be the
	// slice's unique values.
	c := slice.Len()
	m := make(map[interface{}]bool)
	for i := 0; i < c; i++ {
		m[slice.Index(i).Interface()] = true
	}

	i := 0
	mapLen := len(m)

	// Create the output slice and populate it with the map's keys
	out := reflect.MakeSlice(reflect.TypeOf(arr), mapLen, mapLen)
	for k := range m {
		v := reflect.ValueOf(k)
		o := out.Index(i)
		o.Set(v)
		i++
	}

	return out, ok
}

// Intersect returns a slice of values that are present in all of the input slices.
func Intersect(arrs ...interface{}) (reflect.Value, bool) {
	// Create a map to count all the instances of the slice elems.
	arrLength := len(arrs)
	var kind reflect.Kind

	tempMap := make(map[interface{}]int)
	for i, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}

		// check to be sure the type hasn't changed
		if i > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		kind = tempArr.Index(0).Kind()

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			if _, ok := tempMap[tempArr.Index(idx).Interface()]; ok {
				tempMap[tempArr.Index(idx).Interface()]++
			} else {
				tempMap[tempArr.Index(idx).Interface()] = 1
			}
		}
	}

	// Find the keys equal to the length of the input args.
	numElems := 0
	for _, v := range tempMap {
		if v == arrLength {
			numElems++
		}
	}

	i := 0
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), numElems, numElems)
	for key, val := range tempMap {
		if val == arrLength {
			v := reflect.ValueOf(key)
			o := out.Index(i)
			o.Set(v)
			i++
		}
	}

	return out, true
}
