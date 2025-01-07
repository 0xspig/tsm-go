package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/yuin/goldmark"
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
	CONTENT_TYPE_TAG      = 2
	CONTENT_TYPE_CATEGORY = 3
)

// Essential node element/
type Node struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Data_source         string `json:"source"`
	Data_type           int    `json:"data_type"`
	NumberIncomingNodes int    `json:"numIncoming"`
	NumberOutgoingNodes int    `json:"numOutgoing"`
	incomingNodes       NodeList
	outgoingNodes       NodeList
}

func CreateGarden() *Garden {
	return &Garden{
		masterlist: make(map[string]*Node),
		size:       0,
	}
}

// adds node to garden
func (garden *Garden) addNodeToGarden(datatype int, source string, id string, name string) *Node {
	if garden.masterlist[id] != nil {
		fmt.Printf("Node source already exists\n")
		return garden.masterlist[id]
	}
	newNode := new(Node)

	newNode.ID = id
	newNode.Data_source = source
	newNode.Data_type = datatype
	newNode.Name = name
	newNode.NumberIncomingNodes = 0
	newNode.NumberOutgoingNodes = 0

	garden.masterlist[newNode.ID] = newNode
	garden.size += 1

	return newNode

}

// adds node to garden. Source should be filepath relative to root.
func (garden *Garden) AddSourceToGarden(datatype int, source string) *Node {
	if garden.masterlist[source] != nil {
		fmt.Printf("Node source already exists\n")
		return garden.masterlist[source]
	}
	newNode := new(Node)

	newNode.ID = filepath.Base(source)
	newNode.Data_source = source
	newNode.Data_type = datatype
	if datatype == CONTENT_TYPE_MARKDOWN {
		data, err := os.ReadFile(source)
		if err != nil {
			panic(err)
		}
		newNode.Name = scanYAMLFrontMatter(data).Title
	}
	newNode.NumberIncomingNodes = 0
	newNode.NumberOutgoingNodes = 0

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
		if file.IsDir() {
			garden.addNodeToGarden(CONTENT_TYPE_CATEGORY, filepath.Clean(file.Name()), file.Name(), file.Name())
			garden.PopulateGardenFromDir(filepath.Join(source_dir, file.Name()))
		} else {
			relLink := filepath.Clean(filepath.Join(source_dir, file.Name()))
			if err != nil {
				panic(err)
			}
			garden.AddSourceToGarden(CONTENT_TYPE_MARKDOWN, relLink)
		}
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
	master.NumberOutgoingNodes += 1
	outgoing.incomingNodes.AddNodeToList(master)
	outgoing.NumberIncomingNodes += 1
}

// Parses all node sources and populates outgoing and respective incoming connections
func (garden *Garden) ParseAllConnections() {
	for _, node := range garden.masterlist {

		// if datatype is md or html link to parent category
		if node.Data_type < CONTENT_TYPE_TAG {
			data, err := os.ReadFile(node.Data_source)
			if err != nil {

			}
			fileLinks, tagLinks := garden.findLinks(data)

			for _, link := range fileLinks {
				// link[2] is should be the src in the regex function. if this breaks check the regex
				garden.ConnectNodes(node.ID, filepath.Base(link))
			}
			for _, link := range tagLinks {
				// link[2] is should be the src in the regex function. if this breaks check the regex
				garden.ConnectNodes(node.ID, link)
			}

			// parent node (category) directs to child (post)
			category_id, err := filepath.Rel("ui/content", filepath.Dir(node.Data_source))
			if err != nil {
				panic(err)
			}
			garden.ConnectNodes(category_id, node.ID)
		}
	}

}

type YAMLData struct {
	Title    string
	Date     string
	Category string
	Tags     []string
}

// scans file for yaml frontmatter between '---' separators
func scanYAMLFrontMatter(data []byte) *YAMLData {
	var frontMatter YAMLData
	var YAMLBytes []byte

	scanner := bufio.NewScanner(bytes.NewReader(data))
	breakCount := 0
	for breakCount < 2 {
		if !scanner.Scan() {
			break
		}
		YAMLBytes = append(YAMLBytes, scanner.Bytes()...)
		YAMLBytes = append(YAMLBytes, "\n"...)
		if scanner.Text() == "---" {
			breakCount++
		}
	}
	err := yaml.Unmarshal(YAMLBytes, &frontMatter)
	if err != nil {
		panic(err)
	}
	return &frontMatter
}

// parse markdown files for links
func (garden *Garden) findLinks(data []byte) ([]string, []string) {

	frontMatter := scanYAMLFrontMatter(data)

	tagMatches := make([]string, 0)

	for _, tag := range frontMatter.Tags {
		if garden.masterlist[tag] == nil {
			garden.addNodeToGarden(CONTENT_TYPE_TAG, "index.md", tag, "tag:"+tag)
		}
		tagMatches = append(tagMatches, tag)
	}

	// this gets the link value and source '[<value>](<src>)'
	regular_expression, err := regexp.Compile(`\[([^\]]*)\]\(([^\)]*)\)`)

	if err != nil {
		panic(err)
	}
	// substring returns 3 strings for each match 0:full match 1:value 2:src
	matches := regular_expression.FindAllStringSubmatch(string(data), -1)
	matchValues := make([]string, 0)
	for _, match := range matches {
		matchValues = append(matchValues, match[2])
	}

	return matchValues, tagMatches
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

func (garden *Garden) NodeToHTML(nodeID string) []byte {
	node := garden.masterlist[nodeID]
	if node == nil {
		return []byte("<h1>File Not Found</h1>")
	}
	switch node.Data_type {
	case CONTENT_TYPE_MARKDOWN:
		return garden.mdToHTML(node)
	case CONTENT_TYPE_HTML:
		return []byte("<h1>HTML unsupported</h1>")
	case CONTENT_TYPE_TAG:
		return []byte("<h1>Tags currently unsupported</h1>")
	case CONTENT_TYPE_CATEGORY:
		return []byte("<h1>Categories currently unsupported</h1>")
	default:
		return []byte("<h1>File Not Found</h1>")
	}
}
func (garden *Garden) mdToHTML(node *Node) []byte {
	source, err := os.ReadFile(node.Data_source)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
