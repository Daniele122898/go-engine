package mu

import "math"

func RadianToDegree (angle float64) float64 {
	return angle*180/math.Pi
}

func DegreeToRadian (angle float64) float64 {
	return angle / (180/math.Pi)
}
