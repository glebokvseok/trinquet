package event

type Handler interface {
	Handle(wrapper Wrapper)
}
