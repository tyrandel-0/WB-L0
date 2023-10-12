package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"service/server/http-server/model"
	"strings"

	_ "github.com/lib/pq"

	"service/server/storage"
)

type Storage struct {
	db *sql.DB
}

func New(connectionString string) (*Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Init storage successful!")

	stmt, err := db.Prepare(storage.CreateTableOrders)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	log.Println("Create table orders successful!")

	return &Storage{db: db}, nil
}

func (s Storage) GetById(id string) ([]byte, error) {
	var jsonB []byte
	err := s.db.QueryRow(storage.GetByIdFromOrders, id).Scan(&jsonB)
	if err != nil {
		if err != sql.ErrNoRows {
			return jsonB, err
		}
	}

	log.Println("Get order by id successful!")

	return jsonB, nil
}

func (s Storage) GetAll() (map[string][]byte, error) {
	all := make(map[string][]byte)
	rows, err := s.db.Query(storage.GetAllFromOrders)
	if err != nil {
		if err != sql.ErrNoRows {
			return all, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var jsonB []byte

		if err = rows.Scan(&id, &jsonB); err != nil {
			return all, err
		}
		all[id] = jsonB
	}

	if err = rows.Err(); err != nil {
		return all, err
	}

	log.Println("Get all orders successful!")

	return all, nil
}

func (s Storage) AddOrder(order model.Model) error {
	if order.OrderUID == "" {
		return errors.New("orderUID is empty")
	}

	byt, err := json.Marshal(order)
	if err != nil {
		return err
	}

	stmt, err := s.db.Prepare(storage.InsertIntoOrders)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.OrderUID, byt)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input syntax for type json") {
			return errors.New("order not a json")
		} else if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errors.New("not a unique ID")
		}

		return err
	}

	log.Println("Add order successful!")

	return nil
}
