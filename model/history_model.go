package model

import (
	"fmt"
	"os"
)

type Request struct {
	Method string
	URL    string
	Bearer string
	Data   string
}

var requestHistory []Request

func AddToHistory(method, url, bearer, data string) {
	request := Request{Method: method, URL: url, Bearer: bearer, Data: data}
	requestHistory = append(requestHistory, request)
	if len(requestHistory) > 5 {
		requestHistory = requestHistory[1:]
	}
}

func DisplayHistory() {
	fmt.Println("Request History:")
	for i, req := range requestHistory {
		fmt.Printf("%d: %s %s\n", i+1, req.Method, req.URL)
	}
}

func SaveHistoryToFile() {
	tempFile, err := os.CreateTemp("", "request_history_*.txt")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer tempFile.Close()

	for _, req := range requestHistory {
		_, err := tempFile.WriteString(fmt.Sprintf("%s %s\n", req.Method, req.URL))
		if err != nil {
			fmt.Println("Error writing to temporary file:", err)
			return
		}
	}

	fmt.Println("Request history saved to:", tempFile.Name())
}
