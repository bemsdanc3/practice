package entities

import "time"

type Segment struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement"`
	TicketNumber        string    `gorm:"size:255;not null;index:idx_ticket_number"`
	SegmentNumber       int       `gorm:"not null;index:idx_ticket_segment,unique"`
	OperationType       string    `gorm:"size:50;not null"`
	OperationTime       time.Time `gorm:"type:timestamptz;not null"`
	OperationTimeZone   int16     `gorm:"not null"`
	OperationPlace      string    `gorm:"size:255"`
	IsRefunded          bool      `gorm:"default:false"`
	PassengerName       string    `gorm:"size:255"`
	PassengerSurname    string    `gorm:"size:255"`
	PassengerPatronymic string    `gorm:"size:255"`
	DocType             string    `gorm:"size:10"`
	DocNumber           string    `gorm:"size:20;not null;index:idx_doc_number"`
	Birthdate           time.Time
	Gender              string `gorm:"size:1"`
	PassengerType       string `gorm:"size:50"`
	TicketType          int
	AirlineCode         string    `gorm:"size:10;not null"`
	FlightNum           int       `gorm:"not null"`
	DepartPlace         string    `gorm:"size:10;not null"`
	DepartDatetime      time.Time `gorm:"type:timestamptz;not null"`
	ArrivePlace         string    `gorm:"size:10;not null"`
	ArriveDatetime      time.Time `gorm:"type:timestamptz;not null"`
	PNRID               string    `gorm:"size:255"`
}

func (Segment) TableName() string {
	return "segments"
}
