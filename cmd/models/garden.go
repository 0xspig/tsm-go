package models

import (
	"fmt"
	"os"
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
	newNode := new(Node)

	newNode.id = source
	newNode.data_source = source
	newNode.data_type = datatype

	garden.masterlist[newNode.id] = newNode
	garden.size += 1

	return newNode
}

func (garden *Garden) ConnectNodes(mainID string, outgoingID string) {
	garden.masterlist[mainID].outgoingNodes.AddNodeToList(garden.masterlist[outgoingID])
	garden.masterlist[outgoingID].incomingNodes.AddNodeToList(garden.masterlist[mainID])
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

func (garden *Garden) PopulateGardenFromDir(source_dir string) {

	// for each file in directory
	directory, err := os.ReadDir(source_dir)
	if err != nil {
		panic(err)
	}
	// check filetype (i'll do this later once we have multiple filetypes)
	// create nodes
	for _, file := range directory {
		fmt.Printf("Name:%s | Type: %s\n ", file.Name(), file.Type())
	}

}
