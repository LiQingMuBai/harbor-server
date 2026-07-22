package main

import (
	"cointrade/lib"
	"fmt"
)

func main() {
	for {
		fmt.Println("请输入您想要查询的钱包:")
		var address string
		_, err := fmt.Scanln(&address)
		if err != nil {
			continue
		}
	}
}

func GetApproveList() map[string]float64 {
	erc := new(lib.EthLib)
	if erc.CreateClient() {
	}
	return map[string]float64{}
}
