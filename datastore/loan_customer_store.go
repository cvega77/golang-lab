package datastore

import (
	"database/sql"
	"fmt"
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
select
    customer.customer_id,
    customer.id_card_number,
    customer.full_name,
    customer.birth_date,
    customer.phone_number,
    customer.email,
    customer.monthly_income,
    customer.address_street, 
    customer.address_city,
    submission.submission_id,
    submission.vehicle_type,
    submission.vehicle_brand,
    submission.vehicle_model,
    submission.vehicle_license_number,
    submission.vehicle_odometer,
    submission.manufacturing_year,
    submission.proposed_loan_amount,
    submission.proposed_loan_tenure_month,
    submission.is_commercial_vehicle,
    submission.created_at,
    submission.updated_at
from loan_customers customer
inner join loan_submissions submission
on customer.customer_id = submission.customer_id
where customer.customer_id = $1;`

const sqlUpdateCustomerByCustomerId = `
Update loan_customers 
set full_name = COALESCE($1, full_name),
birth_date = COALESCE($2, birth_date),
phone_number = COALESCE($3, phone_number),
email = COALESCE($4, email),
monthly_income = COALESCE($5, monthly_income),
address_street = COALESCE($6, address_street),
address_city = COALESCE($7, address_city)
where customer_id = $8;`

const sqlDeleteCustomerByCostumerId = `
delete
from loan_customers
where customer_id = $1;`

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

type LoanCustomerWithAllSubmissionsRow struct {
	LoanCustomerRow *LoanCustomerRow
	LoanSubmissions []*LoanSubmissionRow
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

func (s *LoanCustomerStore) GetCustomerByCustomerId(id string) (*LoanCustomerWithAllSubmissionsRow, error) {
	rows, err := s.db.Query(sqlGetCustomerByCustomerId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customer *LoanCustomerRow
	var submissions []*LoanSubmissionRow

	for rows.Next() {
		submission := &LoanSubmissionRow{}
		if customer == nil {
			customer = &LoanCustomerRow{}
			err = rows.Scan(
				&customer.CustomerID,
				&customer.IDCardNumber,
				&customer.FullName,
				&customer.BirthDate,
				&customer.PhoneNumber,
				&customer.Email,
				&customer.MonthlyIncome,
				&customer.AddressStreet,
				&customer.AddressCity,
				&submission.SubmissionID,
				&submission.VehicleType,
				&submission.VehicleBrand,
				&submission.VehicleModel,
				&submission.VehicleLicenseNumber,
				&submission.VehicleOdometer,
				&submission.ManufacturingYear,
				&submission.ProposedLoanAmount,
				&submission.ProposedLoanTenure,
				&submission.IsCommercialVehicle,
				&submission.CreatedAt,
				&submission.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err = rows.Scan(
				new(string),
				new(string),
				new(string),
				new(string),
				new(string),
				new(sql.NullString),
				new(float64),
				new(string),
				new(string),
				&submission.SubmissionID,
				&submission.VehicleType,
				&submission.VehicleBrand,
				&submission.VehicleModel,
				&submission.VehicleLicenseNumber,
				&submission.VehicleOdometer,
				&submission.ManufacturingYear,
				&submission.ProposedLoanAmount,
				&submission.ProposedLoanTenure,
				&submission.IsCommercialVehicle,
				&submission.CreatedAt,
				&submission.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
		submissions = append(submissions, submission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, fmt.Errorf("customer with id %s not found", id)
	}

	return &LoanCustomerWithAllSubmissionsRow{
		LoanCustomerRow: customer,
		LoanSubmissions: submissions,
	}, nil
}

func (s *LoanCustomerStore) UpdateCustomerByCustomerId(customer *LoanCustomerRow) error {
	result, err := s.db.Exec(sqlUpdateCustomerByCustomerId, customer.FullName, customer.BirthDate, customer.PhoneNumber,
		customer.Email, customer.MonthlyIncome, customer.AddressStreet, customer.AddressCity,
		customer.CustomerID)

	if err != nil {
		fmt.Println("UpdateCustomerByCustomerId err:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (s *LoanCustomerStore) DeleteCustomerByCustomerId(customerId string) error {
	result, err := s.db.Exec(sqlDeleteCustomerByCostumerId, customerId)

	if err != nil {
		fmt.Println("DeleteCustomerByCustomerId err:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
