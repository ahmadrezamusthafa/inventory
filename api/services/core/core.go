package core

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
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

func ReadSheetProduct(config *configuration.Configuration) (products []dbo.Product, err error) {

	type Credential struct {
		Type         string `json:"type"`
		ProjectID    string `json:"project_id"`
		PrivateKeyID string `json:"private_key_id"`
		PrivateKey   string `json:"private_key"`
		ClientEmail  string `json:"client_email"`
		ClientID     string `json:"client_id"`
		TokenURL     string `json:"token_uri"`
	}

	credential, err := jsoniter.Marshal(Credential{
		Type:         config.SheetCredential.Type,
		ProjectID:    config.SheetCredential.ProjectID,
		PrivateKeyID: config.SheetCredential.PrivateKeyID,
		PrivateKey:   config.SheetCredential.PrivateKey,
		ClientEmail:  config.SheetCredential.ClientEmail,
		ClientID:     config.SheetCredential.ClientID,
		TokenURL:     config.SheetCredential.TokenURL,
	})
	if err != nil {
		return products, errors.New("Marshal error")
	}

	scope := config.SheetCredential.Scope
	if scope == "" {
		return products, errors.New("Scope is required")
	}

	jwtConfig, err := google.JWTConfigFromJSON(credential, scope)
	if err != nil {
		return products, errors.New(fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}

	client := jwtConfig.Client(oauth2.NoContext)

	srv, err := sheets.New(client)
	if err != nil {
		return products, errors.New(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}

	spreadsheetId := config.SheetCredential.SheetID
	readRange := "A2:B"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return products, errors.New(fmt.Sprintf("Unable to retrieve data from sheet: %v", err))
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found")
	} else {
		for _, row := range resp.Values {
			if len(row) >= 2 {
				sku := strings.TrimSpace(fmt.Sprint(row[0]))
				name := strings.TrimSpace(fmt.Sprint(row[1]))

				if sku == "" || name == "" {
					continue
				}

				products = append(products, dbo.Product{
					SKU:  sku,
					Name: name,
				})
			}
		}
	}

	return products, nil
}
