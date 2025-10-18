package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tea struct {
	Type     string   `bson:"type"`     // 类型
	Category string   `bson:"category"` // 分类
	Toppings []string `bson:"toppings"` // 配料
	Price    float32  `bson:"price"`    // 价格
}

func AggregationData() {
	coll := client.Database("kingkong").Collection("teas")
	docs := []any{
		Tea{Type: "Masala", Category: "black", Toppings: []string{"ginger", "pumpkin spice", "cinnamon"}, Price: 6.75},
		Tea{Type: "Gyokuro", Category: "green", Toppings: []string{"berries", "milk foam"}, Price: 5.65},
		Tea{Type: "English Breakfast", Category: "black", Toppings: []string{"whipped cream", "honey"}, Price: 5.75},
		Tea{Type: "Sencha", Category: "green", Toppings: []string{"lemon", "whipped cream"}, Price: 5.15},
		Tea{Type: "Assam", Category: "black", Toppings: []string{"milk foam", "honey", "berries"}, Price: 5.65},
		Tea{Type: "Matcha", Category: "green", Toppings: []string{"whipped cream", "honey"}, Price: 6.45},
		Tea{Type: "Earl Grey", Category: "black", Toppings: []string{"milk foam", "pumpkin spice"}, Price: 6.15},
		Tea{Type: "Hojicha", Category: "green", Toppings: []string{"lemon", "ginger", "milk foam"}, Price: 5.55},
	}
	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Document IDs: %d", result.InsertedIDs...)
}

func AggregationOperator() {
	coll := client.Database("kingkong").Collection("teas")
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "average_price", Value: bson.D{{Key: "$avg", Value: "$price"}}},
			{Key: "type_total", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
	if err != nil {
		log.Fatal(err)
		return
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Printf("Average price of %v tea options: $%v \n", result["_id"], result["average_price"])
		fmt.Printf("Number of %v tea options: %v \n\n", result["_id"], result["type_total"])
	}
}

func AggregationOperator2() {
	coll := client.Database("kingkong").Collection("teas")
	// $match 阶段，用于与 toppings 字段包含“牛奶泡沫（milk foam）”的文档进行匹配
	// $unset 阶段，用于省略 _id 和 category 字段
	// $sort 阶段，用于按升序对 price 和 toppings 进行排序
	// $limit 阶段以显示前两个文档

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "toppings", Value: "milk foam"}}}}
	unsetStage := bson.D{{Key: "$unset", Value: bson.A{"_id", "category"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "price", Value: 1}, {Key: "toppings", Value: 1}}}}
	limitStage := bson.D{{Key: "$limit", Value: 2}}
	// Performs the aggregation and prints the results
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, unsetStage, sortStage, limitStage})
	if err != nil {
		panic(err)
	}
	var results []Tea
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Printf("Tea: %v \nToppings: %v \nPrice: $%v \n\n", result.Type, strings.Join(result.Toppings, ", "), result.Price)
	}
}
