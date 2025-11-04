package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginRsp struct {
	UserId             int64    `json:"uid"`                          // uid
	Token              string   `json:"token"`                        // 授权，后续的接口必传
	Appkey             string   `json:"appkey"`                       // 子应用 与公司关联
	Plat               string   `json:"plat"`                         // B P F端
	Company            string   `json:"-"`                            // 公司全称
	CompanyShort       string   `json:"-"`                            // 公司简称
	Category           []string `json:"-"`                            // 分类
	EmployeeID         string   `json:"-"`                            // 工号
	Position           string   `json:"-"`                            // 职位
	Secret             string   `json:"secret,omitempty"`             // 秘钥，有秘钥消息就加密，没有或者没有就不加密
	Encrypt            int64    `json:"encrypt,omitempty"`            // 当前是否加密 1-加密 其余不加密
	MsgAssistant       int64    `json:"msgAssistant,omitempty"`       // 当前开启群发助手 1-开启 其余不开启
	MsgTranslate       int64    `json:"msgTranslate,omitempty"`       //翻译助手，1-开启翻译助手，其余不开启
	CustomerServiceUrl string   `json:"customerServiceUrl,omitempty"` // 客服url
}

type LoginRspConfig struct {
	SyncUri    string `json:"syncUri"`    // 长轮训监听的地址
	NickGuide  bool   `json:"nickGuide"`  // 是否进入昵称修改引导
	LoginGuide bool   `json:"loginGuide"` // 是否进入登录引导
}

type PcLoginRecordDao struct {
	Id             string         `bson:"_id"`            // DevicedId:Account:Timestamp
	DeviceId       string         `bson:"deviced"`        // 设备ID，唯一
	Account        string         `bson:"account"`        // 账号 可以是用户名/手机号/邮箱
	Timestamp      int64          `bson:"timestamp"`      // 同一次PC端登录，前端多次调用时，该值要保持不变
	Status         int            `bson:"status"`         // 0:pending 1:scan 2:cancel 3:confirm
	LoginRsp       LoginRsp       `bson:"loginRsp"`       // 登登录响应信息
	LoginRspConfig LoginRspConfig `bson:"loginRspConfig"` // 登录响应配置信息
}

func InsertOneX(deviceId string, account string, timestamp int64) error {
	pcl := PcLoginRecordDao{
		Id:        GenerateId(deviceId, account, timestamp),
		DeviceId:  deviceId,
		Account:   account,
		Timestamp: timestamp,
		Status:    0,
	}
	coll := client.Database("kingkong").Collection("loginrecord")
	_, err := coll.InsertOne(context.TODO(), pcl)
	if err != nil {
		return err
	}

	return nil
}

func FindOneX(deviceId string, account string, timestamp int64) (*PcLoginRecordDao, error) {
	var pcl PcLoginRecordDao
	coll := client.Database("kingkong").Collection("loginrecord")
	err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: GenerateId(deviceId, account, timestamp)}}).Decode(&pcl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 表中没有数据，错误返回nil
		} else {
			return nil, err // 查询过程中出现错误，返回原始错误
		}
	}
	return &pcl, nil // 查询到数据
}

func UpdateOneX(deviceId string, account string, timestamp int64, status int) error {
	coll := client.Database("kingkong").Collection("loginrecord")
	_, err := coll.UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: GenerateId(deviceId, account, timestamp)}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}}}},
	)
	if err != nil {
		return err
	}

	return nil
}

func GenerateId(deviceId string, account string, timestamp int64) string {
	originId := fmt.Sprintf("%s:%s:%d", deviceId, account, timestamp)
	encodedId := base64.StdEncoding.EncodeToString([]byte(originId))
	return encodedId
}
