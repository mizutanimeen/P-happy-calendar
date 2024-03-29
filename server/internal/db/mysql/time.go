package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/mizutanimeen/P-happiness-100-strikes/internal/db/model"
)

var (
	timeRecordTable           = os.Getenv("MYSQL_TIME_RECORD_TABLE")
	timeRecordID              = os.Getenv("MYSQL_TIME_RECORD_ID")
	timeRecordUserID          = os.Getenv("MYSQL_USERS_ID")
	timeRecordTime            = os.Getenv("MYSQL_TIME_RECORD_TIME")
	timeRecordInvestmentMoney = os.Getenv("MYSQL_TIME_RECORD_INVESTMENT_MONEY")
	timeRecordRecoveryMoney   = os.Getenv("MYSQL_TIME_RECORD_RECOVERY_MONEY")
)

func (s *Mysql) CreateTimeRecordTable() error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s int(16) AUTO_INCREMENT, %s varchar(32) NOT NULL, %s DATETIME NOT NULL, %s int(64) NOT NULL, %s int(64) NOT NULL, %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (%s), FOREIGN KEY (%s) REFERENCES %s(%s));",
		timeRecordTable, timeRecordID, timeRecordUserID, timeRecordTime, timeRecordInvestmentMoney, timeRecordRecoveryMoney, createAt, updateAt, timeRecordID, timeRecordUserID, userTable, userID)
	if _, err := s.DB.Exec(query); err != nil {
		return fmt.Errorf("error exec: %w", err)
	}

	return nil

}

func (s *Mysql) TimeRecordsGet(userID string, startDate time.Time, endDate time.Time) ([]*model.TimeRecord, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? AND %s >= ? AND %s <= ? ORDER BY %s ASC", timeRecordTable, timeRecordUserID, timeRecordTime, timeRecordTime, timeRecordTime)
	rows, err := s.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error query: %w", err)
	}
	defer rows.Close()

	var timeRecords []*model.TimeRecord
	for rows.Next() {
		var timeRecord model.TimeRecord
		if err := rows.Scan(&timeRecord.ID, &timeRecord.UserID, &timeRecord.Time, &timeRecord.InvestmentMoney, &timeRecord.RecoveryMoney, &timeRecord.Create_at, &timeRecord.Update_at); err != nil {
			return nil, fmt.Errorf("error scan: %w", err)
		}
		timeRecords = append(timeRecords, &timeRecord)
	}
	return timeRecords, nil
}

func (s *Mysql) TimeRecordGetByID(userID, id string) (*model.TimeRecord, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? AND %s = ?", timeRecordTable, timeRecordUserID, timeRecordID)
	row := s.DB.QueryRow(query, userID, id)

	var timeRecord model.TimeRecord
	if err := row.Scan(&timeRecord.ID, &timeRecord.UserID, &timeRecord.Time, &timeRecord.InvestmentMoney, &timeRecord.RecoveryMoney, &timeRecord.Create_at, &timeRecord.Update_at); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("error scan: %w", err)
	}
	return &timeRecord, nil
}

func (s *Mysql) TimeRecordCreate(userID string, time time.Time, investmentMoney int, recoveryMoney int) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s(%s, %s, %s, %s) VALUES(?,?,?,?)", timeRecordTable, timeRecordUserID, timeRecordTime, timeRecordInvestmentMoney, timeRecordRecoveryMoney)
	insert, err := s.DB.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("error prepare: %w", err)
	}

	r, err := insert.Exec(userID, time, investmentMoney, recoveryMoney)
	if err != nil {
		return -1, fmt.Errorf("error insert: %w", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error last insert id: %w", err)
	}
	return id, nil
}

func (s *Mysql) TimeRecordUpdate(userID string, id string, time time.Time, investmentMoney int, recoveryMoney int) error {
	query := fmt.Sprintf("UPDATE %s SET %s=?, %s=?, %s=? WHERE %s=? AND %s=?", timeRecordTable, timeRecordTime, timeRecordInvestmentMoney, timeRecordRecoveryMoney, timeRecordUserID, timeRecordID)
	update, err := s.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error prepare: %w", err)
	}

	if _, err := update.Exec(time, investmentMoney, recoveryMoney, userID, id); err != nil {
		return fmt.Errorf("error update: %w", err)
	}
	return nil
}

func (s *Mysql) TimeRecordDelete(userID, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s=? AND %s=?", timeRecordTable, timeRecordUserID, timeRecordID)
	delete, err := s.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error prepare: %w", err)
	}

	if _, err := delete.Exec(userID, id); err != nil {
		return fmt.Errorf("error delete: %w", err)
	}
	return nil
}
