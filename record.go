package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/mholt/binding"
	"github.com/nlf/boltons"
)

type Record struct {
	ID    string `json:"id"`
	Owner struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"owner"`
	FQDN            string   `json:"fqdn"`
	HandlerHost     string   `json:"handler_host"`
	HandlerPort     int      `json:"handler_port"`
	HandlerProtocol string   `json:"handler_protocol"`
	UpdatedAt       int64    `json:"updated_at"`
	CreatedAt       int64    `json:"created_at"`
	Clients         []string `json:"clients"`
	Blacklist       bool     `json:"blacklist"`
}

type UpdateRecordRequest struct {
	FQDN            string   `json:"fqdn"`
	HandlerHost     string   `json:"handler_host"`
	HandlerPort     int      `json:"handler_port"`
	HandlerProtocol string   `json:"handler_protocol"`
	Blacklist       bool     `json:"blacklist"`
	Clients         []string `json:"clients"`
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

type NewRecordRequest struct {
	FQDN            string `json:"fqdn"`
	HandlerHost     string `json:"handler_host"`
	HandlerPort     int    `json:"handler_port"`
	HandlerProtocol string `json:"handler_protocol"`
}

func (r *NewRecordRequest) FieldMap() binding.FieldMap {
	return binding.FieldMap{}
}

func (r *NewRecordRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
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
	return errs
}

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

func RecordCreateHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		user := &User{}
		user = context.Get(req, "user").(*User)
		recordReq := &NewRecordRequest{}
		if err := binding.Bind(req, recordReq); err.Handle(w) {
			return
		}
		existing, err := FindRecordByFQDN(app.DB, recordReq.FQDN)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
			log.Println(err)
			return
		}
		if existing.ID != "" {
			app.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "fqdn must be unique across the application"})
			return
		}
		record := &Record{}
		record.Owner.Email = user.Email
		record.Owner.ID = user.ID
		now := time.Now().Unix()
		record.CreatedAt = now
		record.UpdatedAt = now
		record.Clients = []string{}
		copier.Copy(record, recordReq)
		if err := app.DB.Save(record); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
			log.Println(err)
			return
		}

		app.Render.JSON(w, http.StatusCreated, record)
	}
}

func RecordIndexHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		records := []Record{}
		if err := app.DB.All(&records); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting records from the database"})
			log.Println(err)
			return
		}
		app.Render.JSON(w, http.StatusOK, records)
	}
}

func RecordShowHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		record := &Record{ID: id}
		if ok, err := app.DB.Exists(record); err != nil || !ok {
			app.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if err := app.DB.Get(record); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting record from the database"})
			log.Println(err)
			return
		}
		app.Render.JSON(w, http.StatusOK, record)
	}
}
func RecordDeleteHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		record := &Record{ID: id}
		if ok, err := app.DB.Exists(record); err != nil || !ok {
			app.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if err := app.DB.Delete(record); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error deleting record from the database"})
			log.Println(err)
			return
		}
		app.Render.JSON(w, http.StatusNoContent, nil)
	}
}

func RecordUpdateHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		record := &Record{ID: id}
		if ok, err := app.DB.Exists(record); err != nil || !ok {
			app.Render.JSON(w, http.StatusNotFound, nil)
			return
		}

		updateReq := &UpdateRecordRequest{}
		if err := binding.Bind(req, updateReq); err.Handle(w) {
			return
		}

		if err := app.DB.Get(record); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting record from the database"})
			log.Println(err)
			return
		}

		if updateReq.FQDN != record.FQDN {
			existing, err := FindRecordByFQDN(app.DB, updateReq.FQDN)
			if err != nil {
				app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
				log.Println(err)
				return
			}
			if existing.ID != "" {
				app.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "fqdn must be unique across the application"})
				return
			}
		}

		if updateReq.Owner.ID != record.ID {
			if ok, err := app.DB.Exists(&User{ID: updateReq.Owner.ID}); err != nil || !ok {
				app.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "owner does not exist"})
				return
			}
		}

		copier.Copy(record, updateReq)
		record.UpdatedAt = time.Now().Unix()
		if err := app.DB.Save(record); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the record"})
			log.Println(err)
			return
		}
		app.Render.JSON(w, http.StatusOK, record)
	}
}
