package appError
import(
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrPasswordInvalid = errors.New("invalid password")
)
