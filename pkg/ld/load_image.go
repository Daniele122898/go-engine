package ld

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func LoadImageData(path string) (*image.RGBA, error) {
	imgFile, err := os.Open(path)

	if err != nil {
		log.Printf("Couldn't find image at path: %s \n %v", path, err)
		return nil, err
	}
	defer imgFile.Close()

	img, format, err := image.Decode(imgFile)
	if err != nil {
		log.Printf("Couldn't decode image with format %s: \n %v", format, err)
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X * 4 {
		log.Printf("unsupported stride %d", rgba.Stride)
		return nil, err
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba, nil
}
