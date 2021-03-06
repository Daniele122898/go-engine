package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"unsafe"
)

// Called first
func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var (
	vertexShaderSource = "#version 460 core\n\nlayout (location = 0) in vec3 aPos;\n\nvoid main() {\n    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n}\000"
)

func main() {
	log.Println("Starting voxel engine")

	err := glfw.Init()
	if err != nil {
		log.Fatal("failed to initialize glfw:\n", err)
	}
	// Add version hints
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// Only needed for mac but we dont care so just keeping it here for future reference
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	defer glfw.Terminate()

	window, err := glfw.CreateWindow(1280, 720, "VoxelEngine", nil, nil)
	if err != nil {
		log.Fatal("failed to create glfw window:\n", err)
	}

	window.MakeContextCurrent()

	// load and init glad/gl
	if err = gl.Init(); err != nil {
		log.Fatal("failed to initialize glad/gl:\n", err)
	}

	gl.Viewport(0,0,1280,720)
	// Register callback in case user changes the window size so gl can update the viewport.
	window.SetFramebufferSizeCallback(FramebufferSizeCallback)

	// Playing around with some triangles
	// Vertices for 1 triangle starting from bottom left going counter clock wise
	vertices := [9]float32 {
		-0.5, -0.5, 0.0,
		 0.5, -0.5, 0.0,
		 0.0,  0.5, 0.0,
	}
	// Tell opengl to create 1 buffer and pass us the ID
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	// We then bind that buffer as a type ARRAY_BUFFER
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	// From now on any call on the ARRAY_BUFFER target will configure this bound buffer.
	// Calling glBufferData will then copy the defined vertex data into the buffers memory
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) << 2, unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	// Dynamically compile the shaders
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	transformedSource := gl.Str(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, &transformedSource, nil)
	gl.CompileShader(vertexShader)
	// Checking compile status
	var success int32
	infoLog := [512]uint8{}
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(vertexShader, 512, nil, &infoLog[0])
		log.Fatalf("vertex shader compilation failed %s", gl.GoStr(&infoLog[0]))
	}

	for !window.ShouldClose() {
		// input
		processInput(window)

		// rendering commands
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// check and call events and swap the buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
}

func FramebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0,0, int32(width), int32(height))
}