package databasetest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type RaterLimit[T any] struct {
	db        *sql.DB
	tableName string
}

func NewRaterLimit[T any](db *sql.DB, tableName string) *RaterLimit[T] {
	return &RaterLimit[T]{
		db:        db,
		tableName: tableName,
	}
}

func (r *RaterLimit[T]) Get(ctx context.Context, key string) (*T, bool) {
	var entity T

	query := fmt.Sprintf("SELECT * FROM %s WHERE Id = ?", r.tableName)
	row := r.db.QueryRowContext(ctx, query, key)

	val := reflect.ValueOf(&entity).Elem()
	fields := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}

	if err := row.Scan(fields...); err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false
		}
		return nil, false
	}
	return &entity, true
}

func (r *RaterLimit[T]) Set(ctx context.Context, key string, entity *T) error {
	_, exist := r.Get(ctx, key)
	if exist {
		return r.update(ctx, key, entity)
	}
	return r.create(ctx, key, entity)
}

func (r *RaterLimit[T]) create(ctx context.Context, key string, entity *T) error {
	val := reflect.ValueOf(entity).Elem()
	typeOfEntity := val.Type()

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typeOfEntity.Field(i)
		var value any
		if field.Name == "Id" {
			value = key
		} else {
			value = val.Field(i).Interface()
		}

		columns = append(columns, field.Name)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "),
	)

	_, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		log.Printf("Erro ao criar registro: %v", err)
		return err
	}

	return nil
}

func (r *RaterLimit[T]) update(ctx context.Context, key string, entity *T) error {
	val := reflect.ValueOf(entity).Elem()
	typeEntity := val.Type()

	var updates []string
	var values []interface{}
	for i := 0; i < val.NumField(); i++ {
		field := typeEntity.Field(i)
		if field.Name != "Id" {
			updates = append(updates, fmt.Sprintf("%s = ?", field.Name))
			values = append(values, val.Field(i).Interface())
		}
	}
	values = append(values, key)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE Id = ?", r.tableName, strings.Join(updates, ", "))

	_, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
