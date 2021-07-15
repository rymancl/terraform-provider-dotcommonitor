package dotcommonitor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//////////////////////////////
// Group validators
//////////////////////////////

// validateGroupAddressNumber ... validates that the given number is:
//  1) a string type
//  2) able to be converted to an integer
func validateGroupAddressNumber() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		// first validate it is of string type
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		// then validate it is a valid number
		if _, err := strconv.Atoi(v); err != nil {
			es = append(es, fmt.Errorf("%s: \"%v\" is not a valid number", k, v))
		}

		return
	}
}

// validateGroupAddressCode ... validates that the given number is:
//  1) a string type
//  2) able to be converted to an integer
func validateGroupAddressCode() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		// first validate it is of string type
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		// then validate it is a valid number
		if _, err := strconv.Atoi(v); err != nil {
			es = append(es, fmt.Errorf("%s: \"%v\" is not a valid code", k, v))
		}

		return
	}
}

//////////////////////////////
// Scheduler validators
//////////////////////////////

// detectInvalidSchedulerWeeklyIntervalDays ... detects if a day string is valid to the API
//  See Weekly_Intervals: https://www.dotcom-monitor.com/wiki/knowledge-base/scheduler-operations/
func detectInvalidSchedulerWeeklyIntervalDays(days []string) []string {
	validDays := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	var invalidDays []string
	for _, item := range days {
		if !stringInList(validDays, item) {
			invalidDays = append(invalidDays, item)
		}
	}
	return invalidDays
}

// validateWeeklyIntervalFrom ... ensure weekly interval to is valid
//  See Weekly_Intervals: https://www.dotcom-monitor.com/wiki/knowledge-base/scheduler-operations/
func validateWeeklyIntervalFrom() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		// first validate it is of string type
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		// then verify input string can be parsed into a Duration
		if d, err := time.ParseDuration(v); err != nil {
			es = append(es, fmt.Errorf("%s: unable to parse \"%v\" as duration, must be in the format of #h#m", k, v))
		} else {
			// then verify minutes are between 0 & 1439
			mins := int(d.Minutes())
			if mins < 0 || mins > 1439 {
				es = append(es, fmt.Errorf("%s: \"%v\", converted to \"%v\" minutes, is invalid - from time must be between 0 & 1439 minutes", k, v, mins))
			}
		}

		return
	}
}

// validateWeeklyIntervalTo ... ensure weekly interval to is valid
//  See Weekly_Intervals: https://www.dotcom-monitor.com/wiki/knowledge-base/scheduler-operations/
func validateWeeklyIntervalTo() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)

		// first validate it is of string type
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		// then verify input string can be parsed into a Duration
		if d, err := time.ParseDuration(v); err != nil {
			es = append(es, fmt.Errorf("%s: unable to parse \"%v\" as duration, must be in the format of #h#m", k, v))
		} else {
			// then verify minutes are between 1 & 1440
			mins := int(d.Minutes())
			if mins < 0 || mins > 1440 {
				es = append(es, fmt.Errorf("%s: \"%v\", converted to \"%v\" minutes, is invalid - to time must be between 1 & 1440 minutes", k, v, mins))
			}
		}
		return
	}
}

// validateExcludedTimeIntervalTimestamp ... ensure excluded time interval is valid
func validateExcludedTimeIntervalTimestamp(i interface{}, k string) (ws []string, errors []error) {
	v := i.(string)

	// validate input can be parsed to iso8601
	_, err := time.Parse(schedulerExcludedTimeIntervalLayout, v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q cannot be parsed as iso8601 Timestamp Format", v))
	}

	return
}

//////////////////////////////
// Filter validators
//////////////////////////////

// validateIgnoreErrorsCodes ... ensures the ignore errors codes are valid
/*  NOTE: This was originally written under the assumption that codes could be input as a range or
		individual values as you can in the console. It turns out the API doesn't actually support
		what the website "accepts" as valid.
		Leaving this function here for when Dotcom supports input as expected.
*/
func validateIgnoreErrorsCodes(i interface{}, k string) (ws []string, errors []error) {
	v, ok := i.(string)

	// first validate it is of string type
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	split := strings.Split(v, "-")

	// then validate if input string is a single integer (ex: 400)
	if len(split) == 1 {
		if _, err := strconv.Atoi(split[0]); err != nil {
			errors = append(errors, fmt.Errorf("%s: \"%v\" single code is not a valid number", k, v))
		}
	} else if len(split) == 2 {
		// then validate if input string is a range (ex: 400-499)
		// ensure each part of the range is a valid integer
		code1, err1 := strconv.Atoi(split[0])
		code2, err2 := strconv.Atoi(split[1])
		if err1 != nil {
			errors = append(errors, fmt.Errorf("%s: \"%v\" left side of code range is not a valid number", k, v))
		}
		if err2 != nil {
			errors = append(errors, fmt.Errorf("%s: \"%v\" right side of code range is not a valid number", k, v))
		}

		// then validate first part of code range is smaller than the second part
		if code1 >= code2 {
			errors = append(errors, fmt.Errorf("%s: \"%v\" left side of code range must be smaller than the right side", k, v))
		}
	} else {
		// then if we fall here, we consider the input invalid
		errors = append(errors, fmt.Errorf("%s: \"%v\" must either be a single erro code or a hyphen-separated range of codes", k, v))
	}

	return
}
