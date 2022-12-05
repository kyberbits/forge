package forge_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/kyberbits/forge/forge"
)

func TestResources(t *testing.T) {
	resources := forge.NewResources([]fs.FS{
		os.DirFS("./test_files/resources"),
	})
	resources.MustParseHTMLTemplate("test.go.tmpl")
	resources.MustOpenFileContents("txt/foo.txt")
	resources.MustOpenDirectory("txt")
}
