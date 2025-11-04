package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensitiveWord struct {
	Appkey       string `bson:"appkey" json:"-"`                            // 应用
	UserId       int64  `bson:"uid" json:"uid"`                             // 创建的用户
	Word         string `bson:"word" json:"word"`                           // 词
	InterceptCnt int64  `bson:"interceptCnt" json:"interceptCnt,omitempty"` // 拦截次数
	PassCnt      int64  `bson:"passCnt" json:"passCnt,omitempty"`           // 通过次数
	RejectCnt    int64  `bson:"rejectCnt" json:"rejectCnt,omitempty"`       // 拒绝次数
	Status       string `bson:"status" json:"status,omitempty"`             // 状态 closed-关闭 open - 开启
	CTime        int64  `bson:"ctime" json:"ctime"`                         // 记录的创建时间
	Method       string `bson:"method" json:"method"`                       //
}

func SensitiveWordInsert() {
	coll := client.Database("app").Collection("sensitive_word")
	sw := SensitiveWord{
		UserId: 1000,
		Appkey: "P312e64",
		Word:   "牛逼2",
		Status: "open",
		CTime:  time.Now().UnixNano() / 1000 / 1000,
	}

	qry := bson.M{"appkey": "P312e64", "word": "牛逼2"}
	opt := &options.UpdateOptions{}
	opt.SetUpsert(true)
	ret, err := coll.UpdateOne(context.TODO(), qry, bson.M{"$setOnInsert": sw}, opt)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(ret.MatchedCount)
}

func IncrSensitiveWordStatics() {
	coll := client.Database("app").Collection("sensitive_word")

	fields := bson.M{}
	fields["interceptCnt"] = 1
	fields["passCnt"] = 2
	fields["rejectCnt"] = 3

	qry := bson.M{"appkey": "P312e64", "word": "牛逼2"}
	ret, err := coll.UpdateOne(context.TODO(), qry, bson.M{"$inc": fields})
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(ret.MatchedCount)
}

func MuteTimesUpdate(appkey string, sensitiveWord string) error {
	coll := client.Database("app").Collection("sensitive_word")

	_, err := coll.UpdateOne(context.TODO(),
		bson.D{{Key: "appkey", Value: appkey}, {Key: "word", Value: sensitiveWord}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "muteCnt", Value: 1}}}},
	)
	if err != nil {
		return err
	}

	return nil
}
