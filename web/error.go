package stormpathweb

import (
	"bytes"
	"fmt"
	"net/http"

	"net/url"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type errorModel struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func newErrorModel(spError stormpath.Error) errorModel {
	return errorModel{
		Status:  spError.Status,
		Message: spError.Message,
	}
}

func unauthorizedRequest(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	contentType := ctx.Value(ResolvedContentType)

	errorModel := buildErrorModelWithCode(fmt.Errorf("Unauthorized"), http.StatusUnauthorized)

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, errorModel, errorModel.Status)
		return
	}
	if contentType == stormpath.TextHTML {
		http.Redirect(w, r, Config.LoginURI+"?next="+getNextURI(r), http.StatusFound)
	}
}

func getNextURI(r *http.Request) string {
	buffer := bytes.Buffer{}

	buffer.WriteString(r.URL.Path)

	if r.URL.RawQuery != "" {
		buffer.WriteByte('?')
		buffer.WriteString(r.URL.RawQuery)
	}

	return url.QueryEscape(buffer.String())
}

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorModel := buildErrorModelWithCode(err, http.StatusBadRequest)

	respondJSON(w, errorModel, errorModel.Status)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	errorModel := buildErrorModelWithCode(fmt.Errorf("Method not allow"), http.StatusMethodNotAllowed)

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, errorModel, errorModel.Status)
		return
	}
	if contentType == stormpath.TextHTML {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func buildErrorModelWithCode(err error, status int) errorModel {
	model := errorModel{
		Status:  status,
		Message: err.Error(),
	}

	spError, ok := err.(stormpath.Error)
	if ok {
		model = newErrorModel(spError)
	}
	return model
}

func buildErrorModel(err error) errorModel {
	return buildErrorModelWithCode(err, http.StatusBadRequest)
}
