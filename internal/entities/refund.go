package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type RefundRequest struct {
	OperationType  string    `json:"operation_type" validate:"required,oneof=refund"`
	OperationTime  time.Time `json:"operation_time" validate:"required"`
	OperationPlace string    `json:"operation_place" validate:"required"`
	TicketNumber   string    `json:"ticket_number" validate:"required"`
}

func (r *RefundRequest) UnmarshalJSON(data []byte) error {
	type Alias RefundRequest
	aux := &struct {
		OperationTime string `json:"operation_time"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, aux.OperationTime)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02T15:04+03:00", aux.OperationTime)
		if err != nil {
			return fmt.Errorf("invalid time format: %v", err)
		}
	}

	r.OperationTime = parsedTime
	return nil
}
