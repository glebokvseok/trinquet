package extensions

func IsNotZero(num int) bool {
	return num != 0
}

func ToNullIfZero(num int) any {
	if num == 0 {
		return nil
	}

	return num
}
