BINARY_NAME=httgo

build:
    go build -o $(BINARY_NAME) .

install:
    go install -o $(BINARY_NAME) .