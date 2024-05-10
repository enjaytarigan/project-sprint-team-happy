package repository

import (
	"database/sql"
	"eniqlostore/internal/entity"
	"errors"
	"fmt"
)

type customerRepository struct {
	db *sql.DB
}

var (
	errCustomerNotFound = errors.New("customer not found")
)

func NewCustomerRepository(db *sql.DB) *customerRepository {
	return &customerRepository{
		db: db,
	}
}

type ICustomerRepository interface {
	Insert(entity.Customer) (entity.Customer, error)
	CheckExistByPhoneNumber(phoneNumber string) (isExist bool, err error)
	FindBy(name string, phoneNumber string) ([]entity.Customer, error)
	GetById(id string) (entity.Customer, error)
}

func (c *customerRepository) GetById(id string) (entity.Customer, error) {
	query := `
		SELECT id, name, phone_number FROM customers WHERE id = $1
	`
	var cust entity.Customer
	err := c.db.QueryRow(query, id).Scan(&cust.ID, &cust.Name, &cust.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Customer{}, errCustomerNotFound
		} else {
			return cust, err
		}
	}
	return cust, err
}

func (c *customerRepository) CheckExistByPhoneNumber(phoneNumber string) (isExist bool, err error) {
	query := `
		SELECT COUNT(phone_number) FROM customers WHERE phone_number = $1
	`

	var count int
	err = c.db.QueryRow(query, phoneNumber).Scan(&count)
	isExist = count > 0

	return isExist, err
}

func (c *customerRepository) Insert(cust entity.Customer) (entity.Customer, error) {
	query := `
		INSERT INTO customers(name, phone_number) 
		VALUES($1, $2)
		RETURNING id
		`

	newCust := cust
	err := c.db.QueryRow(query, cust.Name, cust.PhoneNumber).Scan(&newCust.ID)

	return newCust, err
}

func (c *customerRepository) FindBy(name string, phoneNumber string) ([]entity.Customer, error) {
	values := []any{}

	query := `
		SELECT id, name, phone_number FROM customers WHERE 1=1
	`

	if name != "" {
		values = append(values, fmt.Sprintf("%%%s%%", name))
		query = fmt.Sprintf("%s AND name ILIKE $%d", query, len(values))
	}

	if phoneNumber != "" {
		values = append(values, fmt.Sprintf("%s%%", phoneNumber))
		query = fmt.Sprintf("%s AND phone_number LIKE $%d", query, len(values))
	}

	query += ` ORDER BY created_at DESC`

	rows, err := c.db.Query(query, values...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []entity.Customer = []entity.Customer{}
	for rows.Next() {
		var cust entity.Customer

		if err := rows.Scan(&cust.ID, &cust.Name, &cust.PhoneNumber); err != nil {
			return customers, err
		}

		customers = append(customers, cust)
	}

	if err := rows.Err(); err != nil {
		return customers, err
	}

	return customers, nil
}
