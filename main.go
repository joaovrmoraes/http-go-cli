package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joaovrmoraes/http-go-cli/controller"
)

func main() {
	method, url, bearer, data, save, showInterface := parseFlags()

	if *showInterface {
		return
	}

	if *url == "" {
		fmt.Println("Please provide a URL using the -url flag")
		os.Exit(1)
	}

	if strings.HasPrefix(*url, ":") {
		*url = "http://localhost" + *url
	}

	controller.HandleRequest(*method, *url, *bearer, *data, *save)
}

func parseFlags() (*string, *string, *string, *string, *bool, *bool) {
	method := flag.String("method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	flag.StringVar(method, "m", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")

	url := flag.String("url", "", "URL to make the HTTPS request")

	bearer := flag.String("bearer", "", "Bearer token to authenticate the request")
	flag.StringVar(bearer, "b", "", "Bearer token to authenticate the request")

	data := flag.String("data", "", "Data to send in the request body (for POST, PUT, etc.)")
	flag.StringVar(data, "d", "", "Data to send in the request body (for POST, PUT, etc.)")

	save := flag.Bool("save", false, "Open the JSON response in the text editor")
	flag.BoolVar(save, "s", false, "Open the JSON response in the text editor")

	showInterface := flag.Bool("interface", false, "Show the Bubble Tea interface")
	flag.BoolVar(showInterface, "i", false, "Show the Bubble Tea interface")

	flag.Parse()
	return method, url, bearer, data, save, showInterface
}
