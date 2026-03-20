// https://github.com/qiniu/go-sdk
// go get github.com/qiniu/go-sdk/v7@v7.25.4
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

func main() {
	// 获取上传凭证
	upToken, err := GenerateToken()
	if err != nil {
		return
	}
	fmt.Println(upToken)

	// 配置上传参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 华东存储区域
		UseHTTPS:      true,
		UseCdnDomains: true,
	}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)

	type MyPutRet struct {
		Key    string
		Hash   string
		Fsize  int
		Bucket string
		Name   string
	}

	ret := MyPutRet{}

	// 本地文件路径
	localFile := "./images/touxiang.jpg"
	// 上传到七牛云的文件名，包含user目录前缀
	key := "user/touxiang1.jpg"
	// 可选：设置上传进度
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
	if err != nil {
		fmt.Println("上传失败:", err)
		return
	}

	fmt.Printf("上传成功!\n")
	fmt.Printf("Key: %s\n", ret.Key)
	fmt.Printf("Hash: %s\n", ret.Hash)
	fmt.Printf("文件访问地址: http://%s/%s\n", domain, key)

}

var accessKey string = "xDQ0t244ITe41pucUVw_L3BemEHUo-I8D5X_hSRQ"
var secretKey string = "sYtmYIBxfKAmdhHpPY9pBbVSBy6uaJ2hMhh8snON"
var bucket string = "dongshan2026"
var domain string = "tadlyqdtk.hd-bkt.clouddn.com"

func GenerateToken() (string, error) {
	mac := credentials.NewCredentials(accessKey, secretKey)
	bucket := "dongshan2026"

	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}

	putPolicy.SetReturnBody(`{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`)
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return "", err
	}

	return upToken, nil
}
