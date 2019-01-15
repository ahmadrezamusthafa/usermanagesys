package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"tokopedia.se.training/Project1/usermanagesys/api/repository/dbo"
)

type TokopediaUserRepository struct {
	database *sql.DB
}

func NewTokopediaUserRepository(db *sql.DB) *TokopediaUserRepository {
	return &TokopediaUserRepository{database: db}
}

func (repository *TokopediaUserRepository) Database() *sql.DB {
	return repository.database
}

func (repository *TokopediaUserRepository) GetAll() (user []dbo.TokopediaUser, err error) {

	tokopediaUser := []dbo.TokopediaUser{}
	stmt, err := repository.database.Prepare(`select user_id, user_name, user_email, user_pwd, status, full_name, sex, birth_date, location, msisdn, create_time, update_time
											from ws_user order by full_name asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		data := dbo.TokopediaUser{}
		rows.Scan(
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

		tokopediaUser = append(tokopediaUser, data)

		err = rows.Err()
		if err != nil {
			log.Println(err.Error())
		}
	}

	return tokopediaUser, nil
}
