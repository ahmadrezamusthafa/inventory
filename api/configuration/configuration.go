package configuration

import (
	"encoding/json"
	"github.com/rezamusthafa/inventory/api/configuration/types"
	"os"
	"path"
	"runtime"
)

type Configuration struct {
	App              types.App
	ConnectionString types.ConnectionString
}

func NewConfiguration() (*Configuration, error) {
	configuration := Configuration{}
	_, runningFile, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(runningFile), "./", "config/appSetting.json")
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}
