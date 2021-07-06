package dotcommonitor

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// groupAddressNumberIsValid ... validates that the given number is:
//  1) a string type
//  2) between 1 and 16 characters
//  3) able to be converted to an integer
func groupAddressNumberIsValid() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		// first validate it is of string type
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		// then validate length
		if len(v) < 1 || len(v) > 16 {
			es = append(es, fmt.Errorf("%s: invalid number of characters for number - must be between 1 and 16 inclusive", k))
		}

		// then validate it is a valid number
		if _, err := strconv.Atoi(v); err != nil {
			es = append(es, fmt.Errorf("%s: not a valid number", k))
		}

		return
	}
}

// groupAddressCodeIsValid ... validates that the given number is:
//  1) a string type
//  2) exactly 3 characters
//  3) able to be converted to an integer
func groupAddressCodeIsValid() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		// first validate it is of string type
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		// then validate length
		if len(v) != 3 {
			es = append(es, fmt.Errorf("%s: invalid number of characters for code - must be 3", k))
		}

		// then validate it is a valid number
		if _, err := strconv.Atoi(v); err != nil {
			es = append(es, fmt.Errorf("%s: not a valid code", k))
		}

		return
	}
}

// detectInvalidSchedulerWeeklyIntervalDays ... detects if a day string is valid to the API
//  See Weekly_Intervals: https://www.dotcom-monitor.com/wiki/knowledge-base/scheduler-operations/
func detectInvalidSchedulerWeeklyIntervalDays(days []string) []string {
	validDays := []string{"Su","Mo","Tu","We","Th","Fr","Sa"}
	var invalidDays []string
	for _, item := range days {
		if !stringInList(validDays, item) {
			invalidDays = append(invalidDays, item)
		}
	}
	return invalidDays
}