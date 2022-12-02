package forge_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/kyberbits/forge"
)

func TestLoggerJSON_Emergency(t *testing.T) {
	buffer := bytes.NewBufferString("")
	logger := forge.LoggerJSON{
		Encoder: json.NewEncoder(buffer),
	}
	logger.StandardLogger().Print("test")
	logger.Critical("test Critical", nil)
	logger.Error("test Error", nil)
	logger.Warning("test Warning", nil)
	logger.Info("test Info", nil)
	logger.Debug("test Debug", nil)
}
