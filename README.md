# go-intigriti

Go library and commandline client for interacting with the [Intigriti](https://www.intigriti.com/) external API.

Checkout the autogenerated [SDK docs on pkg.go.dev](https://pkg.go.dev/github.com/hazcod/go-intigriti).

## Commandline client

Usage:
```shell
# list out all company programs
# also try: inti c list
% inti company list-programs

# list out all company submissions across all programs
# also try: inti c sub
% inti company list-submissions

# verify if a specific IP address is linked to an Intigriti user
# also try: inti c ip 1.1.1.1
% inti company check-ip 1.1.1.1
```

### Setup

Ensure the external API enabled on your company account and an integration is created with a redirect URI value of `http://localhost:1337/`.
Afterwards create the following local configuration file:

```yaml
log.level: info
auth:
    client_id: YOUR-CLIENT-ID
    client_secret: YOUR-CLIENT-SECRET
```

For the first call it will ask you to perform browser interaction to authenticate. <br/>
Future calls will not need to since your token will be cached in your configuration file.

## Library 

API Swagger documentation is available on the [ReadMe](https://intigriti.readme.io/reference/introduction).

### Usage
```go
package main

import (
	intigriti "github.com/hazcod/go-intigriti/pkg/api"
	"github.com/hazcod/go-intigriti/pkg/config"
	"log"
)

func main() {
	// this will require manual logon every time your code runs
	// look into persisting the TokenCache so refresh tokens can be saved
	// this will also launch an interactive Browser window to authenticate,
	// look at config.OpenBrowser and config.TokenCache to prevent this
	// or how the cli does it at https://github.com/hazcod/go-intigriti/blob/2eeb6a9fcee42fc4ac1ada7f5dc6d2db5446c15d/cmd/config/config.go#L86
	inti, err := intigriti.New(config.Config{
		Credentials: struct {
			ClientID     string
			ClientSecret string
		}{
		    ClientID: "my-integration-client-id",
		    ClientSecret: "my-integration-client-secret",
		},
	})
	if err != nil { log.Fatal(err) }
	
	programs, err := inti.GetPrograms()
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
