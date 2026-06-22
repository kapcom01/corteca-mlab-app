//go:build !test

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (speedtest *Speedtest) Run() {
	var cmd *exec.Cmd
	var output []byte
	var err error
	speedtest.DiagnosticsState = ERROR_OTHER
	fmt.Printf("Running speedtest\n")
	if input {
		fmt.Printf("Use input file: %s\n", inputFile)
		if inputFile == "-" {
			output, _ = io.ReadAll(os.Stdin)
		} else {
			output, _ = os.ReadFile(inputFile)
		}
	} else {
		fmt.Printf("Executing: %s %s %s %s %s %s %s %s \n", ndt7Binary, ARG0, ARG1, ARG2, ARG3, ARG4, ARG5, ARG6)
		cmd = exec.Command(ndt7Binary, ARG0, ARG1, ARG2, ARG3, ARG4, ARG5, ARG6)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Command Output: %s", string(output))
			c <- Results{err, 0.0, 0.0, 0.0, ""}
			return
		}
	}

	mlab, err := parseOutput(speedtest, output)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		c <- Results{err, 0.0, 0.0, 0.0, ""}
		return
	}

	*speedtest = Speedtest{mlab.ServerFQDN, mlab.Download.Throughput.Value, mlab.Upload.Throughput.Value, mlab.Download.Latency.Value, SPEEDTEST_STATE_COMPLETED}

	fmt.Printf("speedtest download: %f\n", speedtest.Download)
	fmt.Printf("speedtest upload: %f\n", speedtest.Upload)
	fmt.Printf("speedtest latency: %f\n", speedtest.Latency)
	fmt.Printf("speedtest server: %s\n", speedtest.ServerName)
	fmt.Printf("speedtest state: %s\n", speedtest.DiagnosticsState)
	if speedtest == nil || speedtest.Download <= 0 || speedtest.Upload <= 0 || speedtest.Latency <= 0 {
		speedtest.DiagnosticsState = ERROR_OTHER
	}

	c <- Results{nil, speedtest.Download, speedtest.Upload, speedtest.Latency, speedtest.ServerName}
}
