package models

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	g := CreateGarden()

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	home = filepath.Join(home, "prg/tsm/")
	os.Chdir(home)
	// TODO make content dir in config or something to search files in
	// for now Im just going to hack in static
	g.PopulateGardenFromDir("ui/content")
	g.ParseAllConnections()
	g.ConnectNodes("peepee.md", "firstFile.md")

	data, err := g.ExportJSONData()
	if err != nil {
		fmt.Printf("json export error :(")
	}

	fmt.Printf("JSON: %s\n", string(data))

	fmt.Printf("peepee\n")
}
