package models

import (
	"cointrade/utils"
	"fmt"
	"testing"
)

func TestUserModel_EncodePassword(t *testing.T) {

	password := "123456789"
	passwordPrefix := fmt.Sprintf("%s%s", PASSMIX, password)
	passwordSuffix := utils.Md5(fmt.Sprintf("%s%s", PASSMIX, password))

	fmt.Println(passwordPrefix, passwordSuffix)
}
