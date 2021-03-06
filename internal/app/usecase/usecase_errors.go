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
	// Could be a message from PrettyMsg in real
	return fmt.Sprintf("Sorry, you have already saved this url %v ", e.Orig)
}
