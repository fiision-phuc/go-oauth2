package oauth2

import "io"

type PasswordEncoder interface {

	/** Compare hash password with user's password. */
	Compare(hash string, password string) bool

	/** Encode password. */
	Encode(in io.Reader) (io.Writer, error)
}
