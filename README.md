# genjson

## Overview
Generic JSON parse library

## Examples

```
  import (
    "genjson"
    "fmt"
    "logging"
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
  
  func main() {
    root := genjson.Parse(jsonStr)
    s, err := root.QueryString("str")
    if err != nil {
      logging.Fatal(err)
      }
    fmt.Printf("str value: %s\n", s)
    //output: str value: hello world

    i, err := root.QueryInt("int")
    if err != nil {
      logging.Fatal(err)
      }
    fmt.Printf("int value: %d\n", i)
    //output: int value: 1

    f, err := root.QueryFloat("float")
    if err != nil {
      logging.Fatal(err)
      }
    fmt.Printf("float value: %f\n", f)
    //output: float value: 1.23

    b, err := root.QueryBool("bool")
    if err != nil {
      logging.Fatal(err)
      }
    fmt.Printf("bool value: %t\n", b)
    //output: bool value: true
  
    nullNode := root.Query("Null")
    fmt.Printf("nullNode is null: %t\n", nullNode.IsNull())
    //output: nullNode is null: ture

    objField1, err := root.QueryString("obj.field1")
    if err != nil {
      logging.Fatal(err)
      }
    fmt.Printf("object field1 value: %s\n", objField1)
    //output: object field1 value: foo

    arrayNum, err := root.QueryInt("array[0]")
    if err != nil {
      logging.Fatal(err)
      }
    fmt.Printf("array first value: %d\n", arrayNum)
    //output: array first value: 4
    }
```
## Usage
  1. **Parse JSON string**:  
    ```root := genjson.Parse(jsonStr)```  
    You can pass string, []byte, io.Reader into Parse(j interface{}), it returns a *genjson.Node object, **you should always check if the return value is nil**.
  2. **Query node**:  
    ```node := parentNode.Query("xxx.yyy.zzz[0]")```  
    Query(queryStr string) method require a query string and return a *genjson.Node, if the node which is queried not exist, it will return a nil (about [query string format][query string format])
  3. **Query string**:  
    ```s, err := parentNode.QueryString("xxx.yyy.zzz[0]")```  
    QueryString(queryStr string) method require a query string and return string node value, if type of the node is not string or node not exist, it will return a empty string and a error (about [query string format](#query-string-format))
  4. **Query int**:  
    ```i, err := parentNode.QueryInt("xxx.yyy.zzz[0]")```  
    Like QueryString(queryStr string), but QueryInt() method returns a int value
  5. **Query float**:  
    ```f, err := parentNode.QueryFloat("xxx.yyy.zzz[0]")```  
    Like QueryString(queryStr string), but QueryFloat() method returns a float64 value
  6. **Query bool**:  
    ```b, err := parentNode.QueryBool("xxx.yyy.zzz[0]")```  
    Like QueryString(), but QueryBool() method returns a bool value
  7. **Query value**:  
    ```var s string
       err := node.QueryValue(&s, "xxx.yyy.zzz[0]")
    ```  
     QueryValue(v interface{}, queryStr string), require a pointer (*string, *int, *float64, *bool) and a query string, if the type of the pointer not fit to node type or node not exist, it returns a error else return nil 
  8. **Check null node**:  
    ``` result := node.IsNull() ```  
    IsNull() returns true if node is null node, else return false.Before invoke IsNull(), you should ensure node pointer is not nil.
    
## Query String Format
    ``` {
          "xxx": {
                    "yyy": {
                              "zzz": [1, 2, 3, 4]
                           }
                 }
        }```
 If you want to get object field you can use "object.fieldName", if you want to get array element you can use "array[index]" 
    

