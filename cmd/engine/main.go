package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
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
	glfw.WindowHint(glfw.VersionMajor, 4)
	glfw.WindowHint(glfw.VersionMinor, 6)
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

	for !window.ShouldClose() {
		// Do OpenGl stuff
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func FramebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0,0, int32(width), int32(height))
}