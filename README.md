# HTTP Go CLI

HTTP Go CLI is a command-line tool written in Golang that allows you to make HTTP requests directly from the terminal. This tool mimics the behavior of HTTPie, offering a simple and intuitive interface to interact with APIs.

## Features

- Perform HTTP requests (GET, POST, PUT, DELETE, etc.)
- Support for query parameters
- Option to save HTTP responses to files

## Usage

To use HTTP Go CLI, simply run the executable followed by the desired HTTP method and URL. For example:

Perform a GET request

.`/http-go-cli -m GET -url "htttps://api.example.com/resources`

Perform a POST request with JSON data

```
`./http-go-cli -m POST -url https://api.example.com/resources -d '{"key": "value"}'
```
