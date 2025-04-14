package etf

func (t Tuple) small() bool {
	return len(t.v) < 256
}
