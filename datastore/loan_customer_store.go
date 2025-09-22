package datastore

import (
	"database/sql"
)

const sqlUpsertCustomer = `
    INSERT INTO loan_customers (
        customer_id,
        id_card_number,
        full_name,
        birth_date,
        phone_number,
        email,
        monthly_income,
        address_street,
        address_city
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
    ) ON CONFLICT (id_card_number) DO UPDATE SET
        full_name = EXCLUDED.full_name,
        birth_date = EXCLUDED.birth_date,
        phone_number = EXCLUDED.phone_number,
        email = EXCLUDED.email,
        monthly_income = EXCLUDED.monthly_income,
        address_street = EXCLUDED.address_street,
        address_city = EXCLUDED.address_city
    RETURNING customer_id;
`

const sqlGetAllLoanCustomers = `
SELECT
    customer_id,
	id_card_number,
	full_name,
	birth_date,
	phone_number,
	email,
	monthly_income,
	address_street,
	address_city
FROM loan_customers;`

const sqlGetCustomerByCustomerId = `
SELECT
    customer_id,
	id_card_number,
	full_name,
	birth_date,
	phone_number,
	email,
	monthly_income,
	address_street,
	address_city
FROM loan_customers
WHERE customer_id = $1;`

type LoanCustomerRow struct {
	CustomerID    string
	IDCardNumber  string
	FullName      string
	BirthDate     string
	PhoneNumber   string
	Email         sql.NullString
	MonthlyIncome float64
	AddressStreet string
	AddressCity   string
}

type LoanCustomerStore struct {
	db *sql.DB
}

func NewLoanCustomerStore(db *sql.DB) *LoanCustomerStore {
	return &LoanCustomerStore{
		db: db,
	}
}

func (s *LoanCustomerStore) UpsertCustomer(customer *LoanCustomerRow) (string, error) {
	var customerID string
	err := s.db.QueryRow(sqlUpsertCustomer,
		customer.CustomerID,
		customer.IDCardNumber,
		customer.FullName,
		customer.BirthDate,
		customer.PhoneNumber,
		customer.Email,
		customer.MonthlyIncome,
		customer.AddressStreet,
		customer.AddressCity).Scan(&customerID)

	if err != nil {
		return "", err
	}

	return customerID, nil
}

func (s *LoanCustomerStore) GetAllLoanCustomers() ([]*LoanCustomerRow, error) {
	rows, err := s.db.Query(sqlGetAllLoanCustomers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*LoanCustomerRow
	for rows.Next() {
		customer := &LoanCustomerRow{}
		err := rows.Scan(
			&customer.CustomerID,
			&customer.IDCardNumber,
			&customer.FullName,
			&customer.BirthDate,
			&customer.PhoneNumber,
			&customer.Email,
			&customer.MonthlyIncome,
			&customer.AddressStreet,
			&customer.AddressCity,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

func (s *LoanCustomerStore) GetCustomerByCustomerId(id string) (*LoanCustomerRow, error) {
	customer := &LoanCustomerRow{}

	err := s.db.QueryRow(sqlGetCustomerByCustomerId, id).Scan(
		&customer.CustomerID,
		&customer.IDCardNumber,
		&customer.FullName,
		&customer.BirthDate,
		&customer.PhoneNumber,
		&customer.Email,
		&customer.MonthlyIncome,
		&customer.AddressStreet,
		&customer.AddressCity,
	)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
