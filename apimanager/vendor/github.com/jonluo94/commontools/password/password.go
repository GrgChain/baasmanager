package password

import (
	"github.com/alexandrevicenzi/unchained"
	"github.com/jonluo94/commontools/log"
)

var logger = log.GetLogger("password", log.ERROR)

func Encode(password string, saltSize int, hasher string) string {
	hash, err := unchained.MakePassword(password, unchained.GetRandomString(saltSize), hasher)
	if err != nil {
		logger.Errorf("Error encoding password: %s\n", err)
	}
	return string(hash)
}

func Validate(password, cryto string) bool {
	valid, err := unchained.CheckPassword(password, cryto)
	if err != nil {
		logger.Errorf("Error decoding password: %s\n", err)
	}
	return valid
}
