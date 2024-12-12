package models

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	g := CreateGarden()
	g.AddNodeToGarden(CONTENT_TYPE_MARKDOWN, "firstFile.md")
	g.AddNodeToGarden(CONTENT_TYPE_MARKDOWN, "secondFile.md")

	g.ConnectNodes("firstFile.md", "secondFile.md")
	// lets do it again and see if it dupes
	g.ConnectNodes("firstFile.md", "secondFile.md")
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	home = filepath.Join(home, "prg/tsm/")
	os.Chdir(home)
	g.PopulateGardenFromDir("ui/vite/content")

	fmt.Printf("peepee\n")
}
