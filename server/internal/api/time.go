package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mizutanimeen/P-happiness-100-strikes/internal/db"
	"github.com/mizutanimeen/P-happiness-100-strikes/internal/db/model"
)

func TimeRecordsGet(DB db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(CK_USERID).(string)

		startTime, endTime, status, err := getStartEndDateQuery(r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		timeRecords, err := DB.TimeRecordsGet(userID, startTime, endTime)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		records := make(map[string][]model.TimeRecord)
		for _, record := range timeRecords {
			date := record.Time.Format("2006-01-02")
			records[date] = append(records[date], *record)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(records)
	}
}

func TimeRecordsGetByID(DB db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(CK_USERID).(string)

		timeRecordID := chi.URLParam(r, "times_id")

		timeRecord, err := DB.TimeRecordGetByID(userID, timeRecordID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if timeRecord == nil {
			http.Error(w, "Record not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(timeRecord)
	}
}

type createTimeRecordRequest struct {
	DateTime        string `json:"date_time"`
	InvestmentMoney int    `json:"investment_money"`
	RecoveryMoney   int    `json:"recovery_money"`
}

func TimeRecordCreate(DB db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(CK_USERID).(string)

		var timeRecordReq createTimeRecordRequest
		if err := json.NewDecoder(r.Body).Decode(&timeRecordReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		dateTime, err := time.Parse("2006-01-02T15:04:05", timeRecordReq.DateTime)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		id, err := DB.TimeRecordCreate(userID, dateTime, timeRecordReq.InvestmentMoney, timeRecordReq.RecoveryMoney)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		text := fmt.Sprintf(`{"id":%d,"message":"created"}`, id)
		w.Write([]byte(text))
	}
}

type updateTimeRecordRequest struct {
	ID              string `json:"id"`
	DateTime        string `json:"date_time"`
	InvestmentMoney int    `json:"investment_money"`
	RecoveryMoney   int    `json:"recovery_money"`
}

func TimeRecordUpdate(DB db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(CK_USERID).(string)

		var timeRecordReq updateTimeRecordRequest
		if err := json.NewDecoder(r.Body).Decode(&timeRecordReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		dateTime, err := time.Parse("2006-01-02T15:04:05", timeRecordReq.DateTime)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := DB.TimeRecordUpdate(userID, timeRecordReq.ID, dateTime, timeRecordReq.InvestmentMoney, timeRecordReq.RecoveryMoney); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"updated"}`))
	}
}

func TimeRecordDelete(DB db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(CK_USERID).(string)

		timeRecordID := r.URL.Query().Get("time_record_id")
		if timeRecordID == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if err := DB.TimeRecordDelete(userID, timeRecordID); err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"deleted"}`))
	}
}
