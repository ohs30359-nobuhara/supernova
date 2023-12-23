package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type MysqlClient struct {
	client *sqlx.DB
}

type MysqlConnectionOption struct {
	Db     string
	Host   string
	User   string
	Passwd string
}

func NewMySqlClient(option MysqlConnectionOption) (MysqlClient, error) {
	var (
		con *sqlx.DB
		e   error
	)

	con, e = sqlx.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", option.User, option.Passwd, option.Host, option.Db))
	if e != nil {
		return MysqlClient{}, e
	}

	con.SetConnMaxLifetime(time.Minute * 3)
	con.SetMaxOpenConns(10)

	return MysqlClient{
		client: con,
	}, nil
}

func (own MysqlClient) Transaction(fn func(tx *sql.Tx) error) error {
	tx, e := own.client.Begin()
	if e != nil {
		return e
	}

	// エラーが起こった場合はRollbackさせる
	defer func() {
		var unknown error
		p := recover()
		switch {
		// panicが発生した場合に内容を取得して処理を継続させる
		case p != nil:
			unknown = tx.Rollback()
			fmt.Println(p)
		case e != nil:
			unknown = tx.Rollback()
			fmt.Println(e.Error())
		default:
			unknown = tx.Commit()
		}

		if unknown != nil {
			fmt.Println("unknown error: " + unknown.Error())
		}
	}()

	e = fn(tx)
	return e
}

func (own MysqlClient) Exec(tx *sql.Tx, query string, args ...any) (int64, error) {
	var (
		effected sql.Result
		e        error
	)

	if args[0] == nil {
		effected, e = tx.Exec(query)
	} else {
		effected, e = tx.Exec(query, args...)
	}

	if e != nil {
		return 0, e
	}

	cnt, e := effected.RowsAffected()
	if e != nil {
		return 0, e
	}

	return cnt, nil
}

func (own MysqlClient) Select(query string, dataset any, args ...any) error {
	if args == nil || args[0] == nil {
		return own.client.Select(dataset, query)
	} else {
		return own.client.Select(dataset, query, args...)
	}
}
