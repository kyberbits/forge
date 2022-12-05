package forge

import "errors"

const (
	HeaderContentType           = "Content-Type"
	HeaderCacheControl          = "Cache-Control"
	HeaderContentSecurityPolicy = "Content-Security-Policy"
)

const (
	LevelCritical = "CRITICAL"
	LevelError    = "ERROR"
	LevelWarning  = "WARNING"
	LevelInfo     = "INFO"
	LevelDebug    = "DEBUG"
)

var (
	ErrInvalidValue         = errors.New("value must be a non-nil pointer to a struct")
	ErrUnsupportedFieldType = errors.New("field is an unsupported type")
	ErrUnexportedField      = errors.New("field must be exported")
	ErrNotAuthenticated     = errors.New("not authenticated")
)
