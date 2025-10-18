package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type RTCMember struct {
	UserId   int64  `bson:"userId" json:"userId"`     // 用户ID
	SDKToken string `bson:"sdkToken" json:"SDKToken"` // 用户Token
	Status   int32  `bson:"status" json:"status"`     // 状态 -1:canceled  0:pending  1:joined
}

type RTCDao struct {
	Id          int64        `bson:"_id" json:"rtc_sid"`                     // 唯一ID，雪花算法产生
	Type        string       `bson:"type" json:"type,omitempty"`             // 消息类型
	Appkey      string       `bson:"appkey" json:"appkey"`                   // 应用标识
	ClientId    string       `bson:"cid" json:"cid,omitempty"`               // 客户端生成的唯一的消息Id，用于客户端消息去重
	TopicId     string       `bson:"tid" json:"tid,omitempty"`               // 会话ID
	MsgId       string       `bson:"mid" json:"mid,omitempty"`               // 消息ID
	Sender      int64        `bson:"sender" json:"sender,omitempty"`         // 发送方
	STime       int64        `bson:"stime" json:"stime,omitempty"`           // 发送时间 单位毫秒
	MemberCount int32        `bson:"memberCnt" json:"memberCnt,omitempty"`   // 目标用户数
	MemberUids  []int64      `bson:"memberUids" json:"memberUids,omitempty"` // 目标用户id集
	MemberList  []*RTCMember `bson:"memberList" json:"memberList,omitempty"` // 目标用户列表
	Status      int32        `bson:"status" json:"status,omitempty"`         // 状态 -1:canceled  0:pending  1:ready  2:started  3:done
}

func RTCUpdateStauts() {
	coll := client.Database("msg").Collection("rtc")
	// 修改一条记录
	result, err := coll.UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: 1000000001}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 3}}}},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)
}

func RTCCreate() {
	coll := client.Database("msg").Collection("rtc")
	// rtc := RTCDao{
	// 	Id:          1000000001,
	// 	Type:        "communication",
	// 	Appkey:      "BeeAppKey",
	// 	ClientId:    "123-456-789-0001",
	// 	TopicId:     "88-001",
	// 	MsgId:       "99-001",
	// 	Sender:      10001,
	// 	STime:       time.Now().UnixNano() / 1000 / 1000,
	// 	MemberCount: 2,
	// 	MemberUids:  []int64{10001, 10002},
	// 	MemberList: []*RTCMember{
	// 		{
	// 			UserId:   10001,
	// 			SDKToken: "user1_token",
	// 			Status:   1,
	// 		},
	// 		{
	// 			UserId:   10002,
	// 			SDKToken: "user2_token",
	// 			Status:   0,
	// 		},
	// 	},
	// 	Status: 0,
	// }

	rtc := RTCDao{
		Id:          1000000003,
		Type:        "communication",
		Appkey:      "BeeAppKey",
		ClientId:    "123-456-789-0001",
		TopicId:     "88-001",
		MsgId:       "99-001",
		Sender:      10004,
		STime:       time.Now().UnixNano() / 1000 / 1000,
		MemberCount: 2,
		MemberUids:  []int64{10004, 10002},
		MemberList: []*RTCMember{
			{
				UserId:   10004,
				SDKToken: "user1_token",
				Status:   1,
			},
			{
				UserId:   10002,
				SDKToken: "user2_token",
				Status:   0,
			},
		},
		Status: 0,
	}

	result, err := coll.InsertOne(context.TODO(), rtc)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(result.InsertedID)
}

func RTCQuery() {
	coll := client.Database("msg").Collection("rtc")
	// 查询当前用户不属于当前Channel的RTC记录
	cursor, err := coll.Find(context.TODO(), bson.D{
		{Key: "_id", Value: bson.D{{Key: "$ne", Value: 1000000001}}},
		{Key: "status", Value: bson.D{{Key: "$in", Value: []int{0, 1}}}},
		{Key: "memberList.userId", Value: 10002},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// 遍历查询到的记录，并更改当前用户的memberList.status值为-1:canceled

	// 一次性获取
	var results []RTCDao
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		return
	}

	// 遍历查询到的记录，并更改当前用户的memberList.status值为-1:canceled
	for _, res := range results {
		for i, mem := range res.MemberList {
			if mem.UserId == 10002 {
				// 更改其状态值为-1:canceled
				// updateResult, err := coll.UpdateOne(context.TODO(),
				// 	bson.D{{Key: "_id", Value: res.Id}, {Key: "memberList.userId", Value: mem.UserId}},
				// 	bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("memberList.%d.status", i), Value: -1}}}},
				// )
				updateResult, err := coll.UpdateOne(context.TODO(),
					bson.D{{Key: "_id", Value: res.Id}},
					bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("memberList.%d.status", i), Value: -1}}}},
				)
				if err != nil {
					log.Fatal(err)
					return
				}
				fmt.Println(updateResult.ModifiedCount)
			}
		}
	}
}

// coll := client.Database("kingkong").Collection("users")
// 	// 修改一条记录
// 	result, err := coll.UpdateOne(context.TODO(),
// 		bson.D{{Key: "name", Value: "zhangshan"}},
// 		bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 28}}}},
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	fmt.Printf("The number of modified documents: %d\n", result.ModifiedCount)
