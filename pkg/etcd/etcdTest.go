package etcd

import "fmt"

//	func EtcdTest() {
//		EtcdInit("")
//		//Put("test", "2333")
//		//Put("test1/test1_1", "test1_1")
//		//Put("test1/test1_2", "test1_2")
//		//Put("test2/test2_1", "test2_1")
//		//Put("test2/test2_2", "test2_2")
//		//Put("test2/test2_3", "test2_3")
//	}
func main() {
	EtcdInit("")
	Put("test", "test")
	Put("test1/test1_1", "test1_1")
	Put("test1/test1_2", "test1_2")
	Put("test2/test2_1", "test2_1")
	Put("test2/test2_2", "test2_2")
	Put("test2/test2_3", "test2_3")

	val := GetOne("test")
	fmt.Printf("test: %s\n", val)

	fmt.Printf("test1:\n")
	var slice map[string]string
	slice = GetDirectory("test1/")
	for k, v := range slice {
		fmt.Printf("%s: %s\n", k, v)
	}

	fmt.Printf("test2:\n")
	slice = GetDirectory("test2/")
	for k, v := range slice {
		fmt.Printf("%s: %s\n", k, v)
	}
	EtcdDeinit()
}
