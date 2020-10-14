# go-intigriti
Go library for interacting with the [intigriti](https://www.intigriti.com/) API.
Checkout the [docs](https://pkg.go.dev/github.com/hazcod/go-intigriti)!

## Usage
```go
package main

import (
	inti "github.com/hazcod/go-intigriti"
	"log"
)

func main() {
	intigriti := inti.New("my-client-token", "my-client-secret")
	submissions, err := intigriti.GetSubmissions()
	if err != nil { log.Fatal(err) }

	for _, submission := range submissions {
		log.Println(submission.Title)
	}
}
```

## Testing
```shell script
# test on production
TOKEN="foo" SECRET="bar" go test -tags integration -v ./...

# test on staging
TOKEN="foo" SECRET="bar" AUTH_API="testing.api.com" SUB_API="subs.testing.api.com" go test -tags integration -v ./...
```