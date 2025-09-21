package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
)

type LoanSubmitHandler struct {
	CustomerStore   datastore.LoanCustomerStore
	SubmissionStore datastore.LoanSubmissionStore
}

func NewLoanSubmitHandler(
	customerStore datastore.LoanCustomerStore,
	submissionStore datastore.LoanSubmissionStore) *LoanSubmitHandler {
	return &LoanSubmitHandler{
		CustomerStore:   customerStore,
		SubmissionStore: submissionStore,
	}
}

func (h *LoanSubmitHandler) HandleSubmitLoan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method allowed", http.StatusMethodNotAllowed)
		return
	}

	var request LoanSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	loanCustomerRow := convertLoanCustomer(&request.Customer)

	upsertCustomerID, err := h.CustomerStore.UpsertCustomer(loanCustomerRow)
	if err != nil {
		http.Error(w, "Failed to upsert customer", http.StatusInternalServerError)
		return
	}

	loanSubmissionRow := convertLoanProposal(&request.ProposedLoad, upsertCustomerID)
	upsertSubmissionID, err := h.SubmissionStore.UpsertSubmission(loanSubmissionRow)
	if err != nil {
		http.Error(w, "Failed to upsert submission", http.StatusInternalServerError)
		return
	}

	response := LoanSubmitResponse{
		CustomerID:   &upsertCustomerID,
		SubmissionID: &upsertSubmissionID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
