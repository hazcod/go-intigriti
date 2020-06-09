# go-intigriti
Go library to interact with the intigriti API.

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
