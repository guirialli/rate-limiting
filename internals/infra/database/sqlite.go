package database

import (
	"database/sql"
	"fmt"
	"github.com/guirialli/rater_limit/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	database string
}

func NewSqlite(database string) *Sqlite {
	return &Sqlite{
		database: database,
	}
}

func (d *Sqlite) GetConnectionString() string {
	return d.database
}

func (d *Sqlite) GetConnection() (*sql.DB, error) {
	dns := d.GetConnectionString()
	db, err := sql.Open("sqlite3", dns)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Sqlite) TryConnection() error {
	_, err := d.GetConnection()
	return err
}
func (d *Sqlite) InitDatabase(file string) error {
	db, err := d.GetConnection()
	if err != nil {
		return err
	}

	if err := utils.NewDatabaseUtils().ExecScript(db, file); err != nil {
		return err
	}

	fmt.Println("Database Initialized with success!")
	return db.Close()
}
func (d *Sqlite) InitDatabaseGetConnection(file string) (*sql.DB, error) {
	db, err := d.GetConnection()
	if err != nil {
		return nil, err
	}

	if err := utils.NewDatabaseUtils().ExecScript(db, file); err != nil {
		return nil, err
	}

	fmt.Println("Database Initialized with success!")
	return db, nil
}
func (d *Sqlite) Migrate() error {
	return nil
}
