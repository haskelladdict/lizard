// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package statistic provides functions for computing standard
// statistic properties (mean, std, ...) for a slice of floats
//
package statistic

import (
  "bufio"
//  "fmt"
  "log"
  "os"
  "strconv"
  "strings"
)



// column captures a single data column parsed from an input file
type column []float64



// stat describes a struct containing the computed statistics
type stat struct {
  Name string
  Mean float64
  Variance float64
}



// job described the parsing work to be done by a single worker
type job struct {
  fileName string
  colID int
  results chan<- stat
}



// done is used to signal that a worker has finished his assigned
// jobs. The done struct also conveys how many files were successfully
// processed so the analysis routine can do a proper average
type doneStatus struct {
  files_processed int
}



// add_jobs adds all parsing jobs to the work queue (one per data file)
func add_jobs(fileNames []string, colID int, jobs chan<- job,
  result chan<- stat) {
  for _, name := range fileNames {
    jobs <- job{name, colID, result}
  }
  close(jobs)
}



// start_jobs starts jobs still in the queue one by one. Each
// worker processes a separate start_jobs goroutine
func start_jobs(done chan<- doneStatus, jobs <-chan job) {
  num_processed := 0
  for job := range jobs {
    success := job.run()
    if success {
      num_processed++
    }
  }
  done <- doneStatus{num_processed}
}



// run_jobs does the actual processing of a single job descriptor,
// i.e., it parses the file and computes the statistic 
//
// NOTE: The computation of the mean and variance uses Welford's method
//       so we can do away with a single pass through the data.
//       see: Donald Knuth's AOCP, Vol 2, page 232, 3rd edition
//
// NOTE: the column parsing code below can panic in case the user supplies
//       an invalid column ID or the data file itself is damaged. In this
//       case we recover and ignore the whole file.
func (j job) run() bool {

  // ignore data file processed by this job in case of a panic
  defer func() bool {
    if r := recover(); r != nil {
      log.Printf("Warning: Failed to parse file %s. Ignoring file. " +
        "Did you pick the correct column??", j.fileName)
    }
    return false
  }()

  // main processing 
  // NOTE: If filename is empty we assume stdin
  var file *os.File
  if j.fileName == "" {
    file = os.Stdin
  } else {
    var err error
    file, err = os.Open(j.fileName)
    if err != nil {
      log.Printf("Warning: Failed to open file %s. Ignoring file.\n",
        j.fileName)
      return false
    }
  }
  defer file.Close()

  // we are using Welford's method do compute mean and variance
  var count int
  var m_old, s_old, m, s float64

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    elems := strings.TrimSpace(strings.Fields(scanner.Text())[j.colID])
    col, err := strconv.ParseFloat(elems, 64)
    if err != nil {
      log.Printf("Warning: Failed to parse file %s. Ignoring file.\n",
        j.fileName)
      return false
    }

    count++
    if count == 1 {
      m_old = col
      m = col
    } else {
      m = m_old + (col - m_old)/float64(count)
      s = s_old + (col - m_old)*(col - m)
      m_old = m
      s_old = s
    }
  }

  j.results <- stat{j.fileName, m, s/float64(count-1)}
  return true
}



// process_column adds a column to the provided accumulator column
func process_column(result column, acc column) column {

  if len(acc) != 0 {

    if len(result) != len(acc) {
      log.Panic("Mismatched column length in data files. Bailing out...")
    }

    for i, v := range result {
      acc[i] += v
    }
  } else {
    acc = result
  }

  return acc
}



// wait_and_process_results starts with data processing (averaging) while 
// waiting for all workers to finish 
func wait_and_process_results(results <-chan stat, done <-chan doneStatus,
  num_workers int) []stat {

  output := make([]stat, 0)

  for w := 0; w < num_workers; {
    select {  // Blocking
    case result := <-results:
      output = append(output, result)
    case <-done:
      num_workers--
    }
  }

DONE:

  // process any remaining results
  for {
    select {
    case result := <-results:
      output = append(output, result)
    default:
      break DONE
    }
  }

  return output
}



// do_average is the main entry point for doing the averaging spawning
// all involved worker goroutines
//
// NOTE: If the list of fileNames is empty we assume input from stdin
func Statistic(fileNames []string, colID int, numWorkers int) []stat {

  jobs := make(chan job, numWorkers)
  result := make(chan stat, len(fileNames))
  done := make(chan doneStatus, numWorkers)

  go add_jobs(fileNames, colID, jobs, result)
  for i := 0; i < numWorkers; i++ {
    go start_jobs(done, jobs)
  }

  return wait_and_process_results(result, done, numWorkers)
}


