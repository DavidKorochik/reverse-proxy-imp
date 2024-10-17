package errors

import "github.com/rotisserie/eris"

func New(message string) error {
	return eris.New(message)
}

func Wrap(err error, message string) error {
	return eris.Wrap(err, message)
}

func WrapF(err error, message string, args ...any) error {
	return eris.Wrapf(err, message, args)
}
