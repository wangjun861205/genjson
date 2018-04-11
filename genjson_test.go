package genjson

import (
	"fmt"
	"testing"
)

func TestGenJSON(t *testing.T) {
	j := `{
		"test": {
			"hello": "wor\"ld",
			"ni": "hao"
		},
		"xxx": "yyy",
		"zzz": [1, 2, 3, 4],
		"null": null,
		"bool": [true, false],
		}`
	root := Parse(j)
	fmt.Println(root)
	q := root.Query("test.hello")
	fmt.Println(q)
}
