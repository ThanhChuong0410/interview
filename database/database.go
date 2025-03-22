package database

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"context"

	"github.com/chuongthanh0410/interview/config"
	"github.com/jackc/pgx/v5"
)

type globalClient struct {
	*pgx.Conn
}

var DBClient *globalClient

func init() {
	log.Println("NewDBClient")
	err := NewDBClient()
	if err != nil {
		panic(err)
	}
}

func NewDBClient() error {
	if DBClient != nil {
		return nil
	}

	info := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.USER_DATABASE,
		config.PASSWORD_DATABASE,
		config.HOST_DATABASE,
		config.PORT_DATABASE,
		config.DBNAME_DATABASE,
	)
	conn, err := pgx.Connect(context.Background(), info)
	if err != nil {
		return errors.New("cannot to connect to database: " + err.Error())
	}

	DBClient = &globalClient{
		conn,
	}
	return nil
}

func (c *globalClient) SelectOne(table, field, val string) (pgx.Rows, error) {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s';", table, field, val)
	rows, err := c.Query(context.Background(), cmd)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *globalClient) Search(table, cmd string, args []any) (pgx.Rows, error) {
	fmt.Println(cmd, args)
	rows, err := c.Query(context.Background(), cmd, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *globalClient) Count(table string) (int, error) {
	cmd := "SELECT COUNT(*) FROM products"
	var total int
	err := c.QueryRow(context.Background(), cmd).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (c *globalClient) GetInfoSchema(table string) (pgx.Rows, error) {
	cmd := fmt.Sprintf(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = '%s'`, table)
	rows, err := c.Query(context.Background(), cmd)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *globalClient) UpdateOne(table, id, setCmd string, args []interface{}) error {
	cmd := fmt.Sprintf("UPDATE %s %s WHERE id = '%s'", table, setCmd, id)
	fmt.Println(cmd)
	_, err := c.Exec(context.Background(), cmd, args...)
	if err != nil {
		return fmt.Errorf("exec command %s, error %s", cmd, err.Error())
	}
	return nil
}

func (c *globalClient) DeleteOne(table, id string) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", table, id)
	fmt.Println(cmd)
	_, err := c.Exec(context.Background(), cmd)
	if err != nil {
		return fmt.Errorf("exec command %s, error %s", cmd, err.Error())
	}
	return nil
}

func BuildSetCommand(updateData map[string]interface{}) (string, []interface{}) {
	var setClauses []string
	var args []interface{}
	counter := 1

	for column, value := range updateData {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, counter))
		args = append(args, value)
		counter++
	}

	query := fmt.Sprintf("SET %s", strings.Join(setClauses, ", "))

	return query, args
}
