package initialize

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"image2webp/global"
)

func InitConfig() {

	// 如果配置文件在根目录下
	_, err := toml.DecodeFile("config.toml", global.Conf)
	if err != nil {
		zap.S().Panicw("读取配置文件失败", "err", err)
	}
}
