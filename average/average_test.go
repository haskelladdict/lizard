// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package average processes the columns of an arbitrary number of 
// columns based text files.
//
// NOTE: File processing is done via goroutines using a nummber of
//       workers
package average

import (
  "math"
  "testing"
)


// TestWmatch_1 tests the proper number of occurences of string "whale"
func Test_Average_1(t *testing.T) {


  data_files_1 := []string{"test_files/test_data_1.txt"}
  result_1 := Average(data_files_1, 0, 4)
  expected_1 := []float64{1.0, 1.0, 1.0, 1.0}
  if !float_array_equal(result_1, expected_1) {
    t.Error("Parse test 1: Failed to parse input correctly")
  }


  data_files_2 := []string{"test_files/test_data_1.txt",
    "test_files/test_data_2.txt"}
  result_2 := Average(data_files_2, 0, 4)
  expected_2 := []float64{1.5, 1.5, 1.5, 1.5}
  if !float_array_equal(result_2, expected_2) {
    t.Error("Parse test 2: Failed to parse input correctly")
  }


  data_files_3 := []string{"test_files/test_data_3.txt",
    "test_files/test_data_4.txt", "test_files/test_data_5.txt"}
  result_3 := Average(data_files_3, 1, 4)
  expected_3 := []float64{23805.333333333333, 19121.333333333333,
    24376.0000, 12504.0000, 14620.3333333333333, 24463.6666666666666,
    24673.333333333333, 15413.0000, 10786.666666666666, 18102.6666666666666}
  if !float_array_equal(result_3, expected_3) {
    t.Error("Parse test 3: Failed to parse input correctly")
  }
}


// float_array_equal compares to arrays of float for equality
// NOTE: the floating point comparison is currently based on
// the smallest representable float which is probabably not
// the best way to do this.
func float_array_equal(a1, a2 []float64) bool {
  if len(a1) != len(a2) {
    return false
  }

  for i, v := range a1 {
    if math.Abs(a2[i]-v) > math.SmallestNonzeroFloat64 * math.Abs(v) {
      return false
    }
  }

  return true
}

