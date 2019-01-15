package service

import (
	"fmt"
	"github.com/json-iterator/go"
	"net/http"
	"reflect"
	"strconv"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/repository"
	"tokopedia.se.training/Project1/usermanagesys/api/response"
	"tokopedia.se.training/Project1/usermanagesys/api/service/dto"
)

type UserService struct {
	configuration  *configuration.Configuration
	userRepository *repository.TokopediaUserRepository
}

func NewUserService(config *configuration.Configuration,
	userRepo *repository.TokopediaUserRepository) *UserService {
	return &UserService{
		configuration:  config,
		userRepository: userRepo,
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	//POST
	/*data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resWriter.WriteError("Cannot read the data")
		return
	}
	err = jsoniter.Unmarshal([]byte(data), &inputData)
	if err != nil {
		resWriter.WriteError("Unmarshal error")
		return
	}*/

	resWriter := response.New(w)

	defer func() {
		if err := recover(); err != nil {
			resWriter.WriteError(fmt.Sprintf("%v", err))
			return
		}
	}()

	queryValues := r.URL.Query()
	maxResult, err := strconv.Atoi(queryValues.Get("max_result"))
	if err != nil {
		resWriter.WriteError(err.Error())
		return
	}

	finalResult, err := s.userRepository.GetAll(maxResult)
	if err != nil {
		resWriter.WriteError(err.Error())
		return
	}

	resWriter.WriteSuccess(finalResult)
	return
}

func (s *UserService) GetUsersPaging(w http.ResponseWriter, r *http.Request) {

	resWriter := response.New(w)

	defer func() {
		if err := recover(); err != nil {
			resWriter.WriteError(fmt.Sprintf("%v", err))
			return
		}
	}()

	queryValues := r.URL.Query()
	page, err := strconv.Atoi(queryValues.Get("page"))
	if err != nil {
		resWriter.WriteError(err.Error())
		return
	}
	maxResult, err := strconv.Atoi(queryValues.Get("max_result"))
	if err != nil {
		resWriter.WriteError(err.Error())
		return
	}

	var mapFilter = make(map[string]interface{})
	strfilter := queryValues.Get("filter")
	fmt.Println("val:::", strfilter)

	if strfilter != "" {

		var filter *dto.UserFilterDto
		err = jsoniter.Unmarshal([]byte(strfilter), &filter)
		if err != nil {
			resWriter.WriteError(err.Error())
			return
		}

		const tagName = "json"
		val := reflect.ValueOf(filter).Elem()
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)
			tag := typeField.Tag

			fmt.Printf("nama field: %s,\t nil: %v,\t tipe: %v,\t tag: %s\n", typeField.Name, valueField.Elem(), typeField.Type, tag.Get(tagName))
			mapFilter[tag.Get(tagName)] = valueField.Interface()
		}
	}

	finalResult, err := s.userRepository.GetAllPaging(page, maxResult, &mapFilter)
	if err != nil {
		resWriter.WriteError("Data retrieve error")
		return
	}

	resWriter.WriteSuccess(finalResult)
	return
}
