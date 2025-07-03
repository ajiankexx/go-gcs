package appError
import(
	"errors"
)

var (
	ErrorUser = errors.New("error about User")
	ErrorUserNotFound = errors.New("user not found")
	ErrorPasswordInvalid = errors.New("invalid password")
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorWrongPassword = errors.New("wrong password")
	ErrorEmptyPassword = errors.New("empty password")
	ErrorInvalidPassword = errors.New("invalid password")

	ErrorEmailSend = errors.New("email send error")

	ErrorLabel = errors.New("this is a label of error")
	ErrorRedisNotFoundKey = errors.New("not found key in redis")
	ErrorWrongVerifyCode = errors.New("wrong verify code")
	ErrorExpiredVerifyCode = errors.New("Vefify Code has expired")

)
