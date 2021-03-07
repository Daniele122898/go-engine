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
	vertexShaderSource = "#version 460 core\n\nlayout (location = 0) in vec3 aPos;\n\nout vec4 vertexColor; // specify color output for frag shader\n\nvoid main() {\n    gl_Position = vec4(aPos, 1.0);\n    vertexColor = vec4(aPos, 1);\n}\000"
	fragmentShaderSource = "#version 460\n\nout vec4 FragColor;\n\nin vec4 vertexColor; // Input from vert shader. Same name and type\n\nvoid main() {\n    FragColor = vertexColor;\n}\000"
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

	// Dynamically compile the shaders
	// -------------------------------
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
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	transformedFragsource := gl.Str(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, &transformedFragsource, nil)
	gl.CompileShader(fragmentShader)
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(fragmentShader, 512, nil, &infoLog[0])
		log.Fatalf("fragment shader compilation failed %s", gl.GoStr(&infoLog[0]))
	}

	// Linking the shaders into a shader program
	shaderProg := gl.CreateProgram()
	gl.AttachShader(shaderProg, vertexShader)
	gl.AttachShader(shaderProg, fragmentShader)
	gl.LinkProgram(shaderProg)
	gl.GetShaderiv(shaderProg, gl.LINK_STATUS, &success)
	if success == 0 {
		gl.GetProgramInfoLog(shaderProg, 512, nil, &infoLog[0])
		log.Fatalf("shader programm linking failed %s", gl.GoStr(&infoLog[0]))
	}
	// Now delete the shader objects as we already linked them and dont need them
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// Setup Verted data and buffers etc
	// ------------------------------------
	// Create a VAO
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	// Playing around with some triangles
	//vertices := [18]float32 {
	//	// first triangle
	//	 0.5,  0.5, 0.0, // top right
	//	 0.5, -0.5, 0.0, // bottom right
	//	-0.5,  0.5, 0.0, // top let
	//	// second triangle
	//	 0.5, -0.5, 0.0, // bottom right
	//	-0.5, -0.5, 0.0, // bottom let
	//	-0.5,  0.5, 0.0, // top let
	//}

	// Using Element Buffer objects to not specify verticies twice
	vertices := [12]float32 {
		// first triangle
		0.5,  0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom let
		-0.5,  0.5, 0.0, // top let
	}
	indices := [6]int32 {
		0, 1, 3,   // first triangle
		1, 2, 3,    // second triangle
	}

	// Tell opengl to create 1 buffer and pass us the ID
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	// Generate the EBO buffer
	var EBO uint32
	gl.GenBuffers(1, &EBO)
	// We then bind that buffer as a type ARRAY_BUFFER
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	// From now on any call on the ARRAY_BUFFER target will configure this bound buffer.
	// Calling glBufferData will then copy the defined vertex data into the buffers memory
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) << 2, unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	// Bind EBO
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices) << 2, unsafe.Pointer(&indices), gl.STATIC_DRAW)

	// Tell OpenGL how to interpret our vertex data
	// This uses our VBO because its still bound to ARRAY_BUFFER from before
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3 * 4, unsafe.Pointer(nil))
	gl.EnableVertexAttribArray(0)

	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)
	defer gl.DeleteProgram(shaderProg)

	for !window.ShouldClose() {
		// input
		processInput(window)

		// rendering commands
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw triangle
		// We now activate this shader program. Every shader and rendering call after this
		// will now use this program.
		gl.UseProgram(shaderProg)
		gl.BindVertexArray(VAO)
		//gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))

		// check and call events and swap the buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}

	if w.GetKey(glfw.KeyF) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
	if w.GetKey(glfw.KeyP) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func FramebufferSizeCallback(_ *glfw.Window, width int, height int) {
	gl.Viewport(0,0, int32(width), int32(height))
}