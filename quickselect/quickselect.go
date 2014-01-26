// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package quickselect contains routines for sorting and selecting
// the k smallest element in an array of floats using quicksort.
//
package quickselect

import (
  "math/rand"
  "time"
)


// quicksort is the top level driver for quicksort. Its role
// currently mainly is to initialize the random number generator
// before calling the main quicksort routine.
func quicksort(array []float64) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  quicksort_h(array, 0, len(array), r)
}



// quickselect selects the kth smallest item from array using 
// a one sided randomized quicksort routine
// NOTE: k has to be a valid index within slice array
func quickselect(array []float64, k int) float64 {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  return quickselect_h(array, k, 0, len(array), r)
}



// quicksort_h is the main recursive quicksort routine
func quickselect_h(array []float64, k int, first int, last int, r *rand.Rand) float64 {

  if first == last-1 {
    return array[first]
  }

  pivot := partition_items(array, first, last, r)
  if k < pivot {
    return quickselect_h(array, k, first, pivot, r)
  } else if k > pivot {
    return quickselect_h(array, k, pivot+1, last, r)
  } else {
    return array[k]
  }
}



// quicksort_h is the main recursive quicksort routine
func quicksort_h(array []float64, first int, last int, r *rand.Rand) {

  if first >= last-1 {
    return
  }

  pivot := partition_items(array, first, last, r)
  quicksort_h(array, first, pivot+1, r)
  quicksort_h(array, pivot+1, last, r)
}



// partition_items partitions the items according to a chosen
// pivot. The pivot is chosen randomly.
func partition_items(array []float64, first, last int, r *rand.Rand) int {

  // pick random pivot and swap it with the first array element
  pivot := first + r.Intn(last-first-1)
  array[first], array[pivot] = array[pivot], array[first]

  j := first+1
  for i := first+1; i < last; i++ {
    if array[i] < array[first] {
      array[j], array[i] = array[i], array[j]
      j++
    }
  }
  array[j-1], array[first] = array[first], array[j-1]

  return j-1
}






