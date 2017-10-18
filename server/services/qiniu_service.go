package services

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/xrlin/WebIM/server/config"
)

func GenerateUploadToken() string {
	putPolicy := storage.PutPolicy{
		Scope: config.QiniuCfg.Bucket,
	}
	mac := qbox.NewMac(config.QiniuCfg.AccessKey, config.QiniuCfg.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
