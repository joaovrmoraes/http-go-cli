package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Method string
	URL    string
	Bearer string
	Data   string
}

var (
	requestHistory []Request
	historyFile    = "request_history.json"
)

func AddToHistory(method, url, bearer, data string) {
	request := Request{Method: method, URL: url, Bearer: bearer, Data: data}
	requestHistory = append(requestHistory, request)

	if len(requestHistory) > 5 {
		requestHistory = requestHistory[1:]
	}

	SaveHistoryToFile()
}

func DisplayHistory() string {
	var builder strings.Builder
	for i, req := range requestHistory {
		builder.WriteString(fmt.Sprintf("%d: %s %s\n", i+1, req.Method, req.URL))
	}
	return builder.String()
}

func SaveHistoryToFile() {
	data, err := json.MarshalIndent(requestHistory, "", "  ")
	if err != nil {
		fmt.Println("Error saving history:", err)
		return
	}

	err = os.WriteFile(historyFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing history to file:", err)
		return
	}
}

func LoadHistoryFromFile() {
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		requestHistory = []Request{}
		return
	}

	data, err := os.ReadFile(historyFile)
	if err != nil {
		fmt.Println("Error reading history file:", err)
		return
	}

	err = json.Unmarshal(data, &requestHistory)
	if err != nil {
		fmt.Println("Error loading history:", err)
		return
	}
}
