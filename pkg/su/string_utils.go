package su

func CStr(str string) *uint8 {
	cstr := make([]uint8, len(str) + 1)
	copy(cstr, str)
	cstr[len(str)] = 0x00
	return &cstr[0]
}