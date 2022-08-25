package tool

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

func GetRealDecimalValue(value string, decimal int) string {
	if strings.Contains(value, ".") {
		arr := strings.Split(value, ".")
		if len(arr) != 2 {
			return ""
		}
		numOf0 := decimal - len(arr[1])
		return arr[0] + arr[1] + strings.Repeat("0", numOf0)
	} else {
		return value + strings.Repeat("0", decimal)
	}
}
func BuildERC20TransferData(receiver, value string, decimal int) string {
	realValue := GetRealDecimalValue(value, decimal)
	valueBig, _ := new(big.Int).SetString(realValue, 10)
	param1 := common.HexToHash(receiver).String()[2:]
	param2 := common.BytesToHash(valueBig.Bytes()).String()[2:]
	return "0xa9059cbb" + param1 + param2
}
func ErrorPrintln(err error, description string) {
	if err != nil {
		if description != "" {
			fmt.Println(description, ":", err)
		} else {
			fmt.Println(err)
		}
	}
}
