package handler

import (
	"fmt"
	"regexp"
	"testing"
)

func TestEcho(t *testing.T) {
	ethRegex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	//address := "0x7E76e07149671635fbEd2f7Fa74D29bc1234"
	address := "0x96F55350F58beB9B00a871b6C37AF5179c2a4760"
	if !ethRegex.MatchString(address) {

		fmt.Printf("错误的参数")
	}
	fmt.Printf("成功\n")
}
