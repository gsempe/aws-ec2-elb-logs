package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type Item struct {
	t    float64
	line string
}

var (
	maxResponseProcessingTimeResult []Item
	maxBackendProcessingTimeResult  []Item
	maxRequestProcessingTimeResult  []Item
)

func main() {

	var f string
	var r int
	flag.StringVar(&f, "f", "", "ELB log file to parse cf. http://docs.aws.amazon.com/ElasticLoadBalancing/latest/DeveloperGuide/access-log-collection.html")
	flag.IntVar(&r, "r", 3, "number of results wanted")
	flag.Parse()

	if len(f) == 0 {
		fmt.Println("Provide a log file name via the command line")
		os.Exit(1)
	}

	maxRequestProcessingTimeResult = make([]Item, r)
	maxBackendProcessingTimeResult = make([]Item, r)
	maxResponseProcessingTimeResult = make([]Item, r)
	b, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Printf("Unable to open %s. Get the error: %s", f, err.Error())
		os.Exit(1)
	}
	buf := bytes.NewBuffer(b)
	i := 0
	for {
		i++
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Unable to read a line. Got error %s", err.Error())
				os.Exit(1)
			}
			break
		}
		fields := bytes.Fields(line)
		if len(fields) < 7 {
			fmt.Printf("Not enough fields line %d", i)
		}
		// timestamp | elb | client:port | backend:port |request_processing_time | backend_processing_time | response_processing_time | elb_status_code | backend_status_code | received_bytes| sent_bytes | request | user_agent | ssl_cipher | ssl_protocol
		reqpt := fields[4]
		bpt := fields[5]
		resppt := fields[6]
		request_processing_time, err := strconv.ParseFloat(string(reqpt), 64)
		if err != nil {
			fmt.Printf("Invalid request_processing_time line %d", i)
			continue
		}
		maxRequestProcessingTime(string(line), request_processing_time)
		backend_processing_time, err := strconv.ParseFloat(string(bpt), 64)
		if err != nil {
			fmt.Printf("Invalid backend_processing_time line %d", i)
			continue
		}
		maxBackendProcessingTime(string(line), backend_processing_time)
		response_processing_time, err := strconv.ParseFloat(string(resppt), 64)
		if err != nil {
			fmt.Printf("Invalid response_processing_time line %d", i)
			continue
		}
		maxResponseProcessingTime(string(line), response_processing_time)
	}
	fmt.Println("Max request processing times:")
	for _, item := range maxRequestProcessingTimeResult {
		fmt.Printf("  %s", item.line)
	}
	fmt.Println("Max Backend processing times:")
	for _, item := range maxBackendProcessingTimeResult {
		fmt.Printf("  %s", item.line)
	}
	fmt.Println("Max Response processing times:")
	for _, item := range maxResponseProcessingTimeResult {
		fmt.Printf("  %s", item.line)
	}
}

func maxRequestProcessingTime(line string, t float64) {

	var s = -1
	for i, item := range maxRequestProcessingTimeResult {
		if item.t < t {
			s = i
		}
	}
	if s >= 0 {
		newItem := Item{t: t, line: line}
		maxRequestProcessingTimeResult[s] = newItem
	}
}

func maxBackendProcessingTime(line string, t float64) {

	var s = -1
	for i, item := range maxBackendProcessingTimeResult {
		if item.t < t {
			s = i
		}
	}
	if s >= 0 {
		newItem := Item{t: t, line: line}
		maxBackendProcessingTimeResult[s] = newItem
	}
}

func maxResponseProcessingTime(line string, t float64) {

	var s = -1
	for i, item := range maxResponseProcessingTimeResult {
		if item.t < t {
			s = i
		}
	}
	if s >= 0 {
		newItem := Item{t: t, line: line}
		maxResponseProcessingTimeResult[s] = newItem
	}
}
