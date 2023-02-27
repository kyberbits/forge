package forge

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/kyberbits/forge/forgeutils"
)

type HandlerContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func NewHandlerContext(w http.ResponseWriter, r *http.Request) *HandlerContext {
	return &HandlerContext{
		Writer:  w,
		Request: r,
	}
}

func (handlerContext *HandlerContext) ReadBody() []byte {
	bodyBytes, _ := io.ReadAll(handlerContext.Request.Body)
	handlerContext.Request.Body.Close()
	handlerContext.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

func (handlerContext *HandlerContext) DecodeRequest(target any) error {
	body := handlerContext.ReadBody()
	err := json.Unmarshal(body, target)
	if err != nil {
		if err.Error() == "EOF" {
			return errors.New("Request Body Can Not Be Blank")
		}
		return err
	}

	if validator, ok := target.(Validator); ok {
		if err := validator.Validate(); err != nil {
			if _, ok := err.(ValidatorError); ok {
				return err
			}

			return ValidatorError{Message: err.Error()}
		}
	}

	return nil
}

func (handlerContext *HandlerContext) RespondJSON(status int, v any) {
	encoder := json.NewEncoder(handlerContext.Writer)
	if jsonResponse, ok := v.(forgeutils.JsonResponse); ok {
		jsonResponse.ContextID = forgeutils.ContextGetID(handlerContext.Request.Context())
		v = jsonResponse
	}

	handlerContext.Writer.Header().Set("Content-Type", "application/json")
	handlerContext.Writer.WriteHeader(status)
	encoder.Encode(v)
}

func (handlerContext *HandlerContext) RespondHTML(status int, s string) {
	handlerContext.Writer.WriteHeader(status)
	handlerContext.Writer.Write([]byte(s))
}

func (handlerContext *HandlerContext) ExecuteTemplate(tmpl *template.Template, data any) {
	bodyBuffer := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(bodyBuffer, data); err != nil {
		http.Error(handlerContext.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := io.ReadAll(bodyBuffer)
	if err != nil {
		http.Error(handlerContext.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	handlerContext.RespondHTML(http.StatusOK, string(responseBytes))
}
