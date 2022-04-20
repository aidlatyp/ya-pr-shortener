package domain

// URL model represents an url as structure
// to extend it with new properties
type URL struct {
	Orig  string
	Short string
}

func NewURL(original, short string) *URL {
	return &URL{
		Orig:  original,
		Short: short,
	}
}
