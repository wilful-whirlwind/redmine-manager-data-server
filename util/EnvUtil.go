package util

import (
	"os"
)

type EnvUtil struct {
}

func (util EnvUtil) GetPasswordPepper() string {
	return os.Getenv("PASSWORD_PEPPER")
}
