package forgeutils

type JSONResponse struct {
	ContextID string `json:"context_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}
