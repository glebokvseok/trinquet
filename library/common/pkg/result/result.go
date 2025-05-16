package result

type Result[TData any] struct {
	Data TData
	Err  error
}
