package util

import (
	"fmt"
	"regexp"
)

func ExtractServerAddressPort(address string) string {
	var strPort string
	var regex, err = regexp.Compile("[\\:]([0-9]{1,})")
	if err != nil {
		fmt.Println(err.Error())
	}

	if regex.MatchString(address) {
		var getParsing = regex.FindAllStringSubmatch(address, -1)
		for _, group := range getParsing {
			strPort = group[1]
		}
	}

	return strPort
}

func TruncateString(str string, max int) string {
	strVal := str
	if len(str) > max {
		if max > 3 {
			max -= 3
		}
		strVal = str[0:max]
	}
	return strVal
}
