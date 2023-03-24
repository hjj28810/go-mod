package utility

func SliceExceptSame[T any](slice1, slice2 []T, compare func(T, T) bool) []T {
	result := []T{}
	for _, v1 := range slice1 {
		found := false
		for _, v2 := range slice2 {
			if compare(v1, v2) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, v1)
		}
	}
	return result
}

func SliceExcept[T1 any, T2 any, T3 any](slice1 []T1, slice2 []T2, compare func(T1, T2) bool, returnValue func(T1) T3) []T3 {
	result := []T3{}
	for _, v1 := range slice1 {
		found := false
		for _, v2 := range slice2 {
			if compare(v1, v2) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, returnValue(v1))
		}
	}
	return result
}

func SliceIntersectSame[T any](slice1, slice2 []T, compare func(T, T) bool) []T {
	result := []T{}
	for _, v1 := range slice1 {
		for _, v2 := range slice2 {
			if compare(v1, v2) {
				result = append(result, v1)
				break
			}
		}
	}
	return result
}

func SliceIntersect[T1 any, T2 any, T3 any](slice1 []T1, slice2 []T2, compare func(T1, T2) bool, returnValue func(T1, T2) T3) []T3 {
	result := []T3{}
	for _, v1 := range slice1 {
		for _, v2 := range slice2 {
			if compare(v1, v2) {
				result = append(result, returnValue(v1, v2))
				break
			}
		}
	}
	return result
}

func SliceUnion[T any](slice1, slice2 []T, compare func(T, T) bool) []T {
	result := []T{}
	result = append(result, slice1...)
	for _, v2 := range slice2 {
		exist := false
		for _, v1 := range slice1 {
			if compare(v1, v2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, v2)
		}
	}
	return result
}
