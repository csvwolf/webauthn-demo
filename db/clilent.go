package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
)

type Client struct {
	db *sql.DB
}

type ExecContext interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type QueryContext interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

var (
	once sync.Once
	cli  *Client
)

func NewClient() *Client {
	once.Do(func() {
		cli = &Client{}
		dsn := os.Getenv("mysql_dsn")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic("Error connecting to database:" + err.Error())
		}
		cli.db = db
	})
	return cli
}

func GetClient() *Client {
	return cli
}

func (c *Client) GetDB() *sql.DB {
	return c.db
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) DoTransaction(txFunc func(*sql.Tx) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			perr := tx.Rollback()
			fmt.Println("[db][DoTransaction] panic rollback with error:", perr, "panic:", p)
		} else if err != nil {
			rerr := tx.Rollback() // err is non-nil; don't change it
			fmt.Println("[db][DoTransaction] rollback with error:", err, "rollback error:", rerr)
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
			if err != nil {
				fmt.Println("[db][DoTransaction] commit with error:", err)
			}
		}
	}()
	err = txFunc(tx)
	return err
}
