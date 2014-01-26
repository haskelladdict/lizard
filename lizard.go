// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// lizard provides functionality for analyzing column based numerical data
package main

import (
  "fmt"
  "flag"
  "math"
  "runtime"
  "github.com/haskelladdict/lizard/average"
  "github.com/haskelladdict/lizard/statistic"
)



// define variable used in command line parsing
var averageFiles bool
var fileStatistic bool
var columnID int         // id of columne to act on, 0 = leftmost columns
var numWorkers int
var numThreads int


func init() {
  flag.BoolVar(&averageFiles, "a", false, "average columns")
  flag.BoolVar(&fileStatistic, "s", false, "compute file statistics")
  flag.IntVar(&columnID, "c", 0, "column id (default : 0)")
  flag.IntVar(&numWorkers, "w", 4, "number of worker goroutines (default: 4)")
  flag.IntVar(&numThreads, "t", runtime.NumCPU(),
    "maximum number of threads (default: number of CPUs")
}



func main() {

  // parse command line flags 
  flag.Parse()

  // set the number of threads for go runtime
  runtime.GOMAXPROCS(numThreads)

  // if there are no input files we assume stdin
  inputFiles := flag.Args()
  if len(inputFiles) == 0 {
    numWorkers = 1
  } else if len(inputFiles) < numWorkers {
    numWorkers = len(inputFiles)
  }

  if averageFiles {
    avg := average.Average(inputFiles, columnID, numWorkers)
    for _, v := range avg {
      fmt.Printf("%8.4f\n", v)
    }
  }

  if fileStatistic {
    // if no input files are provided we assume the user meant stdin
    // which we signal with an empty string
    if len(inputFiles) == 0 {
      inputFiles = append(inputFiles, "")
    }

    stats := statistic.Statistic(inputFiles, columnID, numWorkers)
    for _, stat := range stats {
      fmt.Printf("%s : %8.8f +/- %8.8f  (mean +/- std)\n", stat.Name,
        stat.Mean, math.Sqrt(stat.Variance))
    }
  }
}


