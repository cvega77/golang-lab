package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alphaloan/vehicle/datastore"
)

type LoanSubmissionHandler struct {
	SubmissionStore datastore.LoanSubmissionStore
}

func NewLoanSubmissionHandler(
	submissionStore datastore.LoanSubmissionStore) *LoanSubmissionHandler {
	return &LoanSubmissionHandler{
		SubmissionStore: submissionStore,
	}
}

func (h *LoanSubmissionHandler) HandleGetAllLoanSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	loanSubmissionRows, err := h.SubmissionStore.GetAllLoanSubmissions()
	if err != nil {
		errMsg := "Failed to get all loan submissions"
		responseBodyErr := GetAllLoanSubmissionsResponse{
			ErrorMessage: &errMsg,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseBodyErr)
		return
	}

	loanSubmissions := make([]LoanSubmission, 0, len(loanSubmissionRows))
	for _, row := range loanSubmissionRows {
		loanSubmissions = append(loanSubmissions, LoanSubmission{
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
	responseBody := GetAllLoanSubmissionsResponse{
		Data: &loanSubmissions,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func (h *LoanSubmissionHandler) HandleSubmissionLoanById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	loanSubmissionId := r.URL.Query().Get("loan_submission_id")

	if !validateLoanSubmissionID(w, loanSubmissionId) {
		return
	}

	loanSubmissionRow, err := h.SubmissionStore.GetLoanSubmissionById(loanSubmissionId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get submission with loan_submission_id %s", loanSubmissionId), http.StatusInternalServerError)
		return
	}

	loanSubmission := LoanSubmission{
		SubmissionID:            loanSubmissionRow.SubmissionID,
		VehicleType:             loanSubmissionRow.VehicleType,
		VehicleBrand:            loanSubmissionRow.VehicleBrand,
		VehicleModel:            loanSubmissionRow.VehicleModel,
		VehicleLicenseNumber:    loanSubmissionRow.VehicleLicenseNumber,
		VehicleOdometer:         loanSubmissionRow.VehicleOdometer,
		ManufacturingYear:       loanSubmissionRow.ManufacturingYear,
		ProposedLoanAmount:      loanSubmissionRow.ProposedLoanAmount,
		ProposedLoanTenureMonth: loanSubmissionRow.ProposedLoanTenure,
		IsCommercialVehicle:     loanSubmissionRow.IsCommercialVehicle,
	}

	response := GetLoanSubmissionsByIdResponse{
		Data: &loanSubmission,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func validateLoanSubmissionID(w http.ResponseWriter, loanSubmissionId string) bool {
	if loanSubmissionId == "" {
		errMsg := "Missing submission_id query parameter"
		responseBodyErr := GetLoanSubmissionsByIdResponse{
			ErrorMessage: &errMsg,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return false
	}
	if !IsValidUUID(loanSubmissionId) {
		errMsg := "Invalid submission_id: " + loanSubmissionId
		responseBodyErr := GetLoanSubmissionsByIdResponse{
			ErrorMessage: &errMsg,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBodyErr)
		return false
	}
	return true
}
