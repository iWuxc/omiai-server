package common

import (
	"fmt"
	"io"
	"mime"
	"omiai-server/pkg/imgutil"
	"omiai-server/pkg/response"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iWuxc/go-wit/log"
)

const (
	MaxUploadSize = 50 * 1024 * 1024 // 50MB
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

var ImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// GetContentType 根据文件扩展名获取 Content-Type
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// 默认 Content-Type
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = "application/octet-stream"
		}
	}
	return contentType
}

// IsImageFile 判断是否为图片文件
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ImageExtensions[ext]
}

func (c *Controller) Upload(ctx *gin.Context) {
	log.Infof("Start processing upload request")
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Errorf("Get form file failed: %v", err)
		response.ErrorResponse(ctx, response.ParamsCommonError, "上传文件不能为空")
		return
	}

	log.Infof("Receiving file: %s, size: %d", file.Filename, file.Size)

	// 1. Validate File Size
	if file.Size > MaxUploadSize {
		response.ErrorResponse(ctx, response.ParamsCommonError, "文件大小不能超过50MB")
		return
	}

	// 2. Validate File Extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedExtensions[ext] {
		response.ErrorResponse(ctx, response.ParamsCommonError, "不支持的文件格式")
		return
	}

	// 3. Open File
	src, err := file.Open()
	if err != nil {
		log.Errorf("Open upload file failed: %v", err)
		response.ErrorResponse(ctx, response.FuncCommonError, "文件打开失败")
		return
	}
	defer src.Close()

	// 4. Process Image if it's an image file
	var uploadReader io.Reader = src
	var finalExt = ext
	var contentType string

	if IsImageFile(file.Filename) {
		log.Infof("Processing image: %s", file.Filename)

		// 处理图片（统一输出 PNG）
		result, err := imgutil.ProcessUpload(src, file)
		if err != nil {
			log.Errorf("Image processing failed: %v", err)
			response.ErrorResponse(ctx, response.FuncCommonError, "图片处理失败: "+err.Error())
			return
		}

		log.Infof("Image processed: %d bytes -> %d bytes, dimensions: %dx%d, format: %s",
			result.OriginSize, result.FinalSize, result.Width, result.Height, result.Format)

		uploadReader = result.Data

		// 根据处理结果设置最终格式
		if result.Format == "jpeg" || result.Format == "jpg" {
			finalExt = ".jpg"
			contentType = "image/jpeg"
		} else {
			finalExt = ".png"
			contentType = "image/png"
		}
	} else {
		contentType = GetContentType(file.Filename)
	}

	// 5. Generate New Filename/Key
	key := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("20060102"), uuid.New().String(), finalExt)

	log.Infof("Uploading file: %s, Content-Type: %s", key, contentType)

	// 6. Save File via Storage Driver
	url, err := c.storage.Put(ctx, key, uploadReader, contentType)
	if err != nil {
		log.Errorf("Storage put failed: %v", err)
		response.ErrorResponse(ctx, response.FuncCommonError, "文件保存失败")
		return
	}

	// 7. Return URL
	response.SuccessResponse(ctx, "上传成功", map[string]string{
		"url": url,
	})
}
