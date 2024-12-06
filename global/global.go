package global

import (
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go/service/s3"

	"image2webp/config"
)

var (
	DB       *gorm.DB
	Conf     = &config.Cfg{}
	S3Client = &s3.S3{}
)
