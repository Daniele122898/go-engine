package mu

import "github.com/go-gl/mathgl/mgl32"

func MultiRotate3D(angle float32, x,y,z float32) mgl32.Mat4 {
	rot := mgl32.Ident4()
	if x > 0 {
		rot = rot.Mul4(mgl32.HomogRotate3DX(angle*x))
	}
	if y > 0 {
		rot = rot.Mul4(mgl32.HomogRotate3DY(angle*y))
	}
	if z > 0 {
		rot = rot.Mul4(mgl32.HomogRotate3DZ(angle*z))
	}
	return rot
}
