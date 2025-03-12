package utils

import "strconv"

func (u *Utils) atoi(s string) int {
	intVal, err := strconv.Atoi(s)
	if err != nil {
		intVal = 0
	}
	return intVal
}

func (u *Utils) Atoi(s string) int {
	return u.atoi(s)
}

func (u *Utils) atoi64(s string) int64 {
	intVal, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		intVal = 0
	}
	return intVal
}

func (u *Utils) atoui8(s string) uint8 {
	intVal, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		intVal = 0
	}
	return uint8(intVal)
}

func (u *Utils) Atob(s string) bool {
	return u.atob(s)
}

func (u *Utils) atob(s string) bool {
	boolVal, err := strconv.ParseBool(s)
	if err != nil {
		boolVal = false
	}
	return boolVal
}
