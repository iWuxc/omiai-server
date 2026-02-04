// Package imgutil 提供图片处理工具函数
// 统一输出为 PNG 格式，支持压缩、裁剪、缩放等功能
package imgutil

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"

	"github.com/disintegration/imaging"
)

// 业务阈值（字节）
const (
	KB           = 1024
	MB           = 1024 * 1024
	Limit20MB    = 20 * MB
	Limit5MB     = 5 * MB
	Limit3MB     = 3 * MB
	Limit500KB   = 500 * KB
	Limit200KB   = 200 * KB
	MinPixelSize = 32  // 最小边长像素
	MaxEdgeSize  = 1500 // 最大边长像素
)

var (
	// ErrOverSingleLimit 超过单张图片上限 20MB
	ErrOverSingleLimit = errors.New("图片超过单张上限 20MB")
	// ErrTooSmall 图片尺寸过小
	ErrTooSmall = errors.New("图片尺寸不能小于 32x32 像素")
	// ErrDecodeFailed 图片解码失败
	ErrDecodeFailed = errors.New("图片解码失败")
	// ErrEncodeFailed 图片编码失败
	ErrEncodeFailed = errors.New("图片编码失败")
)

// ProcessOptions 图片处理选项
type ProcessOptions struct {
	MaxWidth    int  // 最大宽度，0表示不限制
	MaxHeight   int  // 最大高度，0表示不限制
	Quality     int  // 压缩质量 1-100，PNG忽略此参数
	MaxFileSize int  // 最大文件大小（字节），0表示不限制
	MinFileSize int  // 最小文件大小（字节），0表示不限制
	AutoOrient  bool // 自动校正图片方向
}

// DefaultOptions 默认处理选项（统一输出PNG）
var DefaultOptions = &ProcessOptions{
	MaxWidth:    MaxEdgeSize,
	MaxHeight:   MaxEdgeSize,
	Quality:     85,
	MaxFileSize: Limit3MB,
	MinFileSize: 0,
	AutoOrient:  true,
}

// ProcessResult 处理结果
type ProcessResult struct {
	Data       *bytes.Buffer // 处理后的图片数据
	Width      int           // 图片宽度
	Height     int           // 图片高度
	Format     string        // 输出格式（统一为png）
	OriginSize int           // 原始文件大小
	FinalSize  int           // 处理后文件大小
}

// ProcessFromMultipart 处理 multipart 上传的图片
// 统一输出为 PNG 格式
func ProcessFromMultipart(file multipart.File, opts *ProcessOptions) (*ProcessResult, error) {
	if opts == nil {
		opts = DefaultOptions
	}

	// 读取文件数据
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return ProcessFromBytes(data, opts)
}

// ProcessFromBytes 从字节数组处理图片
// 统一输出为 PNG 格式
func ProcessFromBytes(data []byte, opts *ProcessOptions) (*ProcessResult, error) {
	if opts == nil {
		opts = DefaultOptions
	}

	originSize := len(data)

	// 检查原始大小限制
	if originSize > Limit20MB {
		return nil, ErrOverSingleLimit
	}

	// 解码图片
	decodeOpts := imaging.AutoOrientation(false)
	if opts.AutoOrient {
		decodeOpts = imaging.AutoOrientation(true)
	}

	img, err := imaging.Decode(bytes.NewReader(data), decodeOpts)
	if err != nil {
		return nil, ErrDecodeFailed
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 尺寸下限校验
	if width < MinPixelSize || height < MinPixelSize {
		return nil, ErrTooSmall
	}

	// 智能压缩处理
	result, err := smartCompress(img, opts, originSize)
	if err != nil {
		return nil, err
	}

	result.OriginSize = originSize
	return result, nil
}

// smartCompress 智能压缩策略
// 统一输出 PNG 格式
func smartCompress(img image.Image, opts *ProcessOptions, originSize int) (*ProcessResult, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	maxEdge := max(width, height)

	// 第一阶段：尺寸缩放
	if opts.MaxWidth > 0 || opts.MaxHeight > 0 {
		maxW, maxH := opts.MaxWidth, opts.MaxHeight
		if maxW == 0 {
			maxW = MaxEdgeSize
		}
		if maxH == 0 {
			maxH = MaxEdgeSize
		}

		// 如果图片尺寸超过限制，进行等比缩放
		if maxEdge > maxW && maxEdge > maxH {
			if width > height {
				img = imaging.Resize(img, maxW, 0, imaging.Lanczos)
			} else {
				img = imaging.Resize(img, 0, maxH, imaging.Lanczos)
			}
		}
	}

	// 第二阶段：PNG 编码（质量优先）
	// PNG 是无损格式，通过调整压缩级别控制大小
	buf := new(bytes.Buffer)
	encoder := &png.Encoder{
		CompressionLevel: png.BestCompression, // 最高压缩级别
	}

	if err := encoder.Encode(buf, img); err != nil {
		return nil, ErrEncodeFailed
	}

	// 第三阶段：如果文件仍然过大，进行降级处理
	if opts.MaxFileSize > 0 && buf.Len() > opts.MaxFileSize {
		buf = downgradeCompress(img, opts.MaxFileSize)
	}

	// 更新尺寸信息
	bounds = img.Bounds()

	return &ProcessResult{
		Data:      buf,
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		Format:    "png",
		FinalSize: buf.Len(),
	}, nil
}

// downgradeCompress 降级压缩策略
// 当 PNG 过大时，通过降低尺寸来减小文件大小
func downgradeCompress(img image.Image, maxSize int) *bytes.Buffer {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	maxEdge := max(width, height)

	// 逐步降低尺寸直到满足大小要求
	sizes := []int{1200, 1000, 800, 600}
	for _, targetSize := range sizes {
		if maxEdge <= targetSize {
			continue
		}

		var resized image.Image
		if width > height {
			resized = imaging.Resize(img, targetSize, 0, imaging.Lanczos)
		} else {
			resized = imaging.Resize(img, 0, targetSize, imaging.Lanczos)
		}

		buf := new(bytes.Buffer)
		encoder := &png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		if err := encoder.Encode(buf, resized); err == nil && buf.Len() <= maxSize {
			return buf
		}
	}

	// 如果仍不满足，返回最后一次尝试的结果
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf
}

// ValidateProductPixelArea 验证产品图像有效区域（白色像素）
// 适用于黑底白图，通过统计白色像素判断产品有效区域
func ValidateProductPixelArea(img image.Image, minPixelArea int) (bool, int, int, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w == 0 || h == 0 {
		return false, 0, 0, errors.New("图片尺寸无效")
	}

	totalArea := w * h
	validPixels := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			// 白色像素判定（容差范围）
			if rgba.R >= 250 && rgba.G >= 250 && rgba.B >= 250 {
				validPixels++
			}
		}
	}

	return validPixels >= minPixelArea, validPixels, totalArea, nil
}

// Thumbnail 生成缩略图
func Thumbnail(img image.Image, width, height int) image.Image {
	return imaging.Thumbnail(img, width, height, imaging.Lanczos)
}

// Resize 等比缩放图片
func Resize(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// Crop 裁剪图片
func Crop(img image.Image, rect image.Rectangle) image.Image {
	return imaging.Crop(img, rect)
}

// CropCenter 从中心裁剪图片
func CropCenter(img image.Image, width, height int) image.Image {
	return imaging.CropCenter(img, width, height)
}

// Fill 填充到指定尺寸（可能变形）
func Fill(img image.Image, width, height int) image.Image {
	return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
}

// Fit 适应到指定尺寸（保持比例）
func Fit(img image.Image, width, height int) image.Image {
	return imaging.Fit(img, width, height, imaging.Lanczos)
}

// Grayscale 转换为灰度图
func Grayscale(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

// Blur 高斯模糊
func Blur(img image.Image, sigma float64) image.Image {
	return imaging.Blur(img, sigma)
}

// Sharpen 锐化
func Sharpen(img image.Image, sigma float64) image.Image {
	return imaging.Sharpen(img, sigma)
}

// AdjustBrightness 调整亮度 (-100 到 100)
func AdjustBrightness(img image.Image, percentage float64) image.Image {
	return imaging.AdjustBrightness(img, percentage)
}

// AdjustContrast 调整对比度 (-100 到 100)
func AdjustContrast(img image.Image, percentage float64) image.Image {
	return imaging.AdjustContrast(img, percentage)
}

// AdjustSaturation 调整饱和度 (-100 到 100)
func AdjustSaturation(img image.Image, percentage float64) image.Image {
	return imaging.AdjustSaturation(img, percentage)
}

// EncodeToPNG 编码为 PNG 字节数组
func EncodeToPNG(img image.Image, level png.CompressionLevel) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	encoder := &png.Encoder{
		CompressionLevel: level,
	}
	if err := encoder.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf, nil
}

// GetImageInfo 获取图片信息（不处理）
func GetImageInfo(data []byte) (width, height int, format string, err error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return 0, 0, "", err
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), format, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
