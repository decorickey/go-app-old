package models

import (
	"fmt"
	"log"
	"time"
)

func initProgram() {
	var program Program
	tableName := program.getTableName()
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			studio_name TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			performer TEXT NOT NULL,
			vol TEXT NOT NULL,
			debug_date TEXT,
			PRIMARY KEY (studio_name, start_time)
		)
	`, tableName)

	_, err := DbConnection.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
}

type Program struct {
	StudioName string    `json:"studio_name,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	EndTime    time.Time `json:"end_time,omitempty"`
	Performer  string    `json:"performer,omitempty"`
	Vol        string    `json:"vol,omitempty"`
	DebugDate  string    `json:"debug_date,omitempty"`
}

func (p Program) getTableName() string {
	return "program"
}

func (p Program) Create() error {
	query := fmt.Sprintf(`
		INSERT INTO %s (studio_name, start_time, end_time, performer, vol, debug_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`, Program{}.getTableName())

	_, err := DbConnection.Exec(query,
		p.StudioName, p.StartTime.Format(time.RFC3339), p.EndTime.Format(time.RFC3339), p.Performer, p.Vol, p.DebugDate)
	return err
}

func (p Program) Save() error {
	query := fmt.Sprintf(`
		UPDATE %s SET end_time = ?, performer = ?, vol = ?, debug_date = ?
		WHERE studio_name = ? AND start_time = ?
	`, Program{}.getTableName())

	_, err := DbConnection.Exec(query,
		p.EndTime.Format(time.RFC3339), p.Performer, p.Vol, p.DebugDate,
		p.StudioName, p.StartTime.Format(time.RFC3339))
	return err
}

// GetProgram PRIMARY KEYでSELECT
func GetProgram(studioName string, startTime time.Time) *Program {
	query := fmt.Sprintf(`
		SELECT * FROM %s WHERE studio_name = ? AND start_time = ?
	`, Program{}.getTableName())

	row := DbConnection.QueryRow(query, studioName, startTime.Format(time.RFC3339))

	var p Program
	err := row.Scan(&p.StudioName, &p.StartTime, &p.EndTime, &p.Performer, &p.Vol, &p.DebugDate)
	if err != nil {
		return nil
	}
	return &p
}

func GetAllProgram() (dfProgram *DataFrameProgram, err error) {
	query := fmt.Sprintf("SELECT * FROM %s", Program{}.getTableName())

	rows, err := DbConnection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dfProgram = &DataFrameProgram{}
	for rows.Next() {
		p := Program{}
		rows.Scan(&p.StudioName, &p.StartTime, &p.EndTime, &p.Performer, &p.Vol, &p.DebugDate)
		dfProgram.Programs = append(dfProgram.Programs, &p)
	}
	return dfProgram, rows.Err()
}

func GetProgramByPerformer(performer string) (dfProgram *DataFrameProgram, err error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s WHERE performer = ? ORDER BY start_time
	`, Program{}.getTableName())

	rows, err := DbConnection.Query(query, performer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dfProgram = &DataFrameProgram{}
	for rows.Next() {
		p := Program{}
		rows.Scan(&p.StudioName, &p.StartTime, &p.EndTime, &p.Performer, &p.Vol, &p.DebugDate)
		dfProgram.Programs = append(dfProgram.Programs, &p)
	}
	return dfProgram, rows.Err()
}
