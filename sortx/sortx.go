package gutils

// QSortT should sort slices type []T from quickly small to large.
func QSortT(arr interface{}, start2End ...int) {
	var start, end, low, high int

	switch _arr := arr.(type) {
	case []int:
		switch len(start2End) {
		case 0:
			start = 0
			end = len(_arr) - 1
		case 1:
			start = start2End[0]
			end = len(_arr) - 1
		default:
			start = start2End[0]
			end = start2End[1]
		}
		low = start
		high = end
		key := _arr[start]

		for {
			for low < high {
				if _arr[high] < key {
					_arr[low] = _arr[high]
					break
				}
				high--
			}
			for low < high {
				if _arr[low] > key {
					_arr[high] = _arr[low]
					break
				}
				low++
			}
			if low >= high {
				_arr[low] = key
				break
			}
		}

	case []uint64:
		switch len(start2End) {
		case 0:
			start = 0
			end = len(_arr) - 1
		case 1:
			start = start2End[0]
			end = len(_arr) - 1
		default:
			start = start2End[0]
			end = start2End[1]
		}
		low = start
		high = end
		key := _arr[start]

		for {
			for low < high {
				if _arr[high] < key {
					_arr[low] = _arr[high]
					break
				}
				high--
			}
			for low < high {
				if _arr[low] > key {
					_arr[high] = _arr[low]
					break
				}
				low++
			}
			if low >= high {
				_arr[low] = key
				break
			}
		}

	case []string:
		switch len(start2End) {
		case 0:
			start = 0
			end = len(_arr) - 1
		case 1:
			start = start2End[0]
			end = len(_arr) - 1
		default:
			start = start2End[0]
			end = start2End[1]
		}
		low = start
		high = end
		key := _arr[start]

		for {
			for low < high {
				if _arr[high] < key {
					_arr[low] = _arr[high]
					break
				}
				high--
			}
			for low < high {
				if _arr[low] > key {
					_arr[high] = _arr[low]
					break
				}
				low++
			}
			if low >= high {
				_arr[low] = key
				break
			}
		}
	default:
		return
	}

	if low-1 > start {
		QSortT(arr, start, low-1)
	}
	if high+1 < end {
		QSortT(arr, high+1, end)
	}
}
