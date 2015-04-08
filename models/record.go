package models

import (
	"net/http"
	"regexp"

	"github.com/mholt/binding"
	"github.com/nlf/boltons"
)

// Record is a single proxy record used for routing.
type Record struct {
	ID    string `json:"id"`
	Owner struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"owner"`
	FQDN            string `json:"fqdn"`
	HandlerHost     string `json:"handler_host"`
	HandlerPort     int    `json:"handler_port"`
	HandlerProtocol string `json:"handler_protocol"`
	UpdatedAt       int64  `json:"updated_at"`
	CreatedAt       int64  `json:"created_at"`
	Blacklist       bool   `json:"blacklist"`
}

// FindRecordsForOwner returns a list of all records for a given owner by their id.
func FindRecordsForOwner(db *boltons.DB, ID string) ([]Record, error) {
	records := []Record{}
	foundRecords := []Record{}
	if err := db.All(&records); err != nil {
		return foundRecords, err
	}
	for _, r := range records {
		if r.Owner.ID == ID {
			foundRecords = append(foundRecords, r)
		}
	}
	return foundRecords, nil
}

// FindRecordByFQDN returns a single record for the provided fqdn.
func FindRecordByFQDN(db *boltons.DB, fqdn string) (*Record, error) {
	record := Record{}
	records := []Record{}
	if err := db.All(&records); err != nil {
		return &record, err
	}
	for _, r := range records {
		if r.FQDN == fqdn {
			return &r, nil
		}
	}
	return &record, nil
}

// RecordRequest is used for JSON binding during a request to create a new record.
type RecordRequest struct {
	FQDN            string `json:"fqdn"`
	HandlerHost     string `json:"handler_host"`
	HandlerPort     int    `json:"handler_port"`
	HandlerProtocol string `json:"handler_protocol"`
}

func (r *RecordRequest) FieldMap() binding.FieldMap {
	return binding.FieldMap{}
}

func (r *RecordRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if r.FQDN == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"fqdn"},
			Message:    "fqdn must be a valid hostname",
		})
	}
	if ok, err := regexp.Match(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`, []byte(r.HandlerHost)); !ok || err != nil || r.HandlerHost == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"handler_host"},
			Message:    "handler_host must be a valid IP address",
		})
	}
	if r.HandlerPort < 0 || r.HandlerPort > 65535 {
		errs = append(errs, binding.Error{
			FieldNames: []string{"handler_port"},
			Message:    "handler_port must be a valid TCP port",
		})
	}
	if r.HandlerProtocol != "http" && r.HandlerProtocol != "https" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"fqdn"},
			Message:    "handler_protocol must be either http or https",
		})
	}
	return errs
}

// UpdateRecordRequest is used to perform JSON binding when updating a record.
type UpdateRecordRequest struct {
	FQDN            string `json:"fqdn"`
	HandlerHost     string `json:"handler_host"`
	HandlerPort     int    `json:"handler_port"`
	HandlerProtocol string `json:"handler_protocol"`
	Blacklist       bool   `json:"blacklist"`
	Owner           struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"owner"`
}

func (r *UpdateRecordRequest) FieldMap() binding.FieldMap {
	return binding.FieldMap{}
}

func (r *UpdateRecordRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if r.FQDN == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"fqdn"},
			Message:    "fqdn must be a valid hostname",
		})
	}
	if r.HandlerHost == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"handler_host"},
			Message:    "handler_host must be a valid IP address",
		})
	}
	if r.HandlerPort < 0 || r.HandlerPort > 65535 {
		errs = append(errs, binding.Error{
			FieldNames: []string{"handler_port"},
			Message:    "handler_port must be a valid TCP port",
		})
	}
	if r.HandlerProtocol != "http" && r.HandlerProtocol != "https" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"fqdn"},
			Message:    "handler_protocol must be either http or https",
		})
	}
	if r.Owner.ID == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"owner.id"},
			Message:    "owner.id is required",
		})
	}
	if r.Owner.Email == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"email"},
			Message:    "owner.email is required",
		})
	}
	return errs
}
