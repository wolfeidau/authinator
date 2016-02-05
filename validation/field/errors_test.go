/*
Copyright 2014 The Kubernetes Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package field

import (
	"fmt"
	"testing"
)

func TestMakeFuncs(t *testing.T) {
	testCases := []struct {
		fn       func() *Error
		expected ErrorType
	}{
		{
			func() *Error { return Invalid(NewPath("f"), "v", "d") },
			ErrorTypeInvalid,
		},
		{
			func() *Error { return Duplicate(NewPath("f"), "v") },
			ErrorTypeDuplicate,
		},
		{
			func() *Error { return NotFound(NewPath("f"), "v") },
			ErrorTypeNotFound,
		},
		{
			func() *Error { return Required(NewPath("f"), "d") },
			ErrorTypeRequired,
		},
		{
			func() *Error { return InternalError(NewPath("f"), fmt.Errorf("e")) },
			ErrorTypeInternal,
		},
	}

	for _, testCase := range testCases {
		err := testCase.fn()
		if err.Type != testCase.expected {
			t.Errorf("expected Type %q, got %q", testCase.expected, err.Type)
		}
	}
}
