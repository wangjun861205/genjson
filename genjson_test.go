package genjson

import (
	"fmt"
	"log"
	"testing"
)

var jsonStr = `{
    "str": "hello world",
    "int": 1,
    "float": 1.23,
    "bool": true,
    "Null": null,
    "obj":{
        "field1": "foo",
        "field2": 100
        },
    "array":[4, 5, 6, 7, 8]
    }`

func TestGenJSON(t *testing.T) {
	root := Parse(jsonStr)
	s, err := root.QueryString("str")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("str value: %s\n", s)
	//output: str value: hello world

	i, err := root.QueryInt("int")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("int value: %d\n", i)
	//output: int value: 1

	f, err := root.QueryFloat("float")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("float value: %f\n", f)
	//output: float value: 1.23

	b, err := root.QueryBool("bool")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bool value: %t\n", b)
	//output: bool value: true

	nullNode := root.Query("Null")
	fmt.Printf("nullNode is null: %t\n", nullNode.IsNull())
	//output: nullNode is null: ture

	objField1, err := root.QueryString("obj.field1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("object field1 value: %s\n", objField1)
	//output: object field1 value: foo

	arrayNum, err := root.QueryInt("array[0]")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("array first value: %d\n", arrayNum)
	//output: array first value: 4

	m, err := root.QueryMap("obj")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("map value: %v\n", m)
	//output: map value: map[field1:foo field2:100]

	l, err := root.QuerySlice("array")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("array value: %v\n", l)
	//output: array value: array

	fmt.Println(root)
}
