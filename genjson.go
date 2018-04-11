package genjson

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var intRe = regexp.MustCompile(`\d+`)
var floatRe = regexp.MustCompile(`\d+\.\d+`)
var boolRe = regexp.MustCompile(`(true|false)`)
var queryRe = regexp.MustCompile(`[^\.^\[^\]]+`)

type Type int

const (
	Unknown Type = iota
	Map
	Slice
	String
	Int
	Float
	Bool
	// Null
)

type Node struct {
	childrenMap map[interface{}]*Node
	children    []*Node
	content     string
	typ         Type
	isNull      bool
}

func NewNode() *Node {
	return &Node{
		children:    make([]*Node, 0, 64),
		childrenMap: make(map[interface{}]*Node),
	}
}

func (n *Node) readMap(reader *strings.Reader) {
	n.typ = Map
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return
		}
		s := string(r)
		switch s {
		case "{":
			child := NewNode()
			n.children = append(n.children, child)
			child.readMap(reader)
		case "}":
			return
		case "[":
			child := NewNode()
			n.children = append(n.children, child)
			child.readSlice(reader)
		case ",", " ", "\n", "\r", "\t":
			continue
		case "\"":
			child := NewNode()
			n.children = append(n.children, child)
			child.readString(reader)
		case ":":
			child := NewNode()
			child.content = ":"
			n.children = append(n.children, child)
		default:
			reader.UnreadRune()
			child := NewNode()
			n.children = append(n.children, child)
			child.readOther(reader)
		}
	}
}

func (n *Node) readSlice(reader *strings.Reader) {
	n.typ = Slice
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return
		}
		s := string(r)
		switch s {
		case "{":
			child := NewNode()
			n.children = append(n.children, child)
			child.readMap(reader)
		case "[":
			child := NewNode()
			n.children = append(n.children, child)
			child.readSlice(reader)
		case "]":
			return
		case ",", " ", "\n", "\r", "\t":
			continue
		case "\"":
			child := NewNode()
			n.children = append(n.children, child)
			child.readString(reader)
		default:
			reader.UnreadRune()
			child := NewNode()
			n.children = append(n.children, child)
			child.readOther(reader)
		}
	}
}

func (n *Node) readString(reader *strings.Reader) {
	n.typ = String
	var content string
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return
		}
		s := string(r)
		switch s {
		case "\"":
			n.content = content
			return
		case "\\":
			escChar, _, err := reader.ReadRune()
			if err != nil {
				return
			}
			content += string(escChar)
		default:
			content += s
		}
	}
}

func (n *Node) readOther(reader *strings.Reader) {
	var content string
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return
		}
		s := string(r)
		switch s {
		case "}", "]", ",":
			reader.UnreadRune()
			switch {
			case intRe.MatchString(content):
				n.typ = Int
			case floatRe.MatchString(content):
				n.typ = Float
			case boolRe.MatchString(content):
				n.typ = Bool
			case content == "null":
				// n.typ = Null
				n.isNull = true
				return
			}
			n.content = content
			return
		case " ":
			continue
		default:
			content += s
		}
	}
}

func (n *Node) parse() {
	switch n.typ {
	case Map:
		if len(n.children) == 0 {
			n.isNull = true
			return
		}
		for i := 0; i < len(n.children); i += 3 {
			n.children[i+2].parse()
			n.childrenMap[n.children[i].content] = n.children[i+2]
		}
	case Slice:
		if len(n.children) == 0 {
			n.isNull = true
			return
		}
		for i, child := range n.children {
			child.parse()
			n.childrenMap[i] = child
		}
	}
}

func (n *Node) Query(queryStr string) *Node {
	if queryStr == "" {
		return nil
	}
	queryList := queryRe.FindAllString(queryStr, -1)
	currentNode := n
	for _, query := range queryList {
		if index, err := strconv.ParseInt(query, 10, 64); err == nil {
			currentNode = currentNode.childrenMap[int(index)]
		} else {
			currentNode = currentNode.childrenMap[query]
		}
		if currentNode == nil {
			return nil
		}
	}
	return currentNode
}

func (n *Node) GetString() (string, error) {
	if n.typ != String {
		return "", errors.New("Node.GetString(): not valid string node")
	}
	return n.content, nil
}

func (n *Node) GetInt() (int, error) {
	if n.typ != Int {
		return 0, errors.New("Node.GetInt(): not valid int node")
	}
	i64, err := strconv.ParseInt(n.content, 10, 64)
	if err != nil {
		return 0, errors.New("Node.GetInt(): " + err.Error())
	}
	return int(i64), nil
}

func (n *Node) GetFloat() (float64, error) {
	if n.typ != Float {
		return 0.0, errors.New("Node.GetFloat(): not valid float node")
	}
	f64, err := strconv.ParseFloat(n.content, 64)
	if err != nil {
		return 0.0, errors.New("Node.GetFloat(): " + err.Error())
	}
	return f64, nil
}

func (n *Node) GetBool() (bool, error) {
	if n.typ != Bool {
		return false, errors.New("Node.GetBool(): not valid bool node")
	}
	switch n.content {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("Node.GetBool(): not valid bool value")
	}
}

func (n *Node) QueryString(queryStr string) (string, error) {
	sNode := n.Query(queryStr)
	if sNode == nil || sNode.IsNull() {
		return "", errors.New("Node.QueryString(): node not exist")
	}
	return sNode.GetString()
}

func (n *Node) QueryInt(queryStr string) (int, error) {
	iNode := n.Query(queryStr)
	if iNode == nil || iNode.IsNull() {
		return 0, errors.New("Node.QueryInt(): node not exist")
	}
	return iNode.GetInt()
}

func (n *Node) QueryFloat(queryStr string) (float64, error) {
	fNode := n.Query(queryStr)
	if fNode == nil || fNode.IsNull() {
		return 0, errors.New("Node.QueryFloat(): node not exist")
	}
	return fNode.GetFloat()
}

func (n *Node) QueryBool(queryStr string) (bool, error) {
	bNode := n.Query(queryStr)
	if bNode == nil || bNode.IsNull() {
		return false, errors.New("Node.QueryBool(): node not exist")
	}
	return bNode.GetBool()
}

func (n *Node) QueryValue(v interface{}, queryStr string) error {
	switch t := v.(type) {
	case *string:
		s, err := n.QueryString(queryStr)
		if err != nil {
			return err
		}
		*t = s
	case *int:
		i, err := n.QueryInt(queryStr)
		if err != nil {
			return err
		}
		*t = i
	case *float64:
		f, err := n.QueryFloat(queryStr)
		if err != nil {
			return err
		}
		*t = f
	case *bool:
		b, err := n.QueryBool(queryStr)
		if err != nil {
			return err
		}
		*t = b
	default:
		return errors.New("Node.QueryValue(): not valid value type")
	}
	return nil
}

func (n *Node) IsNull() bool {
	return n.isNull
}

func Parse(j interface{}) *Node {
	var reader *strings.Reader
	switch t := j.(type) {
	case string:
		reader = strings.NewReader(t)
	case []byte:
		reader = strings.NewReader(string(t))
	case io.ReadCloser:
		b, err := ioutil.ReadAll(t)
		if err != nil {
			return nil
		}
		reader = strings.NewReader(string(b))
		defer t.Close()
	case io.Reader:
		b, err := ioutil.ReadAll(t)
		if err != nil {
			return nil
		}
		reader = strings.NewReader(string(b))
	default:
		return nil
	}
	root := NewNode()
	r, _, err := reader.ReadRune()
	if err != nil {
		return nil
	}
	s := string(r)
	switch s {
	case "{":
		root.readMap(reader)
	case "]":
		root.readSlice(reader)
	}
	root.parse()
	return root
}
