package typescript

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Interfaces []Interface

func (i Interfaces) CreateFile(targetFile string) error {
	fileBytes := []byte{}
	for _, typescriptInterface := range i {
		fileBytes = append(fileBytes, []byte(typescriptInterface.String())...)
		fileBytes = append(fileBytes, []byte("\n")...)
		fileBytes = append(fileBytes, []byte("\n")...)
	}

	fileBytes = bytes.TrimSpace(fileBytes)
	fileBytes = append(fileBytes, []byte("\n")...)

	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(fileBytes); err != nil {
		return err
	}

	return nil
}

type Interface struct {
	Name   string
	Fields []Field
}

func (i Interface) String() string {
	if len(i.Fields) < 1 {
		return fmt.Sprintf(
			"export interface %s {}",
			i.Name,
		)
	}

	fieldStrings := []string{}
	for _, field := range i.Fields {
		fieldStrings = append(fieldStrings, field.String())
	}

	return fmt.Sprintf(
		"export interface %s {\n\t%s\n}",
		i.Name,
		strings.Join(fieldStrings, "\n\t"),
	)
}
