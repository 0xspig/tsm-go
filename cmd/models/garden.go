package models

import (
	"fmt"
	"os"
	"regexp"
)

// Struct containing a hash table of all nodes in graph
type Garden struct {
	masterlist map[string]*Node
	size       int
}

// Double linked list of nodes
type NodeList struct {
	node *Node
	next *NodeList
	prev *NodeList
}

const (
	CONTENT_TYPE_HTML     = 0
	CONTENT_TYPE_MARKDOWN = 1
)

// Essential node element
type Node struct {
	id            string
	incomingNodes NodeList
	outgoingNodes NodeList

	data_type   int
	data_source string
}

func CreateGarden() *Garden {
	return &Garden{
		masterlist: make(map[string]*Node),
		size:       0,
	}
}

func (garden *Garden) AddNodeToGarden(datatype int, source string) *Node {
	if garden.masterlist[source] != nil {
		fmt.Printf("Node source already exists\n")
		return garden.masterlist[source]
	}
	if garden.masterlist[source] != nil {
		fmt.Printf("Node source already exists\n")
		return garden.masterlist[source]
	}
	newNode := new(Node)

	newNode.id = source
	newNode.data_source = source
	newNode.data_type = datatype

	garden.masterlist[newNode.id] = newNode
	garden.size += 1

	return newNode
}

func (list *NodeList) AddNodeToList(nodeToAdd *Node) {
	if list.node == nil {
		*list = NodeList{node: nodeToAdd, next: nil, prev: nil}
		return
	}
	current := list
	buffer := list.next

	if current.node.id == nodeToAdd.id {
		return
	}

	for buffer != nil {
		if current.node.id == nodeToAdd.id {
			return
		}
		current = current.next
		buffer = current.next
	}

	current.next = &NodeList{node: nodeToAdd, next: nil, prev: current}
}

/*TODO func checkFileType(file) int*/

// Populates garden with nodes generated from source_dir (note: nodes will remain islands until connected)
// Populates garden with nodes generated from source_dir (note: nodes will remain islands until connected)
func (garden *Garden) PopulateGardenFromDir(source_dir string) {

	// for each file in directory
	directory, err := os.ReadDir(source_dir)
	if err != nil {
		panic(err)
	}
	// create nodes
	for _, file := range directory {
		fmt.Printf("Name:%s | Type: %s\n ", file.Name(), file.Type())

		// check filetype (i'll do this later once we have multiple filetypes) asuming md for now
		garden.AddNodeToGarden(CONTENT_TYPE_MARKDOWN, file.Name())
	}
}

// Connect two nodes so that mainID node directs to outgoingID node.
func (garden *Garden) ConnectNodes(mainID string, outgoingID string) {
	garden.masterlist[mainID].outgoingNodes.AddNodeToList(garden.masterlist[outgoingID])
	garden.masterlist[outgoingID].incomingNodes.AddNodeToList(garden.masterlist[mainID])
}

// Parses all node sources and populates outgoing and respective incoming connections
func (garden *Garden) ParseAllConnections() {
	os.Chdir("C:/Users/tmcke/prg/tsm/ui/vite/content")
	for _, node := range garden.masterlist {
		data, err := os.ReadFile(node.data_source)
		if err != nil {
			panic(err)
		}
		links := findLinks(data)

		for _, link := range links {
			// link[2] is should be the src in the regex function. if this breaks check the regex
			garden.ConnectNodes(node.id, link[2])
		}
	}

}

// parse markdown files for links
func findLinks(data []byte) [][]string {
	// this gets the link value and source '[<value>](<src>)'
	regular_expression, err := regexp.Compile(`\[([^\]]*)\]\(([^\)]*)\)`)

	if err != nil {
		panic(err)
	}
	// substring returns 3 strings for each match 0:full match 1:value 2:src
	matches := regular_expression.FindAllStringSubmatch(string(data), -1)
	return matches
}
