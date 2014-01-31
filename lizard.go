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
  "strings"
  "github.com/haskelladdict/lizard/average"
  "github.com/haskelladdict/lizard/statistic"
)



// define variable used in command line parsing
var averageFiles bool
var columnID int         // id of column to act on, 0 = leftmost columns
var fileStatistic bool
var numWorkers int
var numThreads int
var wantMedian bool      // also compute median when computing statistic via -s
                         // NOTE: median is O(n) on average and requires
                         // memory the size of the data array


func init() {
  flag.BoolVar(&averageFiles, "a", false, "average columns")
  flag.BoolVar(&fileStatistic, "s", false, "compute file statistics")
  flag.IntVar(&columnID, "c", 0, "column id (default : 0)")
  flag.BoolVar(&wantMedian, "m", false, "compute median with -s (default: false)")
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

    stats := statistic.Statistic(inputFiles, columnID, wantMedian, numWorkers)
    for _, stat := range stats {
      if wantMedian {
        fmt.Printf("%s : %8.8f +/- %8.8f  (mean +/- std)\n%s   %8.8f (median) \n",
          stat.Name, stat.Mean, math.Sqrt(stat.Variance), 
          strings.Repeat(" ", len(stat.Name)), stat.Median)
      } else {
        fmt.Printf("%s : %8.8f +/- %8.8f  (mean +/- std)\n", stat.Name,
          stat.Mean, math.Sqrt(stat.Variance))
      }
    }
  }
}


