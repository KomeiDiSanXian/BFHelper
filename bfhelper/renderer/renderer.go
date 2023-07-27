// Package renderer 图片和卡片绘制api
package renderer

import (
	"image"
	"io"
	"os"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"github.com/pkg/errors"
)

// Image 图片绘制
type Image struct {
	canvas   *gg.Context
	FontByte []byte
	Err      error
}

// ImageWithXY 含有坐标信息的图片
type ImageWithXY struct {
	Image
	X, Y int
}

// NewImage 创建一个宽、高分别为width、height 的图片
func NewImage(width, height int) *Image {
	return &Image{
		canvas: gg.NewContext(width, height),
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
	err := i.canvas.ParseFontFace(i.FontByte, point)
	if err != nil {
		i.Err = errors.Wrap(err, "Error parsing font face")
		return i
	}
	return i
}

// ToBytes 将图片转换为最大为4MB, 编码质量为70 的jpeg []byte
func (i *Image) ToBytes() ([]byte, error) {
	return imgfactory.ToBytes(i.canvas.Image())
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
	canvasHeight, canvasWidth := float64(i.canvas.H()), float64(i.canvas.W())
	if backgroundHeight/backgroundWidth < canvasHeight/canvasWidth {
		background = imgfactory.Size(background, int(backgroundWidth*canvasHeight/backgroundHeight), int(backgroundHeight*canvasHeight/backgroundHeight)).Image()
		i.canvas.DrawImageAnchored(background, i.canvas.W()/2, i.canvas.H()/2, 0.5, 0.5)
		return i
	}
	background = imgfactory.Size(background, int(backgroundWidth*canvasWidth/backgroundWidth), int(backgroundHeight*canvasWidth/backgroundWidth)).Image()
	i.canvas.DrawImage(background, 0, 0)
	return i
}

// DrawCopyright 绘制底部信息
func (i *Image) DrawCopyright() *Image {
	i.ParseFontFace(28).canvas.SetRGBA255(0, 0, 0, 255)
	copyright := "Created By RemiliaBot Forked From ZeroBot-Plugin"
	i.canvas.DrawStringAnchored(copyright, float64(i.canvas.W()/2+3), float64(i.canvas.H()-70/2+3), 0.5, 0.5)
	i.canvas.SetRGBA255(255, 255, 255, 255)
	i.canvas.DrawStringAnchored(copyright, float64(i.canvas.W()/2), float64(i.canvas.H()-70/2), 0.5, 0.5)
	return i
}

// DrawImages 将imgs 绘制到Image 结构体上的指定坐标
func (i *Image) DrawImages(imgs []ImageWithXY) *Image {
	for _, img := range imgs {
		i.canvas.DrawImage(img.canvas.Image(), img.X, img.Y)
	}
	return i
}
