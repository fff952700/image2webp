package main

import (
	"image2webp/handler"
	"image2webp/initialize"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMysql()
	initialize.InitS3()
	handler.ImageServer()
}
