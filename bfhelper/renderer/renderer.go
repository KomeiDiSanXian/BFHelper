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

// NewImageForImage 使用image.Image 创建图片
func NewImageForImage(im image.Image) *Image {
	return &Image{
		Canvas: gg.NewContextForImage(im),
	}
}

// NewImageWithXY 使用Image 创建位置在 (x, y) 的图片
func NewImageWithXY(img *Image, x, y int) *ImageWithXY {
	return &ImageWithXY{
		Image: img,
		X:     x,
		Y:     y,
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
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open")
	}
	defer f.Close()
	background, _, err := image.Decode(f)
	if err != nil {
		return errors.Wrap(err, "failed to decode")
	}
	backgroundHeight, backgroundWidth := float64(background.Bounds().Dy()), float64(background.Bounds().Dx())
	CanvasHeight, CanvasWidth := float64(i.Canvas.H()), float64(i.Canvas.W())
	if backgroundHeight/backgroundWidth < CanvasHeight/CanvasWidth {
		background = imgfactory.Size(background, int(backgroundWidth*CanvasHeight/backgroundHeight), int(backgroundHeight*CanvasHeight/backgroundHeight)).Image()
		i.Canvas.DrawImageAnchored(background, i.Canvas.W()/2, i.Canvas.H()/2, 0.5, 0.5)
		return nil
	}
	background = imgfactory.Size(background, int(backgroundWidth*CanvasWidth/backgroundWidth), int(backgroundHeight*CanvasWidth/backgroundWidth)).Image()
	i.Canvas.DrawImage(background, 0, 0)
	return nil
}

// DrawCopyright 绘制底部信息
func (i *Image) DrawCopyright(copyright string) {
	i.ParseFontFace(28)
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

// Blur 绘制强度为delta 的高斯模糊
func (i *Image) Blur(delta float64) image.Image {
	return imaging.Blur(i.Canvas.Image(), delta)
}

// Fillet 裁剪图片为圆角矩形
func (i *Image) Fillet(r float64) image.Image {
	return rendercard.Fillet(i.Canvas.Image(), r)
}

// ToImage 转换为image.Image
func (i *Image) ToImage() image.Image {
	return i.Canvas.Image()
}
