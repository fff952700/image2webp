package initialize

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"image2webp/global"
)

func InitS3() {
	// 创建新的 AWS 会话，并指定访问密钥和私密密钥
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(global.Conf.BucketInfo.Region),
		Credentials: credentials.NewStaticCredentials(global.Conf.BucketInfo.AccessKey, global.Conf.BucketInfo.SecretKey, ""), // 使用静态凭证
	})
	if err != nil {
		zap.S().Fatalw("init S3 error", "err", err)
	}

	// 创建 S3 服务客户端
	global.S3Client = s3.New(sess)
}
