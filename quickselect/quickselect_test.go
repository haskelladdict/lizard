// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package quickselect contains routines for sorting and selecting
// the k smallest element in an array of floats using quicksort.
//
package quickselect

import (
  "bufio"
//  "log"
  "math"
  "os"
  "strconv"
  "strings"
  "testing"
)


// Tests
func Test_Average_1(t *testing.T) {

  data_file_1 := "test_files/test_data_1.txt"
  items, err := read_file_into_slice(data_file_1)
  if err != nil {
    t.Error("quickselect test 1 failed - error parsing %v", data_file_1)
  }

  for k := 0; k < 100; k++ {
    result_1 := Quickselect(items, 5)
    if !float_equal(result_1, 6) {
      t.Errorf("quickselect test 1 failed - expected 6 got %v\n", result_1)
    }

    result_2 := Quickselect(items, 6)
    if !float_equal(result_2, 7) {
      t.Errorf("quickselect test 1 failed - expected 6 got %v\n", result_2)
    }

    result_3 := Quickselect(items, 0)
    if !float_equal(result_3, 1) {
      t.Errorf("quickselect test 1 failed - expected 6 got %v\n", result_3)
    }

    result_4 := Quickselect(items, 9)
    if !float_equal(result_4, 10) {
      t.Errorf("quickselect test 1 failed - expected 6 got %v\n", result_4)
    }
  }
}


func Test_Average_2(t *testing.T) {

  items := []float64{77}

  for k := 0; k < 100; k++ {
    result_1 := Quickselect(items, 2)
    if !float_equal(result_1, 77) {
      t.Errorf("quickselect test 2 failed - expected 77 got %v\n", result_1)
    }
  }
}


/*
// Benchmarks
func Benchmark_Average(t *testing.B) {

  data_file_3 := "test_files/test_data_3.txt"
  Statistic([]string{data_file_3}, 1, 4)
}
*/

// Support Functions

// read_file_into_slice reads a single column plain text file
// with doubles into a slice and returns it
func read_file_into_slice(fileName string) ([]float64, error) {

  items := make([]float64, 0)

  file, err := os.Open(fileName)
  if err != nil {
    return items, err
  }

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    elems := strings.TrimSpace(strings.Fields(scanner.Text())[0])
    val, err := strconv.ParseFloat(elems, 64)
    if err != nil {
      return items, err
    }
    items = append(items, val)
  }
  return items, nil
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

