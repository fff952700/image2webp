package handler

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chai2010/webp"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"

	"image2webp/global"
	"image2webp/model"
)

// ImageServer 处理图片信息
func ImageServer() {
	// 获取图片地址
	var imageObj []model.Image

	// 使用 Find 查询数据
	if err := global.DB.Table(global.Conf.FilterInfo.TableName).Select("id,code,image").Where(
		fmt.Sprintf("%s = ?", global.Conf.FilterInfo.ColumnName),
		global.Conf.BucketInfo.ApiId).Find(&imageObj).Error; err != nil {
		zap.L().Error("mysql query failed", zap.Error(err))
		return
	}

	if len(imageObj) == 0 {
		zap.S().Info("No data found")
		return
	}

	// 创建管道用于存储url
	urlChan := make(chan map[string]interface{}, 100) // 使用 interface{} 存储多种类型

	// 这里创建一个 WaitGroup 来等待所有的并发操作完成
	var wg sync.WaitGroup

	// 将所有的图片 URL 放入管道，包含 id
	go func() {
		for _, v := range imageObj {
			// 将图片信息和 id 放入管道中
			urlChan <- map[string]interface{}{
				"code":  v.Code,
				"image": v.Image,
				"id":    v.ID, // 将 id 放入管道
			}
		}
		// 关闭管道，表示没有更多数据
		close(urlChan)
	}()

	for i := 0; i < global.Conf.FilterInfo.WorkNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// 处理从管道中接收到的 URL
			for urlInfo := range urlChan {
				processImage(urlInfo)
			}
		}(i)
	}

	// 等待所有并发操作完成
	wg.Wait()

	// 所有并发任务完成后，可以做一些收尾工作
	zap.S().Info("All image processing completed.")
}

// 处理每个图片 URL 的函数
func processImage(urlInfo map[string]interface{}) {
	code := urlInfo["code"].(string)
	imageURL := urlInfo["image"].(string)
	id := urlInfo["id"].(int64) // 从管道中取出 id

	if len(imageURL) > 0 {
		downloadUrl := ""
		// 下载图片并转换为 WebP 格式
		if strings.HasPrefix(imageURL, "http") {
			downloadUrl = imageURL
		} else {
			downloadUrl = fmt.Sprintf("https://%s%s", global.Conf.BucketInfo.OldBucketName, imageURL)
		}

		if err := downloadAndConvertImage(downloadUrl, code, id); err != nil {
			zap.S().Errorf("Failed to process image for code %s: %v", code, err)
		}
	} else {
		zap.S().Warnf("Invalid image URL for code: %s", code)
	}
}

// downloadAndConvertImage 下载图片并将其转换为 WebP 格式
func downloadAndConvertImage(url, code string, id int64) error {
	// 获取图片格式
	format := strings.ToLower(filepath.Ext(url))
	// 下载图片
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// 检查返回状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch image: status code %d", resp.StatusCode)
	}

	var img image.Image
	// 根据图片扩展名选择解码器
	switch format {
	case ".png":
		img, err = png.Decode(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to decode PNG image: %v", err)
		}
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to decode JPEG image: %v", err)
		}
	case ".webp":
		img, err = webp.Decode(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to decode WebP image: %v", err)
		}
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
	// 将图片编码为 WebP 格式
	webpBytes, err := webp.EncodeRGBA(img, 80)
	if err != nil {
		return fmt.Errorf("failed to encode webp: %v", err)
	}
	parts := strings.Split(code, global.Conf.FilterInfo.SplitFilter)
	rename := fmt.Sprintf("GameID_%s_EN.webp", parts[len(parts)-1])
	//
	//// 上传图片到 S3
	path := global.Conf.BucketInfo.FilePath
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	uploadPath := filepath.Join(path, rename)
	uploadRepPath := strings.ReplaceAll(uploadPath, "\\", "/")
	_, err = global.S3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(global.Conf.BucketInfo.BucketName),
		Key:    aws.String(uploadRepPath), // 设置上传路径
		Body:   bytes.NewReader(webpBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// 更新数据库中的图片路径
	global.DB.Table(global.Conf.FilterInfo.TableName).Where("id = ?", id).Update("image", uploadRepPath)
	// 日志输出
	zap.S().Infof("Image successfully converted to WebP for code: %s and uploaded to S3", code)
	return nil
}
