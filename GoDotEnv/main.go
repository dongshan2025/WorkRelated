// https://github.com/joho/godotenv
// go get github.com/joho/godotenv
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Demo1()
	// Demo2()
	// Demo3()
	Demo4()
}

func Example() {
	err := godotenv.Load() // 将当前目录下的.env配置文件读取到环境变量中
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("name: %s", os.Getenv("name"))
	log.Printf("age: %s", os.Getenv("age"))
	log.Printf("gopath: %s", os.Getenv("GOPATH"))
}
