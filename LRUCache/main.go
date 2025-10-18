// https://github.com/hashicorp/golang-lru
// go get github.com/hashicorp/golang-lru 低版本
package main

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
)

func main() {
	// Demo1()
	Demo2()
}

func Demo1() {
	// 创建一个最大容量为5的LRU缓存
	cache, err := lru.New(5)
	if err != nil {
		fmt.Println("error creating LRU cache: ", err)
		return
	}

	// 添加一些数据到缓存
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")
	cache.Add("key4", "value4")
	cache.Add("key5", "value5")

	// 获取缓存中的数据key1
	if value, ok := cache.Get("key1"); ok {
		fmt.Println("key1: ", value)
	} else {
		fmt.Println("key1 not found in cache")
	}

	// 添加超过容量的数据，触发LRU机制
	cache.Add("key6", "value6")

	// 获取缓存中的数据key2，此时key2应该被移除，因为它是最早添加且没有再次访问的数据
	if value, ok := cache.Get("key2"); ok {
		fmt.Println("key2: ", value)
	} else {
		fmt.Println("key2 not found in cache")
	}
}

func Demo2() {
	// 创建一个带有回调函数的LRU缓存
	onEvicated := func(key interface{}, value interface{}) {
		fmt.Printf("Evicted: key=%v, value=%v\n", key, value)
	}

	cache, err := lru.NewWithEvict(5, onEvicated)
	if err != nil {
		fmt.Println("Error creating LRU cache: ", err)
		return
	}

	// 添加一些数据到缓存中
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")
	cache.Add("key4", "value4")
	cache.Add("key5", "value5")

	// 添加超过容量的数据，触发 LRU 机制及回调函数
	cache.Add("key6", "value6")
}
