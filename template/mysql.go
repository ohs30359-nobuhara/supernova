package template

import (
	"database/sql"
	"fmt"
	"strings"
	"supernova/pkg/database"
)

type MysqlTemplate struct {
	Connection database.MysqlConnectionOption `yaml:"connection"`
	Sql        []string                       `yaml:"sql"`
}

func (t MysqlTemplate) Run() Output {
	var output Output
	client, e := database.NewMySqlClient(t.Connection)
	if e != nil {
		return output.SetBody(OutputBody{
			Body:        []byte("failed to create MySQL client. " + e.Error()),
			Status:      OutputStatusDanger,
			ContentType: OutputTypeText,
		})
	}

	for _, query := range t.Sql {
		if isSelectQuery(query) {
			err := t.runSelectQuery(client, query)
			if err != nil {
				return output.SetBody(OutputBody{
					Body:        []byte("failed to execute SELECT query. " + e.Error()),
					Status:      OutputStatusDanger,
					ContentType: OutputTypeText,
				})
			}
		} else {
			err := t.runNonSelectQuery(client, query)
			if err != nil {
				return output.SetBody(OutputBody{
					Body:        []byte("failed to execute non-SELECT query. " + e.Error()),
					Status:      OutputStatusDanger,
					ContentType: OutputTypeText,
				})
			}
		}
	}

	// TODO: SQLの実行結果を送るようにする
	return output.SetBody(OutputBody{
		Body:        []byte("sql executed successfully"),
		Status:      OutputStatusOK,
		ContentType: OutputTypeText,
	})
}

// runSelectQuery SELECT SQL を実行する
func (t MysqlTemplate) runSelectQuery(client database.MysqlClient, query string) error {
	var dataset interface{}
	if err := client.Select(query, &dataset); err != nil {
		return fmt.Errorf("failed to execute SELECT query: %w", err)
	}
	fmt.Println(query)
	return nil
}

// runNonSelectQuery DELETE, UPDATE, INSERT 系のSQLを実行する
func (t MysqlTemplate) runNonSelectQuery(client database.MysqlClient, query string) error {
	return client.Transaction(func(tx *sql.Tx) error {
		_, err := client.Exec(tx, query)
		return err
	})
}

// isSelectQuery SQLがselectかを判定する
func isSelectQuery(sql string) bool {
	normalizedSQL := strings.TrimSpace(strings.ToUpper(sql))

	// 先頭の単語がSELECTで始まるかどうかを確認します
	return strings.HasPrefix(normalizedSQL, "SELECT")
}
