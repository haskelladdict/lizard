// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package statistic provides functions for computing standard                
// statistic properties (mean, std, ...) for a slice of floats
//
// NOTE: File processing is done via goroutines using a nummber of
//       workers
package statistic

import (
  "math"
  "testing"
)


// Tests for mean and standard deviation
func Test_Average_1(t *testing.T) {

  data_file_1 := "test_files/test_data_1.txt"
  result_1 := Statistic([]string{data_file_1}, 0, true, 4)
  s1_1 := stat{data_file_1, 5.5, 9.166666666666666, 5.5}
  expected_1 := []stat{s1_1}
  if !stat_equal(result_1, expected_1) {
    t.Error("Statistic test 1 failed")
  }

  data_file_2 := "test_files/test_data_2.txt"
  result_2 := Statistic([]string{data_file_2}, 0, true, 4)
  s1_2 := stat{data_file_2, 0.41319134487140002, 0.082911176230414732,
               0.337045349500000}
  expected_2 := []stat{s1_2}
  if !stat_equal(result_2, expected_2) {
    t.Error("Statistic test 2 failed")
  }

  data_file_3 := "test_files/test_data_3.txt"
  result_3 := Statistic([]string{data_file_3}, 1, true, 4)
  s1_3 := stat{data_file_3, 0.49905688017419975, 0.083507191091550331,
               0.498817626000000}
  expected_3 := []stat{s1_3}
  if !stat_equal(result_3, expected_3) {
    t.Error("Statistic test 3 failed")
  }
}


// Benchmarks
func Benchmark_Average(t *testing.B) {

  data_file_3 := "test_files/test_data_3.txt"
  Statistic([]string{data_file_3}, 1, true, 4)
}


// Support Functions
//
// stat_equal compares the entries of a slice of stat structures 
// returned from a call to Average with a reference slice stat structure
func stat_equal(s1, s2 []stat) bool {

  if len(s1) != len(s2) {
    return false
  }

  status := true
  for i := 0; i < len(s1); i++ {
    if s1[i].Name != s2[i].Name {
      status = false
    }

    if !float_equal(s1[i].Mean, s2[i].Mean) {
      status = false
    }

    if !float_equal(s1[i].Variance, s2[i].Variance) {
      status = false
    }

    if !float_equal(s1[i].Median, s2[i].Median) {
      status = false
    }

  }

  return status
}



// float_array_equal compares two float numbers for equality
// NOTE: the floating point comparison is based on an epsilon
//       which was chosen empirically so its not rigorous
func float_equal(a1, a2 float64) bool {
  epsilon := 1e-13
  if math.Abs(a2-a1) > epsilon * math.Abs(a1) {
    return false
  }
  return true
}

