package mu

import "math"

func RadianToDegree64(angle float64) float64 {
	return angle*180/math.Pi
}

func DegreeToRadian64(angle float64) float64 {
	return angle / (180/math.Pi)
}

func RadianToDegree32(angle float32) float32 {
	return angle*180/math.Pi
}

func DegreeToRadian32(angle float32) float32 {
	return angle / (180/math.Pi)
}

