// Copyright (C) 2019 IHS Markit.
// All Rights Reserved
//
// NOTICE: All information contained herein is, and remains
// the property of IHS Markit and its suppliers,
// if any. The intellectual and technical concepts contained
// herein are proprietary to IHS Markit and its suppliers
// and may be covered by U.S. and Foreign Patents, patents in
// process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from IHS Markit.

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
