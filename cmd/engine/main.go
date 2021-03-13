package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
	"unsafe"
	"voxel/pkg/argonaut/graphics"
	"voxel/pkg/mu"
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

	gl.Viewport(0, 0, 1280, 720)
	// Register callback in case user changes the window size so gl can update the viewport.
	window.SetFramebufferSizeCallback(FramebufferSizeCallback)

	gl.Enable(gl.DEPTH_TEST)

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
	//vertices := []float32{
	//	// positions     // colors      // texture
	//	0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
	//	0.5, -0.5, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, // bottom right
	//	-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, // bottom left
	//	-0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, // top left
	//}

	// Vertices for a cube
	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	indices := [6]int32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
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
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)<<2, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	// Bind EBO
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)<<2, unsafe.Pointer(&indices), gl.STATIC_DRAW)

	// Tell OpenGL how to interpret our vertex data
	// This uses our VBO because its still bound to ARRAY_BUFFER from before
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, unsafe.Pointer(nil))
	gl.EnableVertexAttribArray(0)

	// Texture parameters
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Generate and apply texture
	// ------------------------------------
	texture, err := graphics.NewTexture2D(
		"cmd/engine/textures/container.jpg",
		gl.REPEAT,
		gl.REPEAT,
		gl.LINEAR,
		gl.LINEAR,
		0,
		0)

	texture2, err := graphics.NewTexture2D(
		"cmd/engine/textures/awesomeface.png",
		gl.REPEAT,
		gl.REPEAT,
		gl.LINEAR,
		gl.LINEAR,
		1,
		180)

	if err != nil {
		log.Fatalf("failed to load texture: %v", err)
	}

	// Deffered Cleanup
	// -----------------------
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteBuffers(1, &VBO)
	defer shader.Delete()

	// Tell OpenGL which texture unit each shader sampler belongs to.
	shader.Use()
	shader.SetInt("texture1", 0)
	shader.SetInt("texture2", 1)

	// Math and transformations
	// --------------
	// Model matrix that is kinda on the ground on the x axis
	//model := mgl32.HomogRotate3DX(float32(glfw.GetTime()) * mgl32.DegToRad(-55))
	// pretend like our camera goes back which is the same as moving the entire
	// world back in the -z axis
	view := mgl32.Translate3D(0, 0, -3)
	// Create projection matrix
	proj := mgl32.Perspective(mgl32.DegToRad(45), 1280.0/720.0, 0.1, 100)
	vp := proj.Mul4(view)

	cubePos := [][3]float32{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	// Loop
	// --------------
	for !window.ShouldClose() {

		// input
		processInput(window)

		// rendering commands
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw triangle
		// We now activate this shader program. Every shader and rendering call after this
		// will now use this program.
		shader.Use()
		//trans := mgl32.Ident4()
		//trans = trans.Mul4(mgl32.Translate3D(0.5, -0.5, 0.0))
		//trans = trans.Mul4(mgl32.HomogRotate3DZ(float32(glfw.GetTime())))
		//shader.SetMatrix4f("transform", trans)

		//model := mgl32.HomogRotate3DX(float32(glfw.GetTime()) * mgl32.DegToRad(50))
		//model = model.Mul4(mgl32.HomogRotate3DY(float32(glfw.GetTime()) * mgl32.DegToRad(50)))
		//mvp := vp.Mul4(model)

		// Send projection matrices

		texture.UseActive(gl.TEXTURE0)
		texture2.UseActive(gl.TEXTURE1)
		gl.BindVertexArray(VAO)

		for i, pos := range cubePos {
			model := mgl32.Translate3D(pos[0], pos[1], pos[2])
			angle := float32(20.0 * i)
			//model = model.Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(angle)))
			model = model.Mul4(mu.MultiRotate3D(mgl32.DegToRad(angle), 1, 0.3, 0.5))
			mvp := vp.Mul4(model)
			shader.SetMatrix4f("mvp", mvp)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}


		//gl.DrawArrays(gl.TRIANGLES, 0, 36)
		//gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))

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
	gl.Viewport(0, 0, int32(width), int32(height))
}
