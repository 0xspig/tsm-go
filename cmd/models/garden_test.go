package models

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	g := CreateGarden()
	g.AddNodeToGarden(CONTENT_TYPE_MARKDOWN, "firstFile.md")
	g.AddNodeToGarden(CONTENT_TYPE_MARKDOWN, "secondFile.md")

	g.ConnectNodes("firstFile.md", "secondFile.md")

	fmt.Printf("peepee\n")
}
