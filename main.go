// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "conf", "etc/rafall.conf", "config file (in json format)")
}

func main() {
	flag.Parse()
	_, err := NewGenerator(configFile)
	if err != nil {
		panic(err)
	}
}
