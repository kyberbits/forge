package forgeutils

import (
	"io"
	"os"
	"regexp"
)

func EnvironmentSetValueInFile(filePath, key string, value string) {
	targetLine := key + "=" + value

	originalFileString := ""

	file, err := os.Open(filePath)
	if err == nil {
		fileBytes, _ := io.ReadAll(file)
		file.Close()

		originalFileString = string(fileBytes)
	}

	re := regexp.MustCompile(`(?m)^` + key + `=.*$`)
	matchLocation := re.FindStringIndex(originalFileString)

	newFileContents := ""
	if (len(matchLocation)) <= 0 {
		newFileContents += targetLine + "\n"
		newFileContents += originalFileString
	} else {
		newFileContents = originalFileString[0:matchLocation[0]]
		newFileContents += targetLine
		newFileContents += originalFileString[matchLocation[1]:]
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(newFileContents); err != nil {
		panic(err)
	}
}
