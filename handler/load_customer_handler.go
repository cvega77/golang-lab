package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
)

type LoanCustomerHandler struct {
	CustomerStore datastore.LoanCustomerStore
}

func NewLoanCustomerHandler(customerStore datastore.LoanCustomerStore) *LoanCustomerHandler {
	return &LoanCustomerHandler{
		CustomerStore: customerStore,
	}
}

func (h *LoanCustomerHandler) HandleGetAllLoanSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	loanCustomerRows, err := h.CustomerStore.GetAllLoanCustomers()
	if err != nil {
		errMsg := "Failed to get all loan customers"
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	loanCustomers := make([]LoanCustomer, 0, len(loanCustomerRows))
	for _, row := range loanCustomerRows {
		customer := LoanCustomer{
			CustomerID:    row.CustomerID,
			IDCardNumber:  row.IDCardNumber,
			FullName:      row.FullName,
			BirthDate:     row.BirthDate,
			PhoneNumber:   row.PhoneNumber,
			MonthlyIncome: row.MonthlyIncome,
			AddressStreet: row.AddressStreet,
			AddressCity:   row.AddressCity,
		}
		email := row.Email
		if email.Valid {
			customer.Email = &email.String
		}
		loanCustomers = append(loanCustomers, customer)
	}
	responseBody := GetAllLoanCustomersResponse{
		Data: &loanCustomers,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}
