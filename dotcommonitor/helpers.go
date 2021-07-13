package dotcommonitor

import (
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//////////////////////////////
// Common helpers
//////////////////////////////

// convertStringListToIntList ... type asserting a list of string to a list of int
func convertStringListToIntList(csvString string) []int {
	sp := strings.Split(csvString, ",")
	intList := make([]int, len(sp))

	for i, item := range sp {
		intList[i], _ = strconv.Atoi(item)
	}

	return intList
}

// convertInterfaceListToIntList ... type asserting a list of interfaces to a list of int
func convertInterfaceListToIntList(interfaceList []interface{}) []int {
	intList := make([]int, len(interfaceList))

	for i := range interfaceList {
		intList[i] = interfaceList[i].(int)
	}

	return intList
}

// convertInterfaceListToStringList ... type asserting a list of interfaces to a list of string
func convertInterfaceListToStringList(interfaceList []interface{}) []string {
	stringList := make([]string, len(interfaceList))

	for i := range interfaceList {
		stringList[i] = interfaceList[i].(string)
	}

	return stringList
}

// expandIntSet ... type asserting a set to a list of int
func expandIntSet(set *schema.Set) []int {
	return convertInterfaceListToIntList(set.List())
}

// expandStringSet ... type asserting a set to a list of string
func expandStringSet(set *schema.Set) []string {
	return convertInterfaceListToStringList(set.List())
}

// intInList .. checks if the int is in the list of ints
func intInList(intList []int, num int) bool {
	sort.Ints(intList)
	index := sort.Search(len(intList), func(i int) bool { return intList[i] >= num })
	result := (index < len(intList)) && (intList[index] == num)
	return result
}

// stringInList .. checks if the string is in the list of strings
func stringInList(stringList []string, s string) bool {
	sort.Strings(stringList)
	index := sort.Search(len(stringList), func(i int) bool { return stringList[i] >= s })
	result := (index < len(stringList)) && (stringList[index] == s)
	return result
}
