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
  "bufio"
  "log"
  "os"
  "strconv"
  "strings"
)



// column captures a single data column parsed from an input file
type column []float64



// job described the parsing work to be done by a single worker
type job struct {
  fileName string
  colID int
  results chan<- column
}



// done is used to signal that a worker has finished his assigned
// jobs. The done struct also conveys how many files were successfully
// processed so the analysis routine can do a proper average
type doneStatus struct {
  files_processed int
}



// add_jobs adds all parsing jobs to the work queue (one per data file)
func add_jobs(fileNames []string, colID int, jobs chan<- job, 
  result chan<- column) {
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



// run_jobs does the actual processing of a single job descriptonr,
// i.e., it parses the file and the pushes the content into
// the results channel.
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
  file, err := os.Open(j.fileName)
  if err != nil {
    log.Printf("Warning: Failed to open file %s. Ignoring file.\n",
      j.fileName)
    return false
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  output := make([]float64,0)
  for scanner.Scan() {
    col, err := strconv.ParseFloat(strings.Fields(scanner.Text())[j.colID], 64)
    if err != nil {
      log.Printf("Warning: Failed to parse file %s. Ignoring file.\n",
        j.fileName)
      return false
    }

    output = append(output, col)
  }

  j.results <- output
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
func wait_and_process_results(results <-chan column, done <-chan doneStatus,
  num_workers int) column {

  var output column
  num_cols := 0

  for w := 0; w < num_workers; {
    select {  // Blocking
    case result := <-results:
      output = process_column(result, output)
    case d := <-done:
      num_cols += d.files_processed
      num_workers--
    }
  }

DONE:

  // process any remaining results
  for {
    select {
    case result := <-results:
      output = process_column(result, output)
    default:
      break DONE
    }
  }

  num_cols_f := float64(num_cols)
  for i, v := range output {
    output[i] = v / num_cols_f
  }

  return output
}



// do_average is the main entry point for doing the averaging spawning
// all involved worker goroutines
func Average(fileNames []string, colID int, numWorkers int) []float64 {
  jobs := make(chan job, numWorkers)
  result := make(chan column, len(fileNames))
  done := make(chan doneStatus, numWorkers)

  go add_jobs(fileNames, colID, jobs, result)
  for i := 0; i < numWorkers; i++ {
    go start_jobs(done, jobs)
  }

  return wait_and_process_results(result, done, numWorkers)
}


