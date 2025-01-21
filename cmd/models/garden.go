package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"maps"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type StringSet map[string]bool

// Struct containing a hash table of all nodes in graph
type Garden struct {
	Masterlist map[string]*Node
	idList     []string
	size       int
	Tags       StringSet
	Categories StringSet
	Center     string
}

const (
	CONTENT_TYPE_HTML     = 0
	CONTENT_TYPE_MARKDOWN = 1
	CONTENT_TYPE_TAG      = 2
	CONTENT_TYPE_CATEGORY = 3
	CONTENT_TYPE_EXTERNAL = 4
)

type NodeSet map[*Node]bool

// Essential node element/
type Node struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Data_source         string   `json:"source"`
	Data_type           int      `json:"data_type"`
	NumberIncomingNodes int      `json:"numIncoming"`
	NumberOutgoingNodes int      `json:"numOutgoing"`
	IncomingNodes       NodeSet  `json:"-"`
	OutgoingNodes       NodeSet  `json:"-"`
	Metadata            YAMLData `json:"-"`
	ParentGarden        *Garden  `json:"-"`
}

type NodeList []Node

func CreateGarden() *Garden {
	return &Garden{
		Masterlist: make(map[string]*Node),
		size:       0,
		Tags:       make(StringSet),
		Categories: make(StringSet),
	}
}
func (garden *Garden) ContainsID(id string) bool {
	return garden.Masterlist[id] != nil
}

// adds node to garden
func (garden *Garden) addNodeToGarden(datatype int, source string, id string, name string) *Node {
	if garden.Masterlist[id] != nil {
		fmt.Printf("Node source already exists\n")
		return garden.Masterlist[id]
	}
	newNode := new(Node)

	newNode.ID = id
	newNode.Data_source = source
	newNode.Data_type = datatype
	newNode.Name = name
	newNode.OutgoingNodes = NodeSet{}
	newNode.IncomingNodes = NodeSet{}
	newNode.NumberIncomingNodes = 0
	newNode.NumberOutgoingNodes = 0
	newNode.ParentGarden = garden
	switch newNode.Data_type {
	default:
		break
	case CONTENT_TYPE_MARKDOWN:
		data, err := os.ReadFile(source)
		if err != nil {
			panic(err)
		}
		yaml := scanYAMLFrontMatter(data)
		newNode.Metadata = *yaml
	case CONTENT_TYPE_TAG:
		if garden.Tags[id] {
			break
		} else {
			garden.Tags[id] = true
		}
	case CONTENT_TYPE_CATEGORY:
		if garden.Categories[id] {
			break
		} else {
			garden.Categories[id] = true
		}
	}
	garden.Masterlist[newNode.ID] = newNode
	garden.idList = append(garden.idList, id)
	garden.size += 1

	return newNode

}

// adds node to garden. Source should be filepath relative to root.
func (garden *Garden) AddSourceToGarden(datatype int, source string) *Node {
	if garden.Masterlist[source] != nil {
		fmt.Printf("Node source already exists\n")
		return garden.Masterlist[source]
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
		yaml := scanYAMLFrontMatter(data)
		newNode.Name = yaml.Title
		newNode.Metadata = *yaml
	}
	newNode.NumberIncomingNodes = 0
	newNode.NumberOutgoingNodes = 0
	newNode.ParentGarden = garden

	newNode.OutgoingNodes = NodeSet{}
	newNode.IncomingNodes = NodeSet{}

	garden.Masterlist[newNode.ID] = newNode
	garden.idList = append(garden.idList, newNode.ID)
	garden.size += 1

	return newNode
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
	master := garden.Masterlist[mainID]
	outgoing := garden.Masterlist[outgoingID]
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

	maps.Insert(master.OutgoingNodes, maps.All(NodeSet{outgoing: true}))
	master.NumberOutgoingNodes += 1

	maps.Insert(outgoing.IncomingNodes, maps.All(NodeSet{master: true}))
	outgoing.NumberIncomingNodes += 1
}

// Parses all node sources and populates outgoing and respective incoming connections
func (garden *Garden) ParseAllConnections() {
	for _, node := range garden.Masterlist {

		// if datatype is md or html link to parent category
		if node.Data_type < CONTENT_TYPE_TAG {
			data, err := os.ReadFile(node.Data_source)
			if err != nil {

			}
			baseLinks, fullLinks := garden.findLinks(data)

			for _, link := range baseLinks {
				// link[2] is should be the src in the regex function. if this breaks check the regex
				garden.ConnectNodes(node.ID, filepath.Base(link))
			}
			for _, link := range fullLinks {
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
	garden.Center = garden.findCenter()[0]
}

type YAMLData struct {
	Title    string
	Date     string
	Category string
	Tags     []string
	Class    string
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
		if garden.Masterlist[tag] == nil {
			// TODO fix source - currently just placeholder index.md
			garden.addNodeToGarden(CONTENT_TYPE_TAG, "index.md", tag, "Tag: "+tag)
		}
		tagMatches = append(tagMatches, tag)
	}

	// this gets the link value and source '{<value>](<src>)'
	internal_regex, err := regexp.Compile(`\{([^\}]*)\}\(([^\)]*)\)`)

	if err != nil {
		panic(err)
	}
	// substring returns 3 strings for each match 0:full match 1:value 2:src
	matches := internal_regex.FindAllStringSubmatch(string(data), -1)
	matchValues := make([]string, 0)
	for _, match := range matches {
		matchValues = append(matchValues, match[2])
	}

	//same as above for external links
	external_regex, err := regexp.Compile(`\[([^\]]*)\]\(([^\)]*)\)`)
	if err != nil {
		panic(err)
	}
	matches = external_regex.FindAllStringSubmatch(string(data), -1)

	// regex stolen from Berners-Lee
	//$0 = http://www.ics.uci.edu/pub/ietf/uri/#Related
	//$1 = http:
	//$2 = http
	//$3 = //www.ics.uci.edu
	//$4 = www.ics.uci.edu
	//$5 = /pub/ietf/uri/
	//$6 = <undefined>
	//$7 = <undefined>
	//$8 = #Related
	//$9 = Related
	uri_regex, err := regexp.Compile(`^(([^:\/?#]+):)?(\/\/([^\/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		uri := uri_regex.FindStringSubmatch(match[2])
		garden.addNodeToGarden(CONTENT_TYPE_EXTERNAL, uri[0], uri[0], match[1])
		// matches are added to tag matches because the internal file matches get truncated later
		tagMatches = append(tagMatches, uri[0])
	}

	return matchValues, tagMatches
}

func (garden *Garden) shortestPath(origin string, dest string) int {
	if _, exists := garden.Masterlist[dest]; !exists {
		fmt.Printf("Unable to find path of nonexistent node: %s\n", dest)
		return -1
	}
	if _, exists := garden.Masterlist[origin]; !exists {
		fmt.Printf("Unable to find path of nonexistent node: %s\n", origin)
		return -1
	}

	origin_connections := maps.Clone(garden.Masterlist[origin].IncomingNodes)
	maps.Copy(origin_connections, garden.Masterlist[origin].OutgoingNodes)

	if _, exists := origin_connections[garden.Masterlist[dest]]; exists {
		return 1
	}

	stepsMap := map[string]int{origin: 0}
	numSteps := 1

	node_list := origin_connections
	for len(stepsMap) < garden.size {
		for node := range node_list {
			if _, exists := stepsMap[node.ID]; !exists {
				maps.Insert(stepsMap, maps.All(map[string]int{node.ID: numSteps}))
			}
		}

		swap_list := maps.Clone(node_list)
		clear(node_list)
		for c := range swap_list {
			maps.Copy(node_list, c.IncomingNodes)
			maps.Copy(node_list, c.OutgoingNodes)
		}
		maps.DeleteFunc(node_list, func(key *Node, value bool) bool {
			_, exists := stepsMap[key.ID]
			return exists
		})

		numSteps++
	}

	return stepsMap[dest]
}
func (garden *Garden) findCenter() []string {
	minDists := make([][]int, garden.size)
	for i := range garden.idList {
		minDists[i] = make([]int, garden.size)
	}
	for i := range garden.idList {
		for j := range garden.idList {
			if i == j {
				minDists[i][j] = 0
				continue
			}
			_, itoj := garden.Masterlist[garden.idList[i]].IncomingNodes[garden.Masterlist[garden.idList[j]]]
			_, jtoi := garden.Masterlist[garden.idList[j]].IncomingNodes[garden.Masterlist[garden.idList[i]]]
			if itoj || jtoi {
				minDists[i][j] = 1
			} else {
				minDists[i][j] = math.MaxInt
			}
			//minDists[i][j] = garden.shortestPath(garden.idList[i], garden.idList[j])
		}
	}

	for k := range garden.idList {
		for i := range garden.idList {
			for j := range garden.idList {
				if minDists[i][j] > minDists[i][k]+minDists[k][j] {
					//stop integer overflowing
					if minDists[i][k] == math.MaxInt || minDists[k][j] == math.MaxInt {
						continue
					}
					minDists[i][j] = minDists[i][k] + minDists[k][j]
				}
			}
		}

	}

	indexofhip := slices.Index(garden.idList, "hips.md")
	indexofhoney := slices.Index(garden.idList, "honey-hollow.md")

	hiptohoney := minDists[indexofhip][indexofhoney]
	fmt.Printf("distance between hips(%d) and honey(%d) 2: %d\n", indexofhip, indexofhoney, hiptohoney)

	e := slices.Repeat([]int{0}, garden.size)
	for i := range garden.idList {
		for j := range garden.idList {
			if minDists[i][j] > e[i] {
				e[i] = minDists[i][j]
			}
		}
	}
	rad := math.MaxInt
	diam := 0
	for i := range garden.idList {
		if rad > e[i] {
			rad = e[i]
		}
		if diam < e[i] {
			diam = e[i]
		}
	}

	var center []string
	for i := range garden.idList {
		if e[i] == rad {
			center = append(center, garden.idList[i])
		}
	}
	return center
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type GraphData struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func (garden *Garden) genJSONData() ([]byte, error) {
	var data GraphData
	for _, node := range garden.Masterlist {
		data.Nodes = append(data.Nodes, *node)
		for addr, _ := range node.IncomingNodes {
			newLink := Link{Source: node.ID, Target: addr.ID}
			data.Links = append(data.Links, newLink)
		}
	}
	return json.Marshal(data)
}

func (garden *Garden) NodeToHTML(nodeID string) []byte {
	node := garden.Masterlist[nodeID]
	if node == nil {
		return []byte("<h1>File Not Found</h1>")
	}
	switch node.Data_type {
	case CONTENT_TYPE_MARKDOWN:
		return garden.mdToHTML(node)
	case CONTENT_TYPE_HTML:
		return []byte("<h1>HTML unsupported</h1>")
	case CONTENT_TYPE_TAG:
		return garden.tagToHtml(node)
	case CONTENT_TYPE_CATEGORY:
		return garden.catToHtml(node)
	default:
		return []byte("<h1>File Not Found</h1>")
	}
}

func (garden *Garden) NodeLinksToHTML(nodeID string) []byte {
	ts, err := template.ParseFiles("ui/templates/links.template.html")
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = ts.ExecuteTemplate(&buf, "links", garden.Masterlist[nodeID])
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (garden *Garden) mdToHTML(node *Node) []byte {
	source, err := os.ReadFile(node.Data_source)
	if err != nil {
		panic(err)
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	internal_regex, err := regexp.Compile(`\{([^\}]*)\}\(([^\)]*)\)`)
	if err != nil {
		panic(err)
	}

	data := internal_regex.ReplaceAll(buf.Bytes(), []byte(`<a class="internal-link" onClick="targetNode('$2')">$1</a>`))

	data = append([]byte(`{{define "content"}}`), data...)
	data = append(data, []byte(`{{end}}`)...)

	var template_file string

	switch node.Metadata.Class {
	default:
		template_file = "ui/templates/post.template.html"
	case "home":
		template_file = "ui/templates/home.template.html"
	}

	ts, err := template.ParseFiles(template_file, "ui/templates/footer.template.html")
	if err != nil {
		panic(err)
	}
	ts, err = ts.Parse(string(data))
	if err != nil {
		panic(err)
	}
	var template_buf bytes.Buffer
	ts.Execute(&template_buf, node)

	return template_buf.Bytes()
}

func (garden *Garden) tagToHtml(node *Node) []byte {
	ts, err := template.ParseFiles("ui/templates/tag.template.html", "ui/templates/footer.template.html")
	if err != nil {
		return []byte("<h1>Template rendering error</h1>")
	}

	var buf bytes.Buffer
	ts.Execute(&buf, node)

	return buf.Bytes()
}

func (garden *Garden) catToHtml(node *Node) []byte {
	ts, err := template.ParseFiles("ui/templates/cat.template.html", "ui/templates/footer.template.html")
	if err != nil {
		return []byte("<h1>Template rendering error</h1>")
	}

	var buf bytes.Buffer
	ts.Execute(&buf, node)

	return buf.Bytes()

}

func (garden *Garden) GenAssets() {
	// cache node data
	json_data, err := garden.genJSONData()
	if err != nil {
		panic(err)
	}
	os.WriteFile("ui/static/gen/graph-data.json", json_data, 0644)
}
