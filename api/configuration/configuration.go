package configuration

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration/dto"
)

type Configuration struct {
	ConnectionString dto.ConnectionStringDto
	App              dto.AppDto
	NSQ              dto.NSQDto
}

func NewConfiguration()(*Configuration, error) {
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
