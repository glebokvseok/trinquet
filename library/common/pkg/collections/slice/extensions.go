package slice

func HasElements[TValue any](slice []TValue) bool {
	return len(slice) != 0
}

func Map[TIn any, TOut any](slice []TIn, mapFunc func(TIn) TOut) []TOut {
	result := make([]TOut, len(slice))
	for i, element := range slice {
		result[i] = mapFunc(element)
	}

	return result
}

func MapErr[TIn any, TOut any](slice []TIn, mapFunc func(TIn) (TOut, error)) ([]TOut, error) {
	var err error
	result := make([]TOut, len(slice))
	for i, element := range slice {
		result[i], err = mapFunc(element)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
