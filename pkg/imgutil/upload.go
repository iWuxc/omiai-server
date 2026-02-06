// Package imgutil 图片上传处理
// 统一输出 PNG 格式，自动压缩优化
package imgutil

import (
	"bytes"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// UploadProcessResult 上传处理结果
type UploadProcessResult struct {
	Data       *bytes.Buffer // 处理后的图片数据
	Format     string        // 输出格式（统一为png）
	OriginSize int           // 原始大小
	FinalSize  int           // 处理后大小
	Width      int           // 宽度
	Height     int           // 高度
}

// ProcessUpload 处理上传文件（统一输出PNG）
// 适用于业务场景：用户上传头像、图片等
func ProcessUpload(file multipart.File, header *multipart.FileHeader) (*UploadProcessResult, error) {
	// 根据原始文件大小选择不同的压缩策略
	opts := determineProcessOptions(header.Size)

	// 根据文件扩展名优化压缩格式
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == ".jpg" || ext == ".jpeg" {
		opts.TargetFormat = "jpeg"
		// 如果是 JPEG，适当放宽尺寸限制，因为 JPEG 压缩率更高
		if opts.MaxWidth < 1500 && opts.MaxWidth > 0 {
			opts.MaxWidth = 1500
			opts.MaxHeight = 1500
		}
	} else {
		opts.TargetFormat = "png"
	}

	result, err := ProcessFromMultipart(file, opts)
	if err != nil {
		return nil, err
	}

	return &UploadProcessResult{
		Data:       result.Data,
		Format:     result.Format,
		OriginSize: result.OriginSize,
		FinalSize:  result.FinalSize,
		Width:      result.Width,
		Height:     result.Height,
	}, nil
}

// ProcessUploadSimple 简化版处理（使用默认配置）
func ProcessUploadSimple(file multipart.File, header *multipart.FileHeader) (io.Reader, error) {
	result, err := ProcessUpload(file, header)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ProcessUploadWithLimit 带大小限制的处理
// maxSize: 最大文件大小（字节），0表示使用默认值(3MB)
// maxEdge: 最大边长像素，0表示使用默认值(1500)
func ProcessUploadWithLimit(file multipart.File, header *multipart.FileHeader, maxSize, maxEdge int) (*UploadProcessResult, error) {
	opts := &ProcessOptions{
		MaxWidth:    maxEdge,
		MaxHeight:   maxEdge,
		MaxFileSize: maxSize,
		AutoOrient:  true,
	}

	if maxSize == 0 {
		opts.MaxFileSize = Limit3MB
	}
	if maxEdge == 0 {
		opts.MaxWidth = MaxEdgeSize
		opts.MaxHeight = MaxEdgeSize
	}

	result, err := ProcessFromMultipart(file, opts)
	if err != nil {
		return nil, err
	}

	return &UploadProcessResult{
		Data:       result.Data,
		Format:     result.Format,
		OriginSize: result.OriginSize,
		FinalSize:  result.FinalSize,
		Width:      result.Width,
		Height:     result.Height,
	}, nil
}

// ProcessAvatar 处理头像上传（小尺寸）
// 输出：PNG，最大 500KB，尺寸 800x800
func ProcessAvatar(file multipart.File, header *multipart.FileHeader) (*UploadProcessResult, error) {
	opts := &ProcessOptions{
		MaxWidth:    800,
		MaxHeight:   800,
		MaxFileSize: Limit500KB,
		AutoOrient:  true,
	}

	result, err := ProcessFromMultipart(file, opts)
	if err != nil {
		return nil, err
	}

	return &UploadProcessResult{
		Data:       result.Data,
		Format:     result.Format,
		OriginSize: result.OriginSize,
		FinalSize:  result.FinalSize,
		Width:      result.Width,
		Height:     result.Height,
	}, nil
}

// ProcessPhoto 处理照片上传（中等尺寸）
// 输出：PNG，最大 2MB，尺寸 1500x1500
func ProcessPhoto(file multipart.File, header *multipart.FileHeader) (*UploadProcessResult, error) {
	opts := &ProcessOptions{
		MaxWidth:    MaxEdgeSize,
		MaxHeight:   MaxEdgeSize,
		MaxFileSize: 2 * MB,
		AutoOrient:  true,
	}

	result, err := ProcessFromMultipart(file, opts)
	if err != nil {
		return nil, err
	}

	return &UploadProcessResult{
		Data:       result.Data,
		Format:     result.Format,
		OriginSize: result.OriginSize,
		FinalSize:  result.FinalSize,
		Width:      result.Width,
		Height:     result.Height,
	}, nil
}

// ProcessGallery 处理图库/相册图片（大尺寸）
// 输出：PNG，最大 5MB，尺寸 2000x2000
func ProcessGallery(file multipart.File, header *multipart.FileHeader) (*UploadProcessResult, error) {
	opts := &ProcessOptions{
		MaxWidth:    2000,
		MaxHeight:   2000,
		MaxFileSize: Limit5MB,
		AutoOrient:  true,
	}

	result, err := ProcessFromMultipart(file, opts)
	if err != nil {
		return nil, err
	}

	return &UploadProcessResult{
		Data:       result.Data,
		Format:     result.Format,
		OriginSize: result.OriginSize,
		FinalSize:  result.FinalSize,
		Width:      result.Width,
		Height:     result.Height,
	}, nil
}

// determineProcessOptions 根据原始大小确定处理选项
func determineProcessOptions(originSize int64) *ProcessOptions {
	switch {
	case originSize < Limit200KB:
		// 小图：不压缩尺寸，仅转换格式
		return &ProcessOptions{
			MaxWidth:    0,
			MaxHeight:   0,
			MaxFileSize: 0,
			AutoOrient:  true,
		}
	case originSize < Limit500KB:
		// 中等图：限制最大边1500
		return &ProcessOptions{
			MaxWidth:    MaxEdgeSize,
			MaxHeight:   MaxEdgeSize,
			MaxFileSize: Limit500KB,
			AutoOrient:  true,
		}
	case originSize < Limit3MB:
		// 大图：限制最大边1500，最大3MB
		return &ProcessOptions{
			MaxWidth:    MaxEdgeSize,
			MaxHeight:   MaxEdgeSize,
			MaxFileSize: Limit3MB,
			AutoOrient:  true,
		}
	default:
		// 超大图（>3MB）：严格限制在 1MB 以内
		return &ProcessOptions{
			MaxWidth:    1200, // 降低最大边长起始值
			MaxHeight:   1200,
			MaxFileSize: Limit1MB, // 目标大小 1MB
			AutoOrient:  true,
		}
	}
}
