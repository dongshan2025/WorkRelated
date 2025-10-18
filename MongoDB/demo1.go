package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name   string   `bson:"name"`
	Age    int      `bson:"age"`
	Email  string   `bson:"email"`
	Family []string `bson:"family"`
}

// 注意：两种插入方式对空值的处理方式不同
// bson.D插入时，family字段没有值，在数据库中的值为：(N/A)
// user插入时，family字段没有值，在数据库中的值为：(Null)
func InsertOne() {
	coll := client.Database("kingkong").Collection("users")
	// 插入一条记录
	result, err := coll.InsertOne(context.TODO(),
		bson.D{
			{Key: "name", Value: "lisi"},
			{Key: "age", Value: 21},
			{Key: "email", Value: "lisi@gmail.com"},
		})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Inserted a single document: %v\n", result.InsertedID)
	// 插入一条记录
	user := User{
		Name:  "liuer",
		Age:   22,
		Email: "liuer@gmail.com",
	}
	result, err = coll.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Inserted a single document: %v\n", result.InsertedID)
}

func InsertMultiple() {
	coll := client.Database("kingkong").Collection("users")
	// 插入多条记录
	result, err := coll.InsertMany(context.TODO(), []interface{}{
		bson.D{{Key: "name", Value: "zhangshan"}, {Key: "age", Value: 17}, {Key: "email", Value: "zhangshan@gmail.com"}},
		bson.D{{Key: "name", Value: "lisi"}, {Key: "age", Value: 18}, {Key: "email", Value: "lisi@gmail.com"}},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Inserted multiple document: %v\n", result.InsertedIDs...)
	// 插入多条记录
	users := []interface{}{
		User{Name: "wangwu", Age: 19, Email: "wangwu@gmail.com", Family: []string{"fathor"}},
		User{Name: "zhaoliu", Age: 20, Email: "zhaoliu@gmail.com", Family: []string{"mother"}},
	}
	result, err = coll.InsertMany(context.TODO(), users)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Inserted multiple document: %v\n", result.InsertedIDs...)

}

func FindOne() {
	coll := client.Database("kingkong").Collection("users")
	// 查询一条记录
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: "zhangshan"}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(result)
	// 查询一条记录
	var user User
	err = coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: "zhangshan"}}).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Name: %s, Age: %d, Email: %s", user.Name, user.Age, user.Email)
}

func FindMultiple() {
	coll := client.Database("kingkong").Collection("users")
	// 查询多条记录
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 20}}}})
	if err != nil {
		log.Fatal(err)
		return
	}
	// 迭代获取
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(result)
	}
	// 一次性获取
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(results)

	// 结构体迭代获取
	cursor, err = coll.Find(context.TODO(), bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 20}}}})
	if err != nil {
		log.Fatal(err)
		return
	}

	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Name: %s, Age: %d, Email: %s\n", user.Name, user.Age, user.Email)
	}

	// 一次性获取
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(users)
}

func UpdateOne() {
	coll := client.Database("kingkong").Collection("users")
	// 修改一条记录
	result, err := coll.UpdateOne(context.TODO(),
		bson.D{{Key: "name", Value: "zhangshan"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 28}}}},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)
}

func UpdateOneMulField() {
	coll := client.Database("kingkong").Collection("users")
	// 修改一条记录
	result, err := coll.UpdateOne(context.TODO(),
		bson.D{{Key: "name", Value: "lisi"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 28}, {Key: "email", Value: "lisi28@gmail.com"}}}},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)
}

func UpdatMultiple() {
	coll := client.Database("kingkong").Collection("users")
	// 修改多条记录
	result, err := coll.UpdateMany(context.TODO(),
		bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 19}}}, {Key: "age", Value: bson.D{{Key: "$lte", Value: 20}}}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 27}}}},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)
}

func UpdatMultipleArrayField() {
	coll := client.Database("kingkong").Collection("users")
	result, err := coll.UpdateMany(context.TODO(),
		bson.D{}, // 没有加查询条件，所有记录的family字段都增加brother家庭成员
		bson.D{{Key: "$push", Value: bson.D{{Key: "family", Value: "brother"}}}},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)
}

// 将满足查询条件的第一个文档替换为新的文档
// 如果需要替换多个文档，需要使用循环多次调用ReplaceOne方法
// 只返回操作结果
func ReplaceOne() {
	coll := client.Database("kingkong").Collection("users")
	result, err := coll.ReplaceOne(context.TODO(),
		bson.D{{Key: "name", Value: "zhangshan"}},
		bson.D{{Key: "name", Value: "zhangshan2"}, {Key: "age", Value: 29}, {Key: "email", Value: "zhangshan2@gmail.com"}, {Key: "family", Value: []string{"father"}}},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The number of replace documents: %d\n", result.MatchedCount)
	fmt.Printf("The number of replace documents: %d\n", result.ModifiedCount)
}

// 将满足查询条件的第一个文档替换为新的文档
// 返回替换前的文档（根据选项设置）
func FindOneAndReplace() {
	coll := client.Database("kingkong").Collection("users")
	var user User
	// 默认返回替换前的旧文档
	err := coll.FindOneAndReplace(context.TODO(),
		bson.D{{Key: "name", Value: "zhangshan"}},
		bson.D{{Key: "name", Value: "zhangshan2"}, {Key: "age", Value: 15}, {Key: "email", Value: "zhangshan@gmail.com"}},
	).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(user)

	// 设置返回替换后的新文档
	opts := options.FindOneAndReplace().SetReturnDocument(options.After)
	err = coll.FindOneAndReplace(context.TODO(),
		bson.D{{Key: "name", Value: "zhangshan2"}},
		bson.D{{Key: "name", Value: "zhangshan"}, {Key: "age", Value: 17}, {Key: "email", Value: "zhangshan@gmail.com"}}, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(err)
			return
		}
		log.Fatal(err)
		return
	}
	fmt.Println(user)

}

func DeleteOne() {
	coll := client.Database("kingkong").Collection("users")
	// 删除一条文档
	result, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "name", Value: "zhangshan"}})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(result.DeletedCount)
}

func DeleteMultiple() {
	coll := client.Database("kingkong").Collection("users")
	// 删除多条记录
	result, err := coll.DeleteMany(context.TODO(), bson.D{{Key: "age", Value: bson.D{{Key: "$lte", Value: 19}}}})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(result.DeletedCount)
}

func BulkWrite() {
	coll := client.Database("kingkong").Collection("users")
	// 批量写入
	// 匹配name为“wangwu”的文档，并将其替换为新文档
	// 匹配name为“zhaoliu”的文档，并将age值更新为39
	// 插入新文档“zhangshan”
	models := []mongo.WriteModel{
		mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "name", Value: "wangwu"}}).SetReplacement(User{Name: "wangwu2", Age: 29, Email: "wangwu2@gmail.com", Family: []string{"father"}}),
		mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "name", Value: "zhaoliu"}}).SetUpdate(bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 39}}}}),
		mongo.NewInsertOneModel().SetDocument(User{Name: "zhangshan", Age: 18, Email: "zhangshan@gmail.com", Family: []string{"father", "monther"}}),
	}
	// 指定按顺序执行批量操作
	opts := options.BulkWrite().SetOrdered(true)

	result, err := coll.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("MatchedCount: %d, ModifiedCount: %d, InsertedCount: %d, UpsertedCount: %d", result.MatchedCount, result.ModifiedCount, result.InsertedCount, result.UpsertedCount)
}

func MonitorDataChanges() {
	coll := client.Database("kingkong").Collection("users")
	pipeline := mongo.Pipeline{bson.D{{Key: "$match", Value: bson.D{{Key: "operationType", Value: "insert"}}}}}
	cursor, err := coll.Watch(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cursor.Close(context.Background())

	// 循环获取变化事件
	for cursor.Next(context.Background()) {
		var event interface{}
		if err := cursor.Decode(&event); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Change detected: %v\n", event)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func CountDocument() {
	coll := client.Database("kingkong").Collection("users")
	count, err := coll.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Document count: %d", count)
}

func DistinctDocument() {
	coll := client.Database("kingkong").Collection("users")
	result, err := coll.Distinct(context.TODO(), "name", bson.D{})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(result...)
}

func LimitNumber() {
	coll := client.Database("kingkong").Collection("users")
	cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetLimit(2))
	if err != nil {
		log.Fatal(err)
		return
	}

	// 迭代获取
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(result)
	}
}

func SkipNumber() {
	coll := client.Database("kingkong").Collection("users")
	cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetSkip(1))
	if err != nil {
		log.Fatal(err)
		return
	}

	// 迭代获取
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(result)
	}
}

func SortDocument() {
	coll := client.Database("kingkong").Collection("users")
	cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{{Key: "age", Value: -1}})) // 1:ascending -1:descending
	if err != nil {
		log.Fatal(err)
		return
	}

	// 迭代获取
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(user)
	}
}

func ProjectDocument() {
	coll := client.Database("kingkong").Collection("users")
	cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.D{{Key: "age", Value: 0}, {Key: "_id", Value: 0}}))
	if err != nil {
		log.Fatal(err)
		return
	}

	// 迭代获取
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(result)
	}
}

func CreateIndex() {
	coll := client.Database("kingkong").Collection("users")
	// 创建联合索引：name字段升序，age字段降序
	model := mongo.IndexModel{Keys: bson.D{{Key: "name", Value: 1}, {Key: "age", Value: -1}}}
	name, err := coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Index name: %s", name) // name_1_age_-1
}

func SearchText() {
	coll := client.Database("kingkong").Collection("users")
	// only searches fields with text indexes
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: "beagle"}}}})
	if err != nil {
		log.Fatal(err)
		return
	}

	// 迭代获取
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(result)
	}
}
