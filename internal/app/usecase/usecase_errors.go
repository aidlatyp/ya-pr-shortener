package usecase

import "fmt"

// ErrAlreadyExists represents usecase layer error with wrapped
// context error, can be logged and contain User understandable message
type ErrAlreadyExists struct {
	Err            error
	ExistShortenID string
	Orig           string
	// can hold PrettyMsg as property field
	// and be changed and filled  with usecase layer which knows
	// how to deal with user.(??)

	// Any restrictions of this approach?
	PrettyMsg string // Not used for simplicity now
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("duplication error %v", e.ExistShortenID)
}

func (e ErrAlreadyExists) Pretty() string {
	return fmt.Sprintf("Sorry, you have already saved this url %v ", e.Orig)
}

type ErrURLDeleted struct {
	Err       error
	ShortID   string
	Orig      string
	PrettyMsg string
}

func (e ErrURLDeleted) Error() string {
	return fmt.Sprintf("short key id  %v marked deleted, original is %v", e.ShortID, e.Orig)
}

func (e ErrURLDeleted) Pretty() string {
	return "Unfortunately, requested url was already been deleted"
}
