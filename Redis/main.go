// https://github.com/redis/go-redis
// go get github.com/go-redis/redis/v8 低版本
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

	timestamp := 1760507538
	t := time.Unix(int64(timestamp), 0)
	format := t.Format("2006-01-02 15:04:05")

	fmt.Println(format)
	// Example()
	HashExample()
}

func HashExample() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// count, err := client.HIncrBy(context.Background(), "1:d5a233a635c3baed5a152dacd6181671", "count", 1).Result()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// fmt.Println("count: ", count)

	// userId := 1001
	// err = client.HSet(context.Background(), "1:d5a233a635c3baed5a152dacd6181671", userId, 1760507513).Err()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// userId := 1
	// err := client.HSet(client.Context(), "2:d5a233a635c3baed5a152dacd6181671", map[string]interface{}{
	// 	"count":            1,
	// 	fmt.Sprint(userId): 1760507513},
	// ).Err()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// err := client.HSet(client.Context(), "1:d5a233a635c3baed5a152dacd6181671", "count", 1, "1", 1760507513).Err()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	data, err := client.HGetAll(context.Background(), "1:d5a233a635c3baed5a152dacd6181671").Result()
	if err != nil {
		log.Fatal(err)
		return
	}

	var minTimestamp int64 = math.MaxInt64
	for k, v := range data {
		if k == "count" {
			client.HDel(context.TODO(), "1:d5a233a635c3baed5a152dacd6181671", "count")
		} else {
			timestamp, _ := strconv.ParseInt(v, 10, 64)

			if timestamp < minTimestamp {
				minTimestamp = timestamp
			}
			client.HDel(context.TODO(), "1:d5a233a635c3baed5a152dacd6181671", k)
		}
	}

	// for k := range data {
	// 	if k != "count" {
	// 		uid, _ := strconv.ParseInt(k, 10, 64)
	// 		client.HSet(context.TODO(), "1:d5a233a635c3baed5a152dacd6181671", uid, 1768569873)
	// 	}
	// }

	// count, err := client.HIncrBy(context.Background(), "2:d5a233a635c3baed5a152dacd6181671", "count", 1).Result()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// fmt.Println("count: ", count)

	// // 根据key和field字段查询field字段的值
	// countStr, err := client.HGet(context.Background(), "2:d5a233a635c3baed5a152dacd6181671", "count").Result()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// fmt.Println("count: ", countStr)

	// client.HDel(context.Background(), "2:d5a233a635c3baed5a152dacd6181671", "count", "1").Result()
}

func Example() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Fatal(err)
		return
	}

	val, err := rdb.Get(context.Background(), "key").Result()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("key = ", val)
}
