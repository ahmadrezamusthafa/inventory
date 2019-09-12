package core

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
	"math/rand"
	"os"
	"reflect"
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

func ExportDataToCSV(title [][]string, fileName string, data interface{}) (err error) {

	if fileName == "" {
		return errors.New("File name is required")
	}

	var (
		cache         = make(map[int]string)
		tableHeaders  []string
		tableContents [][]string
	)

	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		return errors.New("Reflection attributes is invalid")
	}

	t = t.Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		headerTag, ok := field.Tag.Lookup("header")
		if !ok {
			continue
		}

		displayHeaderName := field.Name
		if ok && headerTag != "" {
			displayHeaderName = headerTag
		}

		cache[i] = field.Name
		tableHeaders = append(tableHeaders, displayHeaderName)
	}

	s := reflect.ValueOf(data)
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)

		rowData := []string{}
		for i := 0; i < t.NumField(); i++ {
			if _, ok := cache[i]; ok {
				rowData = appendInterfaceToSliceString(rowData, v.Field(i).Interface())
			}
		}

		tableContents = append(tableContents, rowData)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//title
	for _, val := range title {
		err = writer.Write(val)
		if err != nil {
			return err
		}
	}

	//space
	for i := 0; i < 2; i++ {
		err = writer.Write([]string{""})
		if err != nil {
			return err
		}
	}

	//header
	err = writer.Write(tableHeaders)
	if err != nil {
		return err
	}

	//body
	for _, val := range tableContents {
		err = writer.Write(val)
		if err != nil {
			return err
		}
	}

	return
}

func appendInterfaceToSliceString(str []string, val interface{}) []string {
	switch val.(type) {
	case time.Time:
		tim := val.(time.Time)
		if tim.IsZero() {
			str = append(str, "")
		} else {
			str = append(str, tim.Format(util.TimestampWithSecond))
		}

	default:
		str = append(str, fmt.Sprint(val))
	}
	return str
}
