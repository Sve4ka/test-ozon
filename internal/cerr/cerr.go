package cerr

import "errors"

var (
	ErrNotFound                  = errors.New("link not found")
	ErrShortCodeAlreadyExists    = errors.New("short code already exists")
	ErrOriginalLinkAlreadyExists = errors.New("original link already exists")
	ErrNotGenerated              = errors.New("not generated")
)
