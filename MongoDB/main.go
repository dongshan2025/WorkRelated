// https://github.com/mongodb/mongo-go-driver
// go get go.mongodb.org/mongo-driver/mongo 低版本
// go get github.com/joho/godotenv 使用 godotenv 包从环境变量中读取 MongoDB 连接字符串，避免在源代码中嵌入凭据。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var uri = "mongodb://admin:123456@localhost:27017/?readPreference=primary&ssl=false"
// var uri = "mongodb://127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019/?replicaSet=myReplicaSet&directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.8"
var uri = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=myReplicaSet"

// var uri = "mongodb://localhost:27018/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.8"

// clusterUri := "mongodb://username:password@host1:27017,host2:27017,host3:27017/?replicaSet=myReplicaSet"
var client *mongo.Client

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
}

func main() {
	// InsertOne()
	// InsertMultiple()
	// FindOne()
	// FindMultiple()
	// UpdateOne()
	// UpdatMultiple()
	// UpdatMultipleArrayField()
	// ReplaceOne()
	// FindOneAndReplace()
	// DeleteOne()
	// DeleteMultiple()
	// BulkWrite()
	// MonitorDataChanges()
	// CountDocument()
	// DistinctDocument()
	// LimitNumber()
	// SkipNumber()
	// SortDocument()
	// ProjectDocument()
	// CreateIndex()
	// AggregationData()
	// AggregationOperator()
	// AggregationOperator2()

	// RTCCreate()
	// RTCQuery()
	// RTCUpdateStauts()
	UpdateOneMulField()
}

func Example() {
	// 设置上下文，用于连接超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 确保在函数结束时取消上下文，避免资源泄漏

	// 构建连接字符串，根据你的实际配置进行修改
	uri := "mongodb://admin:123456@localhost:27017/?readPreference=primary&ssl=false"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Connected to MongoDB!")

	// 插入一条记录
	collection := client.Database("kingkong").Collection("users")
	user := bson.D{{Key: "name", Value: "zhangshan"}, {Key: "age", Value: 30}, {Key: "email", Value: "zhangshan@gmail.com"}}
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inserted document with ID: %v\n", result.InsertedID)

	// 查询一条记录
	name := "zhangshan"
	var res bson.M
	err = collection.FindOne(context.TODO(), bson.D{{Key: "name", Value: "zhangshan"}}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", name)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
