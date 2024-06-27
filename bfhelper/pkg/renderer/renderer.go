// Package renderer 图片和卡片绘制api
package renderer

import (
	"image"
	"io"
	"os"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"github.com/FloatTech/rendercard"
	"github.com/pkg/errors"
)

// Image 图片绘制
type Image struct {
	Canvas   *gg.Context
	FontByte []byte
}

// ImageWithXY 含有坐标信息的图片
type ImageWithXY struct {
	*Image
	X, Y int
}

// NewImage 创建一个宽、高分别为width、height 的空白图片
func NewImage(width, height int) *Image {
	return &Image{
		Canvas: gg.NewContext(width, height),
	}
}

// NewImageByImage 使用image.Image 创建图片
func NewImageByImage(im image.Image) *Image {
	return &Image{
		Canvas: gg.NewContextForImage(im),
	}
}

// NewImageByPath 使用path下的图片创建图片
func NewImageByPath(path string) *Image {
	img, _ := ReadFromPath(path)
	return NewImageByImage(img)
}

// NewImageWithXY 使用Image 创建位置在 (x, y) 的图片
func NewImageWithXY(img *Image, x, y int) *ImageWithXY {
	return &ImageWithXY{
		Image: img,
		X:     x,
		Y:     y,
	}
}

// ReadFromPath 读取路径为path 的图片
func ReadFromPath(path string) (image.Image, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// WithXY 转换Imgage 为ImageWithXY
func (i *Image) WithXY(x, y int) *ImageWithXY {
	return &ImageWithXY{
		i, x, y,
	}
}

// LoadFont 加载字体
func (i *Image) LoadFont(fontPath string) error {
	f, err := os.Open(fontPath)
	if err != nil {
		return errors.Wrap(err, "failed to open")
	}
	defer f.Close()
	i.FontByte, err = io.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "failed to read font")
	}
	return nil
}

// ParseFontFace 使用Image.FontByte 字体，大小为point
func (i *Image) ParseFontFace(point float64) error {
	err := i.Canvas.ParseFontFace(i.FontByte, point)
	if err != nil {
		return errors.Wrap(err, "Error parsing font face")
	}
	return nil
}

// ToBytes 将图片转换为最大为4MB, 编码质量为70 的jpeg []byte
func (i *Image) ToBytes() ([]byte, error) {
	return imgfactory.ToBytes(i.Canvas.Image())
}

// DrawBackground 绘制背景
func (i *Image) DrawBackground(path string) error {
	background, err := ReadFromPath(path)
	if err != nil {
		return errors.Wrap(err, "failed to decode")
	}
	bgimg := NewImageWithXY(NewImageByImage(background), 0, 0)
	bgimg.ScaleToSize(i.Canvas.W(), i.Canvas.H())
	i.DrawImages([]*ImageWithXY{bgimg})
	return nil
}

// DrawCopyright 绘制底部信息
//
// copyright: 绘制的文字  point: 字体大小
func (i *Image) DrawCopyright(copyright string, point float64) {
	_ = i.ParseFontFace(point)
	i.Canvas.SetRGBA255(0, 0, 0, 255)
	i.Canvas.DrawStringAnchored(copyright, float64(i.Canvas.W()/2+3), float64(i.Canvas.H()-70/2+3), 0.5, 0.5)
	i.Canvas.SetRGBA255(255, 255, 255, 255)
	i.Canvas.DrawStringAnchored(copyright, float64(i.Canvas.W()/2), float64(i.Canvas.H()-70/2), 0.5, 0.5)
}

// DrawImages 将imgs 绘制到Image 结构体上的指定坐标
func (i *Image) DrawImages(imgs []*ImageWithXY) *Image {
	for _, img := range imgs {
		if img != nil {
			i.Canvas.DrawImage(img.ToImage(), img.X, img.Y)
		}
	}
	return i
}

// // Blur 绘制强度为delta 的高斯模糊
// func (i *Image) Blur(delta float64) image.Image {
// 	return imaging.Blur(i.Canvas.Image(), delta)
// }

// Fillet 裁剪图片为圆角矩形
func (i *Image) Fillet(r float64) image.Image {
	return rendercard.Fillet(i.Canvas.Image(), r)
}

// ToImage 转换为image.Image
func (i *Image) ToImage() image.Image {
	return i.Canvas.Image()
}

// ToBase64 转换图片为b64
func (i *Image) ToBase64() ([]byte, error) {
	return imgfactory.ToBase64(i.Canvas.Image())
}

// ScaleByPercent 百分比缩放
func (i *Image) ScaleByPercent(factor float64) *Image {
	width, height := i.Canvas.W(), i.Canvas.H()
	newWidth, newHeight := float64(width)*factor, float64(height)*factor
	newCanvas := gg.NewContext(int(newWidth), int(newHeight))
	newCanvas.Scale(factor, factor)
	newCanvas.DrawImage(i.ToImage(), 0, 0)
	i.Canvas = newCanvas
	return i
}

// ScaleToSize 缩放到指定宽高
func (i *Image) ScaleToSize(targetWidth, targetHeight int) *Image {
	newCanvas := gg.NewContext(targetWidth, targetHeight)
	scaleX := float64(targetWidth) / float64(i.Canvas.W())
	scaleY := float64(targetHeight) / float64(i.Canvas.H())
	newCanvas.Scale(scaleX, scaleY)
	newCanvas.DrawImage(i.ToImage(), 0, 0)
	i.Canvas = newCanvas
	return i
}
