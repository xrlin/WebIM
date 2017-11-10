package main

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

func GenerateUploadToken() string {
	putPolicy := storage.PutPolicy{
		Scope: QiniuCfg.Bucket,
	}
	mac := qbox.NewMac(QiniuCfg.AccessKey, QiniuCfg.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
