//Back to routes.go,
//we did not just return router but
//app.recoverPanic(router). What is
//recoverPanic? You asked. It's a
//middleware that does what its name
//suggests: gracefully recover from panic!
//We don't want to shut out our users without
//letting them know there was an issue.

package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
