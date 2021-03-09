package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"unsafe"
	"voxel/pkg/argonaut/graphics"
)

// Called first
func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

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
	shader, err := graphics.NewShader("cmd/engine/shaders/simple_vert.glsl", "cmd/engine/shaders/simple_frag.glsl")
	if err != nil {
		log.Fatal("shader creation failed", err)
	}


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
	vertices := []float32 {
		// positions     // colors      // texture
	 	 0.5,  0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		 0.5, -0.5, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, // bottom left
		-0.5,  0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, // top left
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
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) << 2, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	// Bind EBO
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices) << 2, unsafe.Pointer(&indices), gl.STATIC_DRAW)

	// Tell OpenGL how to interpret our vertex data
	// This uses our VBO because its still bound to ARRAY_BUFFER from before
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8 * 4, unsafe.Pointer(nil))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8 * 4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	// Texture parameters
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8 * 4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// Generate and apply texture
	// ------------------------------------
	texture, err := graphics.NewTexture2D(
		"cmd/engine/textures/container.jpg",
		gl.REPEAT,
		gl.REPEAT,
		gl.LINEAR,
		gl.LINEAR)

	if err != nil {
		log.Fatalf("failed to load texture: %v", err)
	}

	// Deffered Cleanup
	// -----------------------
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)
	defer shader.Delete()

	// Loop
	// --------------
	for !window.ShouldClose() {
		// input
		processInput(window)

		// rendering commands
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw triangle
		// We now activate this shader program. Every shader and rendering call after this
		// will now use this program.
		shader.Use()

		texture.Use()
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