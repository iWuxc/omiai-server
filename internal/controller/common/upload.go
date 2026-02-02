package common

import (
	"fmt"
	"omiai-server/pkg/response"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iWuxc/go-wit/log"
)

const (
	MaxUploadSize = 5 * 1024 * 1024 // 5MB
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
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
		response.ErrorResponse(ctx, response.ParamsCommonError, "文件大小不能超过5MB")
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

	// 4. Generate New Filename/Key
	key := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)

	// 5. Save File via Storage Driver
	url, err := c.storage.Put(ctx, key, src)
	if err != nil {
		log.Errorf("Storage put failed: %v", err)
		response.ErrorResponse(ctx, response.FuncCommonError, "文件保存失败")
		return
	}

	// 6. Return URL
	response.SuccessResponse(ctx, "上传成功", map[string]string{
		"url": url,
	})
}
