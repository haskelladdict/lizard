// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// lizard provides functionality for analyzing column based numerical data
package main

import (
  "fmt"
  "flag"
  "github.com/haskelladdict/lizard/average"
)



// define variable used in command line parsing
var averageFiles bool
var columnID int         // id of columne to act on, 0 = leftmost columns
var numWorkers int

func init() {
  flag.BoolVar(&averageFiles, "a", false, "average columns (short)")
  flag.BoolVar(&averageFiles, "average", false, "average columns")

  flag.IntVar(&columnID, "c", 0, "column id (short)")
  flag.IntVar(&columnID, "columnID", 0, "column id")

  flag.IntVar(&numWorkers, "w", 4, "number of worker goroutines (short)")
  flag.IntVar(&numWorkers, "workers", 4, "number of worker goroutines")
}



func main() {

  // parse command line flags 
  flag.Parse()
  if len(flag.Args()) == 0 {
    flag.Usage()
    return
  }


  if averageFiles {
    avg := average.Average(flag.Args(), columnID, numWorkers)
    for _, v := range avg {
      fmt.Printf("%8.4f\n", v)
    }
  }
}


