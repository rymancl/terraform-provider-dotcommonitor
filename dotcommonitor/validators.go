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
			// then verify minutes are between 0 & 1438
			mins := int(d.Minutes())
			if mins < 0 || mins > 1438 {
				es = append(es, fmt.Errorf("%s: \"%v\", converted to \"%v\" minutes, is invalid - from time must be between 0 & 1438 minutes", k, v, mins))
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
			// then verify minutes are between 1 & 1439
			mins := int(d.Minutes())
			if mins < 1 || mins > 1439 {
				es = append(es, fmt.Errorf("%s: \"%v\", converted to \"%v\" minutes, is invalid - to time must be between 1 & 1439 minutes", k, v, mins))
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
func validateIgnoreErrorsCodes(i interface{}, k string) (ws []string, errors []error) {
	v, ok := i.(string)

	// first validate it is of string type
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	parts := strings.Split(v, IgnoreErrorsCodesSeparator)

	for _, item := range parts {
		var singleErrors []error
		var rangeErrors []error

		// first see if current item is a single error code
		if _, errSingle := strconv.Atoi(item); errSingle != nil {
			// item not a valid single code
			singleErrors = append(singleErrors, fmt.Errorf("%s: \"%v\" unable to parse input as single error code", k, v))

			// try item as a range
			r := strings.Split(item, IgnoreErrorsCodesRangeSeparator)
			if len(r) == 2 {
				// validate if input string is a range (ex: 400-499)
				// ensure each part of the range is a valid integer
				from, errFrom := strconv.Atoi(r[0])
				to, errTo := strconv.Atoi(r[1])
				if errFrom != nil {
					rangeErrors = append(rangeErrors, fmt.Errorf("%s: \"%v\" left side \"%v\" of code range is not a valid number", k, v, from))
				}
				if errTo != nil {
					rangeErrors = append(rangeErrors, fmt.Errorf("%s: \"%v\" right side \"%v\" of code range is not a valid number", k, v, to))
				}

				// then validate first part of code range is smaller than the second part
				if from >= to {
					rangeErrors = append(rangeErrors, fmt.Errorf("%s: \"%v\" left side \"%v\" of code range must be smaller than the right side \"%v\"", k, v, from, to))
				}
			} else {
				// range doesn't include two codes
				rangeErrors = append(rangeErrors, fmt.Errorf("%s: \"%v\" a valid range must have two codes separated by a hyphen", k, v))
			}

			// bubble up errors only if all parsing attempts fail
			if len(singleErrors) > 0 && len(rangeErrors) > 0 {
				errors = append(errors, singleErrors...)
				errors = append(errors, rangeErrors...)
			}
		}
	}

	return
}
