package etf

type (
	Binstr string
)

func NewBinstr(v []byte) Binstr {
	return Binstr(v)
}
