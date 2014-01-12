// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// lizard provides functionality for analyzing column based numerical data
package main

import (
  "fmt"
  "os"
  "github.com/haskelladdict/lizard/average"
)


func main() {

  avg := average.Average(os.Args[1:], 0, 8)

  for _, v := range avg {
    fmt.Printf("%8.4f\n", v)
  }
}


