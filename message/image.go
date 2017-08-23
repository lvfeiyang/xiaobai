package message

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/lvfeiyang/xiaobai/common/config"
)

func QiniuToken(bucket string) string {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(config.ConfigVal.QiniuAK, config.ConfigVal.QiniuSK)
	return putPolicy.UploadToken(mac)
}
