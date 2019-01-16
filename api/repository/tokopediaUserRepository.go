package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"strings"
	"time"
	"tokopedia.se.training/Project1/usermanagesys/api/repository/dbo"
)

const DateFormat string = "2006-01-02"
const DateTimeFormat string = "2006-01-02 15:04:05"

type TokopediaUserRepository struct {
	database *sql.DB
}

func NewTokopediaUserRepository(db *sql.DB) *TokopediaUserRepository {
	return &TokopediaUserRepository{database: db}
}

func (repository *TokopediaUserRepository) Database() *sql.DB {
	return repository.database
}

func (repository *TokopediaUserRepository) GetAll(maxResult int) (user []dbo.TokopediaUser, err error) {

	tokopediaUser := []dbo.TokopediaUser{}
	query := fmt.Sprintf(`select user_id, user_name, user_email, user_pwd, status, full_name, sex, birth_date, location, msisdn, create_time, update_time
											from ws_user order by user_id asc limit %d`, maxResult)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	rows, err := repository.database.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		data := dbo.TokopediaUser{}
		err := rows.Scan(
			&data.UserId,
			&data.UserName,
			&data.UserEmail,
			&data.UserPassword,
			&data.Status,
			&data.FullName,
			&data.Sex,
			&data.BirthDate,
			&data.Location,
			&data.MSISDN,
			&data.CreateTime,
			&data.UpdateTime)

		if err != nil {
			return nil, err
		}

		var age time.Duration
		if data.BirthDate != nil {
			age = time.Since(*data.BirthDate)
			age = age / time.Hour / 24 / 365

			a := data.BirthDate.Format(DateFormat)
			data.StrBirthDate = &a
		}
		data.Age = &age

		if data.CreateTime != nil {
			a := data.CreateTime.Format(DateTimeFormat)
			data.StrCreateTime = &a
		}
		if data.UpdateTime != nil {
			a := data.UpdateTime.Format(DateTimeFormat)
			data.StrUpdateTime = &a
		}

		tokopediaUser = append(tokopediaUser, data)
	}

	return tokopediaUser, nil
}

func (repository *TokopediaUserRepository) GetAllPaging(page int, maxResult int, filter *map[string]interface{}) (user []dbo.TokopediaUser, err error) {

	page = (maxResult * page) - maxResult
	tokopediaUser := []dbo.TokopediaUser{}
	query := fmt.Sprintf(`select user_id, user_name, user_email, user_pwd, status, full_name, sex, birth_date, location, msisdn, create_time, update_time
											from ws_user %s order by user_id asc limit %d offset %d`, generateFilterWhereStatement(filter), maxResult, page)

	fmt.Println(query)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	rows, err := repository.database.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		data := dbo.TokopediaUser{}
		err := rows.Scan(
			&data.UserId,
			&data.UserName,
			&data.UserEmail,
			&data.UserPassword,
			&data.Status,
			&data.FullName,
			&data.Sex,
			&data.BirthDate,
			&data.Location,
			&data.MSISDN,
			&data.CreateTime,
			&data.UpdateTime)

		if err != nil {
			return nil, err
		}

		var age time.Duration
		if data.BirthDate != nil {
			age = time.Since(*data.BirthDate)
			age = age / time.Hour / 24 / 365

			a := data.BirthDate.Format(DateFormat)
			data.StrBirthDate = &a
		}
		data.Age = &age

		if data.CreateTime != nil {
			a := data.CreateTime.Format(DateTimeFormat)
			data.StrCreateTime = &a
		}
		if data.UpdateTime != nil {
			a := data.UpdateTime.Format(DateTimeFormat)
			data.StrUpdateTime = &a
		}

		tokopediaUser = append(tokopediaUser, data)
	}

	return tokopediaUser, nil
}

func generateFilterWhereStatement(strColumn *map[string]interface{}) string {
	if strColumn == nil {
		return ""
	}

	if len(*strColumn) > 0 {
		index := 0
		prefix := ""
		for key, val := range *strColumn {
			key = strings.ToLower(key)

			if val != nil && !reflect.ValueOf(val).IsNil() {
				if index == 0 {
					prefix = " where "
				}

				var str string
				strFormat := "%s%s=%v"
				switch val.(type) {
				case string:
					str = fmt.Sprintf("'%s%v%s'", "%", val.(string), "%")
					strFormat = "%slower(%s) like lower(%v)"
				case *string:
					str = fmt.Sprintf("'%s%v%s'", "%", *val.(*string), "%")
					strFormat = "%slower(%s) like lower(%v)"
				case int:
					str = fmt.Sprintf("%v", val.(int))
				case *int:
					str = fmt.Sprintf("%v", *val.(*int))
				case time.Time:
					str = fmt.Sprintf("%v", val.(time.Time))
				case *time.Time:
					str = fmt.Sprintf("%v", *val.(*time.Time))
				default:
					str = fmt.Sprintf("%v", *val.(*int))
				}

				fmt.Printf("key: %s | value: %v", key, str)

				prefix += fmt.Sprintf(strFormat, func() string {
					if index != 0 {
						return " and "
					}
					return ""
				}(), key, str)
				index++
			}
		}
		return prefix
	}

	return ""
}
