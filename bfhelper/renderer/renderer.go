// Package renderer 图片和卡片绘制api
package renderer

import (
	"image"
	"io"
	"os"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"github.com/FloatTech/rendercard"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

// Image 图片绘制
type Image struct {
	Canvas   *gg.Context
	FontByte []byte
	Err      error
}

// ImageWithXY 含有坐标信息的图片
type ImageWithXY struct {
	*Image
	X, Y int
}

// NewImage 创建一个宽、高分别为width、height 的图片
func NewImage(width, height int) *Image {
	return &Image{
		Canvas: gg.NewContext(width, height),
	}
}

// NewImageWithXY 创建一个宽、高分别为width、height 位置在(x,y)的图片
func NewImageWithXY(width, height, x, y int) *ImageWithXY {
	return &ImageWithXY{
		Image: NewImage(width, height),
		X:     x,
		Y:     y,
	}
}

// LoadFont 加载字体
func (i *Image) LoadFont(fontPath string) *Image {
	f, err := os.Open(fontPath)
	if err != nil {
		i.Err = errors.Wrap(err, "failed to open")
		return i
	}
	defer f.Close()
	i.FontByte, err = io.ReadAll(f)
	if err != nil {
		i.Err = errors.Wrap(err, "failed to read font")
		return i
	}
	return i
}

// ParseFontFace 使用Image.FontByte 字体，大小为point
func (i *Image) ParseFontFace(point float64) *Image {
	err := i.Canvas.ParseFontFace(i.FontByte, point)
	if err != nil {
		i.Err = errors.Wrap(err, "Error parsing font face")
		return i
	}
	return i
}

// ToBytes 将图片转换为最大为4MB, 编码质量为70 的jpeg []byte
func (i *Image) ToBytes() ([]byte, error) {
	return imgfactory.ToBytes(i.Canvas.Image())
}

// DrawBackground 绘制背景
func (i *Image) DrawBackground(path string) *Image {
	f, err := os.Open(path)
	if err != nil {
		i.Err = errors.Wrap(err, "failed to open")
		return i
	}
	defer f.Close()
	background, _, err := image.Decode(f)
	if err != nil {
		i.Err = errors.Wrap(err, "failed to decode")
		return i
	}
	backgroundHeight, backgroundWidth := float64(background.Bounds().Dy()), float64(background.Bounds().Dx())
	CanvasHeight, CanvasWidth := float64(i.Canvas.H()), float64(i.Canvas.W())
	if backgroundHeight/backgroundWidth < CanvasHeight/CanvasWidth {
		background = imgfactory.Size(background, int(backgroundWidth*CanvasHeight/backgroundHeight), int(backgroundHeight*CanvasHeight/backgroundHeight)).Image()
		i.Canvas.DrawImageAnchored(background, i.Canvas.W()/2, i.Canvas.H()/2, 0.5, 0.5)
		return i
	}
	background = imgfactory.Size(background, int(backgroundWidth*CanvasWidth/backgroundWidth), int(backgroundHeight*CanvasWidth/backgroundWidth)).Image()
	i.Canvas.DrawImage(background, 0, 0)
	return i
}

// DrawCopyright 绘制底部信息
func (i *Image) DrawCopyright() *Image {
	i.ParseFontFace(28).Canvas.SetRGBA255(0, 0, 0, 255)
	copyright := "Created By RemiliaBot Forked From ZeroBot-Plugin"
	i.Canvas.DrawStringAnchored(copyright, float64(i.Canvas.W()/2+3), float64(i.Canvas.H()-70/2+3), 0.5, 0.5)
	i.Canvas.SetRGBA255(255, 255, 255, 255)
	i.Canvas.DrawStringAnchored(copyright, float64(i.Canvas.W()/2), float64(i.Canvas.H()-70/2), 0.5, 0.5)
	return i
}

// DrawImages 将imgs 绘制到Image 结构体上的指定坐标
func (i *Image) DrawImages(imgs []ImageWithXY) *Image {
	for _, img := range imgs {
		i.Canvas.DrawImage(img.Canvas.Image(), img.X, img.Y)
	}
	return i
}

// Blur 绘制强度为delta 的高斯模糊
func (i *Image) Blur(delta float64) image.Image {
	return imaging.Blur(i.Canvas.Image(), delta)
}

// Fillet 裁剪图片为圆角矩形
func (i *Image) Fillet(r float64) image.Image {
	return rendercard.Fillet(i.Canvas.Image(), r)
}
