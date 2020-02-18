package utils

import (
	"sort"
	"strconv"
	"strings"
)

//string => []int64  1,2,3,4 => []int64{1,2,3,4}
func convertStringToInt64Slice(s string) ([]int64, error) {
	var list []int64
	for _, v := range strings.Split(strings.TrimSpace(s), ",") {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		} else {
			list = append(list, i)
		}
	}
	return list, nil
}

//[]int64 => string []int64{1,2,3,4} => "1,2,3,4"
func ConvertInt64SliceToString(param []int64) string {
	result := make([]string, len(param))
	for i, v := range param {
		result[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(result, ",")
}

//删除重复的int64
func RemoveDuplicateInt64s(originInt64s []int64) []int64 {
	sort.Slice(originInt64s, func(i, j int) bool {
		return originInt64s[i] > originInt64s[j]
	})
	var int64s []int64
	var tempInt64 int64 = 0
	for _, i := range originInt64s {
		if tempInt64 != i {
			tempInt64 = i
			int64s = append(int64s, i)
		}
	}
	return int64s
}
