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

func (e *Endpoint) getLocalHandler(state string, resultChan chan callbackResult, doneChan chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			e.Logger.WithField("path", r.URL.Path).Debug("invalid callback path")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.URL.Query().Get("state") != state {
			e.Logger.WithField("required_state", state).WithField("given_state", r.URL.Query().Get("state")).
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

		e.Logger.WithField("code", r.URL.Query().Get("code")).Debug("callback successfully got code")
	})
}

func (e *Endpoint) listenForCallback(localPort uint, state string, resultChan chan callbackResult) {
	e.Logger.WithField("port", localPort).Debug("listening for callback for new authorization code")

	doneChan := make(chan struct{}, 2)

	srv := http.Server{}
	srv.Addr = fmt.Sprintf("localhost:%d", localPort)
	srv.Handler = e.getLocalHandler(state, resultChan, doneChan)

	go func() {
		<-doneChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		_ = srv.Shutdown(ctx)
		e.Logger.Debug("shut down local callback listener")
		cancel() // just to fix govet
	}()

	err := srv.ListenAndServe()
	resultChan <- callbackResult{Error: err}
	e.Logger.Debug("returning from listenForCallback")
}
