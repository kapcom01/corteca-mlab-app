//go:build !test

package main

import (
	"flag"
	"fmt"
)

var inputFile string
var input bool
var port string
var ndt7Binary string

type Results struct {
	err        error
	Download   float64
	Upload     float64
	Latency    float64
	ServerName string
}

type Status struct {
	SpeedtestRunning bool
	ResultsValid     bool
	SpeedtestStatus  string
	SpeedtestResults Results
}

var status Status
var c chan Results = make(chan Results)

func init() {
	var help = flag.Bool("help", false, "Show help")
	flag.BoolVar(&input, "input", false, "read file from input instead of running the ndt7-client")
	flag.StringVar(&inputFile, "file", "-", "file name to parse, default is - for stdandard input")
	flag.StringVar(&ndt7Binary, "command", "/bin/ndt7-client", "path of ndt7-client executable")
	flag.StringVar(&port, "port", "18000", "Web UI port, default is 18000")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
}

func main() {
	fmt.Printf("Starting webserver\n")
	RunServer()
}
