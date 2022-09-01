package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	htmlAutoClose = `
	<html>
		<body><h1>Successfully authenticated, you may close this page.</h1></body>
		<script>window.close();</script>
	</html>
	`
)

type callbackResult struct {
	Error error
	Code  string
}

// this is the local http listener which will be called after successfully authenticating to Intigriti
// here we will compare the state parameter to prevent csrf and extract the authorization code
func (e *Endpoint) getLocalHandler(uri, state string, resultChan chan callbackResult, doneChan chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != uri {
			e.logger.WithField("path", r.URL.Path).Debug("invalid callback path")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.URL.Query().Get("state") != state {
			e.logger.WithField("required_state", state).WithField("given_state", r.URL.Query().Get("state")).
				Warn("invalid state provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/html; charset=UTF-8")
		_, _ = w.Write([]byte(htmlAutoClose))

		resultChan <- callbackResult{
			Code:  r.URL.Query().Get("code"),
			Error: nil,
		}

		doneChan <- struct{}{}

		e.logger.WithField("code", r.URL.Query().Get("code")).Debug("callback successfully got code")
	})
}

// helper function that creates the callback listener and waits until a response is received or timeout expires
func (e *Endpoint) listenForCallback(uri, localHost string, localPort uint, state string, resultChan chan callbackResult) {
	e.logger.WithField("port", localPort).Debug("listening for callback for new authorization code")

	doneChan := make(chan struct{}, 2)

	srv := http.Server{}
	srv.Addr = fmt.Sprintf("%s:%d", localHost, localPort)
	srv.Handler = e.getLocalHandler(uri, state, resultChan, doneChan)

	go func() {
		<-doneChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		_ = srv.Shutdown(ctx)
		e.logger.Debug("shut down local callback listener")
		cancel() // just to fix govet
	}()

	err := srv.ListenAndServe()
	resultChan <- callbackResult{Error: err}
	e.logger.Debug("returning from listenForCallback")
}
