package extensions

func IsEmpty(str string) bool {
	return str == ""
}

func IsNotEmpty(str string) bool {
	return str != ""
}

func SafeUnwrap(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

func ToNullIfEmpty(str string) any {
	if IsEmpty(str) {
		return nil
	}

	return str
}
