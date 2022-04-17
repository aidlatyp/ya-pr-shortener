package domain

// Url model represents an url as structure
type Url struct {
	Orig  string
	Short string
}

func NewUrl(original, short string) *Url {
	return &Url{
		Orig:  original,
		Short: short,
	}
}
