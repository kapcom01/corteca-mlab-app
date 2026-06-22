//go:build !test

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Page struct {
	SpeedtestData  Status
	ButtonDisabled bool
	TestStatus     string
}

var tpl *template.Template

var _speedtestValues = Results{
	Download: 0.0,
	Upload:     0.0,
	Latency:    0.0,
	ServerName: ""}

var _speedtest_status = Status{SpeedtestRunning: false,
	ResultsValid:     false,
	SpeedtestResults: _speedtestValues}

var _speedtest_page = Page{SpeedtestData: _speedtest_status,
	ButtonDisabled: false,
	TestStatus:     "not started"}

func init() {
	tpl = template.Must(template.ParseFiles("/www/home.html"))
}

func handleResults(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate(w, &_speedtest_page)
	case http.MethodPost:
		_speedtest_page = Page{SpeedtestData: Status{true, false, "Speedtest Running", Results{}},
			ButtonDisabled: true, TestStatus: "Speedtest Running"}
		go runSpeedtest()
		renderTemplate(w, &_speedtest_page)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if _speedtest_page.SpeedtestData.SpeedtestRunning == true {
			w.WriteHeader(http.StatusTooEarly)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func renderTemplate(w http.ResponseWriter, page *Page) {

	err := tpl.ExecuteTemplate(w, "home.html", &page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func RunServer() {
	fs := http.FileServer(http.Dir("/www/") + "/static")
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleResults)
	http.HandleFunc("/status", handleStatus)
	fmt.Println("Server Listens on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runSpeedtest() {
	_speedtest_status.SpeedtestRunning = true
	_speedtest_page.TestStatus = "Speedtest Running"
	_speedtest_page.ButtonDisabled = true
	var test Speedtest
	go test.Run()
	select {
	case result := <-c:
		if result.err != nil {
			_speedtest_page = Page{SpeedtestData: Status{false, false, "Error", result},
				ButtonDisabled: false, TestStatus: "Error"}
		} else {
			_speedtest_page = Page{SpeedtestData: Status{false, true, "Finished", result},
				ButtonDisabled: false, TestStatus: "Finished"}
		}
	case <-time.After(60 * 2 * time.Second):
		fmt.Println("Error fetching results!")
	}

}
