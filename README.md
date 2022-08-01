# go-intigriti
Go library and commandline client for interacting with the [intigriti](https://www.intigriti.com/) v1 and v2 API.
Checkout the [docs](https://pkg.go.dev/github.com/intigriti/sdk-go)!


## Library 

API documentation is available on the [ReadMe](https://dash.readme.com/project/intigriti/v2.0/overview).

### Usage
```go
package main

import (
	inti "github.com/hazcod/go-intigriti/v2"
	"log"
)

func main() {
	intigriti, err := inti.New("my-client-token", "my-client-secret", nil) // or a logrus.Logger
	if err != nil { log.Fatal(err) }
	
	programs, err := intigriti.GetPrograms()
	if err != nil { log.Fatal(err) }

	for _, program := range programs {
		log.Println(program.Name)
	}
}
```

### Testing
```shell script
# test on production using inti.yml
go test -tags integration -v ./...

# test on staging using inti.yml
INTI_TOKEN_URL=="testing.api.com" INTI_AUTH_URL=="subs.testing.api.com" INTI_API_URL="api.testing.com" go test -tags integration -v ./...
```

## Commandline client

Usage:
```shell
% inti company list-programs
% inti company list-submissions
```