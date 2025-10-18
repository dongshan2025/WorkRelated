package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func Demo1() {
	log.Printf("name :%s", os.Getenv("name"))
	log.Printf("age: %s", os.Getenv("age"))
}

func Demo2() {
	// 没有指定".env"文件，该文件也会自动加载进来
	err := godotenv.Load("common", "dev.env", "production.env")
	if err != nil {
		log.Fatal(err)
		return
	}

	// database参数在"dev.env"和"production.env"中都存在，谁先加载就使用哪个的配置
	log.Printf("name: %s, age: %s, sex: %s, city: %s, database: %s", os.Getenv("name"), os.Getenv("age"), os.Getenv("sex"), os.Getenv("city"), os.Getenv("database"))
}

func Demo3() {
	err := godotenv.Load("test.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("code: %s, version: %s", os.Getenv("code"), os.Getenv("version"))
}

// 可以不将.env文件内容存入环境变量，使用godotenv.Read()返回一个map[string]string，可直接使用
func Demo4() {
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("name: %s, age: %s", myEnv["name"], myEnv["age"])
}
