package utils

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
)

type Database struct{}

func NewDatabaseUtils() *Database {
	return &Database{}
}

func (f *Database) InitDatabase(db *sql.DB, file string) error {
	initFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("cold not open file: %s", err.Error())
	}
	defer initFile.Close()
	reader := bufio.NewReader(initFile)
	var sqlStmt string

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("unknow error on read line: %s", err.Error())
		}
		sqlStmt += line
		if line == ");\n" || line == ";\n" || line == ");" {
			if _, err := db.Exec(sqlStmt); err != nil {
				return err
			}
			sqlStmt = ""
		}
	}
	return nil
}
