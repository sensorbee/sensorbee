package core

type Box interface {
	Process(t *Tuple, s Sink) error
}

type BoxFunc func(t *Tuple, s Sink) error

func (b BoxFunc) Process(t *Tuple, s Sink) error {
	return b(t, s)
}
