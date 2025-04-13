package etr

type (
	Binstr string
)

func NewBinstr(v []byte) Binstr {
	return Binstr(v)
}
