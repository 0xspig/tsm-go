package models

import (
	"encoding/json"
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
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Data_source         string `json:"source"`
	data_type           int
	numberIncomingNodes int
	numberOutgoingNodes int
	incomingNodes       NodeList
	outgoingNodes       NodeList
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

	newNode.ID = source
	newNode.Data_source = source
	newNode.data_type = datatype
	newNode.numberIncomingNodes = 0
	newNode.numberOutgoingNodes = 0

	garden.masterlist[newNode.ID] = newNode
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

	if current.node.ID == nodeToAdd.ID {
		return
	}

	for buffer != nil {
		if current.node.ID == nodeToAdd.ID {
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
	master := garden.masterlist[mainID]
	outgoing := garden.masterlist[outgoingID]
	//verify that both IDs exist
	err := 0
	if master == nil {
		fmt.Printf("Error: nil node ID - %s: %p\n", mainID, master)
		err = 1
	}
	if outgoing == nil {
		fmt.Printf("Error: nil node ID - %s: %p\n", outgoingID, outgoing)
		err = 1
	}
	if err == 1 {
		return
	}
	master.outgoingNodes.AddNodeToList(outgoing)
	master.numberOutgoingNodes += 1
	outgoing.incomingNodes.AddNodeToList(master)
	outgoing.numberIncomingNodes += 1
}

// Parses all node sources and populates outgoing and respective incoming connections
func (garden *Garden) ParseAllConnections() {
	os.Chdir("C:/Users/tmcke/prg/tsm/ui/content")
	for _, node := range garden.masterlist {
		data, err := os.ReadFile(node.Data_source)
		if err != nil {
			panic(err)
		}
		links := findLinks(data)

		for _, link := range links {
			// link[2] is should be the src in the regex function. if this breaks check the regex
			garden.ConnectNodes(node.ID, link[2])
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

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type GraphData struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func (garden *Garden) ExportJSONData() ([]byte, error) {
	data := GraphData{}
	for _, node := range garden.masterlist {
		data.Nodes = append(data.Nodes, *node)
		for link := node.incomingNodes; link.node != nil; {
			newLink := Link{Source: node.ID, Target: link.node.ID}
			data.Links = append(data.Links, newLink)
			if link.next != nil {
				link = *link.next
			} else {
				break
			}
		}
	}
	return json.Marshal(data)
}
