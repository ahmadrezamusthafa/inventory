package core

import (
	"fmt"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func IsValidSKU(sku string) bool {
	matched, _ := regexp.MatchString(`[A-Z]{3}[\-][A-Z][0-9]{8}[\-][A-Z0-9]{2}[\-][A-Z0-9]{3}`, sku)
	return matched
}

func GenerateSKU(product inputs.Product) string {
	if product.Size == "" || product.Color == "" {
		return ""
	}

	prefix := "SSI-D"
	rand.Seed(time.Now().UnixNano())
	min, max := 11111111, 99999999
	middle := fmt.Sprintf("%d", rand.Intn(max-min+1)+min)

	for {
		product.Size += product.Size
		if len(product.Size) >= 2 {
			break
		}
	}

	for {
		product.Color += product.Color
		if len(product.Color) >= 3 {
			break
		}
	}

	product.Size = util.TruncateString(product.Size, 2)
	product.Color = util.TruncateString(product.Color, 3)
	return strings.ToUpper(fmt.Sprintf("%s%s-%s-%s", prefix, middle, product.Size, product.Color))
}
