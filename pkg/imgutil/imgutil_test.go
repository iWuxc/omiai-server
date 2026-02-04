package imgutil

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// createTestImage 创建测试图片
func createTestImage(width, height int) *bytes.Buffer {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// 填充渐变色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: 128,
				A: 255,
			})
		}
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf
}

func TestProcessFromBytes(t *testing.T) {
	// 创建测试图片数据
	testData := createTestImage(2000, 1500)

	tests := []struct {
		name    string
		data    []byte
		opts    *ProcessOptions
		wantErr bool
	}{
		{
			name: "normal image with default options",
			data: testData.Bytes(),
			opts: DefaultOptions,
		},
		{
			name: "small image no resize",
			data: createTestImage(500, 400).Bytes(),
			opts: &ProcessOptions{
				MaxWidth:    1500,
				MaxHeight:   1500,
				AutoOrient:  true,
			},
		},
		{
			name: "oversize image",
			data: make([]byte, Limit20MB+1),
			opts: DefaultOptions,
			wantErr: true,
		},
		{
			name: "too small image",
			data: createTestImage(10, 10).Bytes(),
			opts: DefaultOptions,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessFromBytes(tt.data, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// 验证输出格式
			if result.Format != "png" {
				t.Errorf("Expected PNG format, got %s", result.Format)
			}

			// 验证尺寸不超过限制
			if result.Width > MaxEdgeSize || result.Height > MaxEdgeSize {
				t.Errorf("Image dimensions exceed limit: %dx%d", result.Width, result.Height)
			}

			// 验证文件大小不超过限制
			if result.FinalSize > Limit3MB {
				t.Errorf("File size exceeds limit: %d bytes", result.FinalSize)
			}

			t.Logf("Processed: %dx%d, %d bytes (from %d bytes)",
				result.Width, result.Height, result.FinalSize, result.OriginSize)
		})
	}
}

func TestResize(t *testing.T) {
	testData := createTestImage(1000, 800)
	img, _, err := image.Decode(bytes.NewReader(testData.Bytes()))
	if err != nil {
		t.Fatalf("Failed to decode test image: %v", err)
	}

	resized := Resize(img, 500, 0)
	bounds := resized.Bounds()
	
	if bounds.Dx() != 500 {
		t.Errorf("Expected width 500, got %d", bounds.Dx())
	}
	if bounds.Dy() != 400 {
		t.Errorf("Expected height 400, got %d", bounds.Dy())
	}
}

func TestThumbnail(t *testing.T) {
	testData := createTestImage(1000, 800)
	img, _, err := image.Decode(bytes.NewReader(testData.Bytes()))
	if err != nil {
		t.Fatalf("Failed to decode test image: %v", err)
	}

	thumb := Thumbnail(img, 100, 100)
	bounds := thumb.Bounds()
	
	if bounds.Dx() > 100 {
		t.Errorf("Thumbnail width exceeds limit: %d", bounds.Dx())
	}
	if bounds.Dy() > 100 {
		t.Errorf("Thumbnail height exceeds limit: %d", bounds.Dy())
	}
}

func TestGrayscale(t *testing.T) {
	testData := createTestImage(100, 100)
	img, _, err := image.Decode(bytes.NewReader(testData.Bytes()))
	if err != nil {
		t.Fatalf("Failed to decode test image: %v", err)
	}

	gray := Grayscale(img)
	bounds := gray.Bounds()
	
	// 检查中心点是否为灰度
	center := bounds.Max.X / 2
	rgba := gray.At(center, center)
	c := color.RGBAModel.Convert(rgba).(color.RGBA)
	
	// 灰度图的 R=G=B
	if c.R != c.G || c.G != c.B {
		t.Error("Image is not grayscale")
	}
}

func TestGetImageInfo(t *testing.T) {
	testData := createTestImage(800, 600)
	
	width, height, format, err := GetImageInfo(testData.Bytes())
	if err != nil {
		t.Fatalf("GetImageInfo failed: %v", err)
	}
	
	if width != 800 {
		t.Errorf("Expected width 800, got %d", width)
	}
	if height != 600 {
		t.Errorf("Expected height 600, got %d", height)
	}
	if format != "png" {
		t.Errorf("Expected format png, got %s", format)
	}
}

func BenchmarkProcessFromBytes(b *testing.B) {
	testData := createTestImage(2000, 1500)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ProcessFromBytes(testData.Bytes(), DefaultOptions)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
