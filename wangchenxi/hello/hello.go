package main

import "fmt"

func main() {
	// tmp := fmt.Sprint()
	tmp := fmt.Sprintf("%v,%v", "hello", "test")
	fmt.Println(tmp)
	testStr := make([]int, 2)
	testStr[0] = 1
	testStr[1] = 1
	// testStr[2] = 1
	// testStr[3] = 1
	// testStr[4] = 1
	fmt.Println(testStr, len(testStr), cap(testStr))
	t := append(testStr, 1)
	fmt.Println(testStr, len(testStr), cap(testStr))
	fmt.Println(t, len(t), cap(t))
	mapslice := make([]map[string]interface{}, 2, 3)
	mapslice[0] = make(map[string]interface{}, 2)
	mapslice[0]["hello"] = 123
	mapslice[1] = make(map[string]interface{}, 2)
	mapslice[1]["hello1"] = 123
	tm := make(map[string]interface{}, 2)
	tmt := append(mapslice, tm)
	tm["hello2"] = 123
	fmt.Println(mapslice)
	fmt.Println(tmt)
	tmt[0]["test"] = 321
	fmt.Println(mapslice)
	fmt.Println(tmt)

}
