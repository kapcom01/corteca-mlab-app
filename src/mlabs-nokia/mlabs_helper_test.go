package main

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	testCases := []struct {
		name string
		line string
		want string
	}{
		{
			name: "ValidErrorWithRecognizableMessage",
			line: `{"Key":"error","Value":{"Failure":"no available M-Lab servers, server misbehaving"}}`,
			want: "server misbehaving",
		},
		{
			name: "ErrorNotInMappings",
			line: `{"Key":"error","Value":{"Failure":"unknownError"}}`,
			want: "",
		},
		{
			name: "InvalidJSON",
			line: `{"Key":"error", "Value":}`,
			want: "",
		},
		{
			name: "NoErrorKey",
			line: `{"Key":"info","Value":{"Failure":"no such host"}}`,
			want: "",
		},
		{
			name: "EmptyString",
			line: ``,
			want: "",
		},
		{
			name: "NonJSONString",
			line: `This is not a JSON string`,
			want: "",
		},
		{
			name: "ErrorKeyWithInvalidStructure",
			line: `{"Key":"error","Value":"invalidStructure"}`,
			want: "",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLine(tt.line)
			if got != tt.want {
				t.Errorf("parseLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOutput(t *testing.T) {
	testCases := []struct {
		name      string
		output    []byte
		want      Mlab
		wantErr   bool
		diagState string
	}{
		{
			name:   "SuccessfulParsing",
			output: []byte(`{"ServerFQDN":"ndt-mlab1-ber01.mlab-oti.measurement-lab.org","ServerIP":"2600:1901:81f0:266::","Download":{"Throughput":{"Value":10.207,"Unit":"Mbit/s"},"Latency":{"Value":17.951,"Unit":"ms"},"Retransmission":{"Value":0.1196,"Unit":"%"}},"Upload":{"Throughput":{"Value":12.375,"Unit":"Mbit/s"},"Latency":{"Value":23,"Unit":"ms"},"Retransmission":{"Value":0,"Unit":""}}}`),
			want: Mlab{
				ServerFQDN: "ndt-mlab1-ber01.mlab-oti.measurement-lab.org",
				ServerIP:   "2600:1901:81f0:266::",
				Download:   Stream{Throughput: ValueUnit{Value: 10.207, Unit: "Mbit/s"}, Latency: ValueUnit{Value: 17.951, Unit: "ms"}, Retransmission: ValueUnit{Value: 0.1196, Unit: "%"}},
				Upload:     Stream{Throughput: ValueUnit{Value: 12.375, Unit: "Mbit/s"}, Latency: ValueUnit{Value: 23, Unit: "ms"}, Retransmission: ValueUnit{Value: 0, Unit: ""}},
			},
			wantErr:   false,
			diagState: "", // no error sets an empty string
		},
		{
			name:      "InvalidJSONFormat",
			output:    []byte(`{"ServerFQDN":"ndt-mlab1-ber01.mlab-oti.measurement-lab.org",`),
			want:      Mlab{},
			wantErr:   true,
			diagState: ERROR_OTHER,
		},
		{
			name:      "ValidJSONWithUnexpectedStructure",
			output:    []byte(`{"UnexpectedKey":"unexpectedValue"}`),
			want:      Mlab{},
			wantErr:   true,
			diagState: ERROR_OTHER,
		},
		{
			name:      "JSONWithPartialData",
			output:    []byte(`{"ServerFQDN":"ndt-mlab1-ber01.mlab-oti.measurement-lab.org"}`),
			want:      Mlab{ServerFQDN: "ndt-mlab1-ber01.mlab-oti.measurement-lab.org"},
			wantErr:   false,
			diagState: "", // partial data doesn't set an error state
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			speedtest := &Speedtest{}
			got, err := parseOutput(speedtest, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOutput() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOutput() got = %v, want %v", got, tt.want)
			} else if speedtest.DiagnosticsState != tt.diagState {
				t.Errorf("parseOutput() diagnostics state got = %v, want %v", speedtest.DiagnosticsState, tt.diagState)
			}
		})
	}
}

func TestHandleOutputErrors(t *testing.T) {
	testCases := []struct {
		name          string
		output        string
		expectedState string
	}{
		{
			name:          "NoErrorsInOutput",
			output:        `{"ServerFQDN":"example.com","ServerIP":"192.0.2.1"}`,
			expectedState: ERROR_OTHER,
		},
		{
			name:          "RecognizedErrorIOTimeout",
			output:        `{"Key":"error","Value":{"Test":"download","Failure":"i/o timeout"}}`,
			expectedState: ERROR_TIMEOUT,
		},
		{
			name:          "RecognizedErrorInvalidPort",
			output:        `{"Key":"error","Value":{"Test":"upload","Failure":"invalid port"}}`,
			expectedState: ERROR_OTHER,
		},
		{
			name:          "RecognizedErrorServerMisbehaving",
			output:        `{"Key":"error","Value":{"Test":"connect","Failure":"server misbehaving"}}`,
			expectedState: ERROR_INITCONNECTION,
		},
		{
			name:          "RecognizedErrorInvalidURLEscape",
			output:        `{"Key":"error","Value":{"Test":"setup","Failure":"invalid URL escape"}}`,
			expectedState: ERROR_OTHER,
		},
		{
			name:          "RecognizedErrorNoSuchHost",
			output:        `{"Key":"error","Value":{"Test":"resolve","Failure":"no such host"}}`,
			expectedState: ERROR_RESOLVE,
		},
		{
			name:          "UnrecognizedError",
			output:        `{"Key":"error","Value":{"Test":"download","Failure":"unknown error"}}`,
			expectedState: ERROR_OTHER,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			speedtest := &Speedtest{}
			handleOutputErrors([]byte(tt.output), speedtest)
			if speedtest.DiagnosticsState != tt.expectedState {
				t.Errorf("%s: expected state %s, got %s", tt.name, tt.expectedState, speedtest.DiagnosticsState)
			}
		})
	}
}
