package graphics

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"log"
	"unsafe"
	"voxel/pkg/ld"
)

type Texture2D struct {
	Id uint32
}

func NewTexture2D(texturePath string, wrapS, wrapT, minFilter, magFilter int32) (*Texture2D, error) {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// Set the texture wrapping/filtering options
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, wrapT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, minFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, magFilter)
	// Load the actual image data
	rgba, err := ld.LoadImageData(texturePath)
	if err != nil {
		log.Printf("failed to load texture:\n %v", err)
		return nil, err
	}
	// Generate the texture
	size := rgba.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&rgba.Pix[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Unbind texture so we can create and modify other ones
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return &Texture2D{
		Id: texture,
	}, nil
}

// Use/activate the texture
func (t *Texture2D) Use() {
	gl.BindTexture(gl.TEXTURE_2D, t.Id)
}