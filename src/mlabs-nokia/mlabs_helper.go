package main

import (
	"encoding/json"
	"errors"
	"strings"
)

type ValueUnit struct {
	Value float64 `json:"Value"`
	Unit  string  `json:"Unit"`
}

type Stream struct {
	Throughput     ValueUnit `json:"Throughput"`
	Latency        ValueUnit `json:"Latency"`
	Retransmission ValueUnit `json:"Retransmission"`
}

type Mlab struct {
	ServerFQDN string `json:"ServerFQDN"`
	ServerIP   string `json:"ServerIP"`
	Download   Stream `json:"Download"`
	Upload     Stream `json:"Upload"`
}

var errorMappings = map[string]string{
	"i/o timeout":        ERROR_TIMEOUT,
	"invalid port":       ERROR_OTHER,
	"server misbehaving": ERROR_INITCONNECTION,
	"invalid URL escape": ERROR_OTHER,
	"no such host":       ERROR_RESOLVE,
	"other error":        ERROR_OTHER,
}

const (
	ERROR_TIMEOUT        string = "Error_Timeout"
	ERROR_RESOLVE               = "Error_CannotResolveHostName"
	ERROR_OTHER                 = "Error_Other"
	ERROR_NOROUTETOHOST         = "Error_NoRouteToHost"
	ERROR_INITCONNECTION        = "Error_InitConnectionFailed"
)

const (
	ARG0 string = "-format"
	ARG1        = "json"
	ARG2        = "-quiet"
	ARG3        = "-anonymize.ip"
	ARG4        = "netblock"
	ARG5        = "-scheme"
	ARG6        = "ws"
)

type Speedtest struct {
	ServerName       string
	Download         float64
	Upload           float64
	Latency          float64
	DiagnosticsState string
}

const (
	SPEEDTEST_STATE_REQUESTED = "Requested"
	SPEEDTEST_STATE_COMPLETED = "Completed"
)

func parseLine(line string) string {
	var lineData map[string]interface{}

	if err := json.Unmarshal([]byte(line), &lineData); err != nil {
		// Skips line if it's not a valid JSON
		return ""
	}

	// Searches for "Key" as "error"
	if key, isString := lineData["Key"].(string); isString && key == "error" {
		if value, valueIsMap := lineData["Value"].(map[string]interface{}); valueIsMap {
			if failure, failureExists := value["Failure"].(string); failureExists {
				for errorMessage := range errorMappings {
					if strings.Contains(failure, errorMessage) {
						return errorMessage
					}
				}
			}
		}
	}

	return ""
}

func parseOutput(speedtest *Speedtest, output []byte) (Mlab, error) {
	var mlab Mlab
	var jsonData map[string]interface{}
	if jsonErr := json.Unmarshal(output, &jsonData); jsonErr != nil {
		handleOutputErrors(output, speedtest)
		return Mlab{}, errors.New("Speedtest Failed")
	}

	if _, ok := jsonData["ServerFQDN"]; !ok {
		speedtest.DiagnosticsState = ERROR_OTHER // Set the diagnostics state
		return Mlab{}, errors.New("unexpected JSON structure")
	}

	err := json.Unmarshal(output, &mlab)
	if err != nil {
		speedtest.DiagnosticsState = ERROR_OTHER
		return mlab, err
	}
	return mlab, nil
}

func handleOutputErrors(output []byte, speedtest *Speedtest) {
	speedtest.DiagnosticsState = ERROR_OTHER
	// Split the error in separate lines
	outputLines := strings.Split(string(output), "\n")

	// Checks each error line for specific mlab-error
	for _, line := range outputLines {
		mlabError := parseLine(line)
		if mlabError != "" {
			speedtest.DiagnosticsState = errorMappings[mlabError]
		}
	}
}
