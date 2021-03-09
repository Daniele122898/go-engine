package graphics

import (
	"errors"
	"github.com/go-gl/gl/v4.6-core/gl"
	"io/ioutil"
	"log"
	"voxel/pkg/su"
)

type Shader struct {
	Id uint32
}

// Read and build the shader
func NewShader(vertexPath string, fragmentPath string) (*Shader, error) {
	// Retrieve the vertex and fragment source code from their respective paths
	vertSource, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		log.Printf("couldn't find vertex shader on path %s :\n %v", vertexPath, err)
		return nil, err
	}
	// null terminate the "string"
	vertSource = append(vertSource, 0x00)
	fragSource, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		log.Printf("couldn't find fragment shader on path %s :\n %v", fragmentPath, err)
		return nil, err
	}
	// null terminate the "string"
	fragSource = append(fragSource, 0x00)

	// compile the shaders
	var vertex, fragment uint32
	var success int32
	infoLog := [512]uint8{}

	// Compile vertex shader
	vertex = gl.CreateShader(gl.VERTEX_SHADER)
	vp, free := gl.Strs(string(vertSource))
	gl.ShaderSource(vertex, 1, vp, nil)
	gl.CompileShader(vertex)
	gl.GetShaderiv(vertex, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(vertex, 512, nil, &infoLog[0])
		errStr := gl.GoStr(&infoLog[0])
		log.Printf("failed to compile vertex shader:\n %s", errStr)
		return nil, errors.New(errStr)
	}
	free()

	// Compile fragment shader
	fragment = gl.CreateShader(gl.FRAGMENT_SHADER)
	fp, free := gl.Strs(string(fragSource))
	gl.ShaderSource(fragment, 1, fp, nil)
	gl.CompileShader(fragment)
	gl.GetShaderiv(fragment, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(fragment, 512, nil, &infoLog[0])
		errStr := gl.GoStr(&infoLog[0])
		log.Printf("failed to compile fragment shader:\n %s", errStr)
		return nil, errors.New(errStr)
	}
	free()

	// Create shader program
	shaderProg := gl.CreateProgram()
	gl.AttachShader(shaderProg, vertex)
	gl.AttachShader(shaderProg, fragment)
	gl.LinkProgram(shaderProg)
	gl.GetShaderiv(shaderProg, gl.LINK_STATUS, &success)
	if success == 0 {
		gl.GetProgramInfoLog(shaderProg, 512, nil, &infoLog[0])
		errStr := gl.GoStr(&infoLog[0])
		log.Printf("shader programm linking failed %s", errStr)
		return nil, errors.New(errStr)
	}
	// Now delete the shader objects as we already linked them and dont need them
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	return &Shader{
		Id: shaderProg,
	}, nil
}

func (s *Shader) Delete()  {
	gl.DeleteProgram(s.Id)
}

// Use/activate the shader
func (s *Shader) Use() {
	gl.UseProgram(s.Id)
}

// Uniform utility functions

func (s *Shader) SetBool(name string, value bool) {
	var ibool int32
	if value {
		ibool = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(s.Id, su.CStr(name)), ibool)
}

func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.Id, su.CStr(name)), value)
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.Id, su.CStr(name)), value)
}


