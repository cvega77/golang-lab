package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
)

type LoanCustomerHandler struct {
	CustomerStore   datastore.LoanCustomerStore
	SubmissionStore datastore.LoanSubmissionStore
}

func NewLoanCustomerHandler(
	customerStore datastore.LoanCustomerStore,
	submissionStore datastore.LoanSubmissionStore) *LoanCustomerHandler {
	return &LoanCustomerHandler{
		CustomerStore:   customerStore,
		SubmissionStore: submissionStore,
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

func (h *LoanCustomerHandler) HandleGetCustomerAndSubmissionById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	customerID := r.PathValue("customerID")
	if !IsValidUUID(customerID) {
		errMsg := "Invalid customer ID: " + customerID
		responseBodyErr := GetAllLoanCustomersResponse{
			ErrorMessage: &errMsg,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	loanCustomerWithAllSubmissionsRow, err := h.CustomerStore.GetCustomerByCustomerId(customerID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Failed to get Customer with CustomerId: %s", customerID), http.StatusInternalServerError)
		return
	}

	loanCustomer := LoanCustomer{
		CustomerID:    loanCustomerWithAllSubmissionsRow.LoanCustomerRow.CustomerID,
		IDCardNumber:  loanCustomerWithAllSubmissionsRow.LoanCustomerRow.IDCardNumber,
		FullName:      loanCustomerWithAllSubmissionsRow.LoanCustomerRow.FullName,
		BirthDate:     loanCustomerWithAllSubmissionsRow.LoanCustomerRow.BirthDate,
		PhoneNumber:   loanCustomerWithAllSubmissionsRow.LoanCustomerRow.PhoneNumber,
		MonthlyIncome: loanCustomerWithAllSubmissionsRow.LoanCustomerRow.MonthlyIncome,
		AddressStreet: loanCustomerWithAllSubmissionsRow.LoanCustomerRow.AddressStreet,
		AddressCity:   loanCustomerWithAllSubmissionsRow.LoanCustomerRow.AddressCity,
	}
	if loanCustomerWithAllSubmissionsRow.LoanCustomerRow.Email.Valid {
		loanCustomer.Email = &loanCustomerWithAllSubmissionsRow.LoanCustomerRow.Email.String
	}
	customerAndSubmissions := CustomerAndSubmissions{
		Customer: &loanCustomer,
	}
	if err != nil {
		fmt.Sprintf("Failed to get Loan Submission Customer with CustomerId: %s", customerID)
	}

	loadSubmissions := make([]LoanSubmission, 0, len(loanCustomerWithAllSubmissionsRow.LoanSubmissions))
	for _, row := range loanCustomerWithAllSubmissionsRow.LoanSubmissions {
		loadSubmissions = append(loadSubmissions, LoanSubmission{
			SubmissionID:            row.SubmissionID,
			VehicleType:             row.VehicleType,
			VehicleBrand:            row.VehicleBrand,
			VehicleModel:            row.VehicleModel,
			VehicleLicenseNumber:    row.VehicleLicenseNumber,
			VehicleOdometer:         row.VehicleOdometer,
			ManufacturingYear:       row.ManufacturingYear,
			ProposedLoanAmount:      row.ProposedLoanAmount,
			ProposedLoanTenureMonth: row.ProposedLoanTenure,
			IsCommercialVehicle:     row.IsCommercialVehicle,
		})
	}
	customerAndSubmissions.Submissions = &loadSubmissions

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customerAndSubmissions)
}
