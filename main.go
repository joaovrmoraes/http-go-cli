package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
)

func main() {
	// add flags
	method := flag.String("method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	flag.StringVar(method, "m", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")

	url := flag.String("url", "", "URL to make the HTTPS request")

	bearer := flag.String("bearer", "", "Bearer token to authenticate the request")

	data := flag.String("data", "", "Data to send in the request body (for POST, PUT, etc.)")
	flag.StringVar(data, "d", "", "Data to send in the request body (for POST, PUT, etc.)")

	save := flag.Bool("save", false, "Open the JSON response in the text editor")
	flag.BoolVar(save, "s", false, "Open the JSON response in the text editor")

	flag.Parse()

	// Verify if the URL is empty
	if *url == "" {
		fmt.Println("Please provide a URL using the -url flag")
		os.Exit(1)
	}

	if strings.HasPrefix(*url, ":") {
		*url = "http://localhost" + *url
	}

	// start the timer
	start := time.Now()

	// create new request
	var req *http.Request
	var err error
	if *data != "" {
		req, err = http.NewRequest(*method, *url, bytes.NewBuffer([]byte(*data)))
	} else {
		req, err = http.NewRequest(*method, *url, nil)
	}
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
		os.Exit(1)
	}

	// Add the Content-Type header if the data is not empty
	if *data != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Add("Authorization", "Bearer "+*bearer)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making the request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start).Round(time.Millisecond)
	fmt.Printf("\033[1;32mTime elapsed: %v\033[0m\n", elapsed)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body: %v\n", err)
		os.Exit(1)
	}

	// Convert the response body to JSON
	var jsonObj interface{}

	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		fmt.Printf("Error converting the response body to JSON: %v\n", err)
		os.Exit(1)
	}

	// Format the JSON with colors
	f := colorjson.NewFormatter()
	f.Indent = 2
	coloredJSON, err := f.Marshal(jsonObj)
	if err != nil {
		fmt.Printf("Error formatting the JSON: %v\n", err)
		os.Exit(1)
	}

	if *save {
		// Save the JSON to a temporary file
		tmpfile, err := os.CreateTemp("", "response*.json")

		if err != nil {
			fmt.Printf("Error creating the temporary file: %v\n", err)
			os.Exit(1)
		}

		jsonBytes, err := json.Marshal(jsonObj)
		if err != nil {
			fmt.Printf("Error converting the JSON to bytes: %v\n", err)
			os.Exit(1)
		}

		if _, err := tmpfile.Write(jsonBytes); err != nil {
			fmt.Printf("Error writing to the temporary file: %v\n", err)
			os.Exit(1)
		}

		if err := tmpfile.Close(); err != nil {
			fmt.Printf("Error closing the temporary file: %v\n", err)
			os.Exit(1)
		}

		// Open the temporary file in the editor
		cmd := exec.Command("code", tmpfile.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error opening the text editor: %v\n", err)
			os.Exit(1)
		}

	} else {
		// print the colored JSON
		fmt.Println(string(coloredJSON))
	}

}
