package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"

	"github.com/fatih/color"
)

type Request struct {
	Method string
	URL    string
	Bearer string
	Data   string
}

func PrintStatus(statusCode int) {
	var statusColor *color.Color

	switch {
	case statusCode >= 200 && statusCode < 300:
		statusColor = color.New(color.FgGreen)
	case statusCode >= 300 && statusCode < 400:
		statusColor = color.New(color.FgYellow)
	case statusCode >= 400 && statusCode < 500:
		statusColor = color.New(color.FgRed)
	default:
		statusColor = color.New(color.FgMagenta)
	}

	statusColor.Printf("Status: %d\n", statusCode)
}

func PrintHeaders(headers http.Header) {
	type header struct {
		Key   string
		Value []string
	}
	var headerList []header
	for key, value := range headers {
		headerList = append(headerList, header{Key: key, Value: value})
	}

	sort.Slice(headerList, func(i, j int) bool {
		return headerList[i].Key < headerList[j].Key
	})

	for _, h := range headerList {
		fmt.Printf("%s: %v\n", h.Key, h.Value)
	}
}

func CleanBody(input []byte) []byte {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAll(input, nil)
}

func FormatJSON(body []byte) ([]byte, error) {
	cleanBody := CleanBody(body)

	var jsonObj interface{}
	err := json.Unmarshal(cleanBody, &jsonObj)
	if err != nil {
		return nil, fmt.Errorf("error converting the response body to JSON: %v", err)
	}

	formattedJSON, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error formatting JSON: %v", err)
	}

	return formattedJSON, nil
}

func SaveToFile(coloredJSON []byte) {
	tmpfile, err := os.CreateTemp("", "response*.json")
	if err != nil {
		fmt.Printf("Error creating the temporary file: %v\n", err)
		return
	}

	if _, err := tmpfile.Write(coloredJSON); err != nil {
		fmt.Printf("Error writing to the temporary file: %v\n", err)
		return
	}

	if err := tmpfile.Close(); err != nil {
		fmt.Printf("Error closing the temporary file: %v\n", err)
		return
	}

	cmd := exec.Command("code", tmpfile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error opening the text editor: %v\n", err)
		return
	}
}
