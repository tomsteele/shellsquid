package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/mholt/binding"
	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/models"
)

func isSameAsListener(listener string, req *models.RecordRequest) bool {
	return listener == req.HandlerHost+":"+strconv.Itoa(req.HandlerPort)
}

// CreateRecord handles a request to create a new record.
func CreateRecord(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		user := &models.User{}
		user = context.Get(req, "user").(*models.User)
		recordReq := &models.RecordRequest{}
		if err := binding.Bind(req, recordReq); err.Handle(w) {
			return
		}
		if isSameAsListener(server.Config.Proxy.HTTP.Listener, recordReq) {
			server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Handler Host and Handler Port must not be the same as HTTP Listener"})
			return
		}
		if isSameAsListener(server.Config.Proxy.SSL.Listener, recordReq) {
			server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Handler Host and Handler Port must not be the same as SSL Listener"})
			return
		}
		existing, err := models.FindRecordByFQDN(server.DB, recordReq.FQDN)
		if err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
			log.Println(err)
			return
		}
		if existing.ID != "" {
			server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "fqdn must be unique across the application"})
			return
		}
		now := time.Now().Unix()
		record := &models.Record{
			CreatedAt: now,
			UpdatedAt: now,
		}
		record.Owner.Email = user.Email
		record.Owner.ID = user.ID
		if err := copier.Copy(record, recordReq); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
			return
		}

		if err := server.DB.Save(record); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
			return
		}

		server.Render.JSON(w, http.StatusCreated, record)
	}
}

// IndexRecord handles a request to return a list of all records.
func IndexRecord(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		records := []models.Record{}
		if err := server.DB.All(&records); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting records from the database"})
			return
		}
		server.Render.JSON(w, http.StatusOK, records)
	}
}

// ShowRecord handles a request to return a single record provided by the mux parameter id.
func ShowRecord(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		record := &models.Record{ID: id}
		if ok, err := server.DB.Exists(record); err != nil || !ok {
			server.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if err := server.DB.Get(record); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting record from the database"})
			log.Println(err)
			return
		}
		server.Render.JSON(w, http.StatusOK, record)
	}
}

// DeleteRecord handles a request to delete a single record provided the mux parameter id.
func DeleteRecord(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		record := &models.Record{ID: id}
		if ok, err := server.DB.Exists(record); err != nil || !ok {
			server.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if err := server.DB.Delete(record); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error deleting record from the database"})
			log.Println(err)
			return
		}
		server.Render.JSON(w, http.StatusNoContent, nil)
	}
}

// UpdateRecord handles a request to update a single record provided the mux parameter id.
func UpdateRecord(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		record := &models.Record{ID: id}
		if ok, err := server.DB.Exists(record); err != nil || !ok {
			server.Render.JSON(w, http.StatusNotFound, nil)
			return
		}

		updateReq := &models.UpdateRecordRequest{}
		if err := binding.Bind(req, updateReq); err.Handle(w) {
			return
		}

		if err := server.DB.Get(record); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting record from the database"})
			log.Println(err)
			return
		}

		if updateReq.FQDN != record.FQDN {
			existing, err := models.FindRecordByFQDN(server.DB, updateReq.FQDN)
			if err != nil {
				server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the record to the database"})
				log.Println(err)
				return
			}
			if existing.ID != "" {
				server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "fqdn must be unique across the application"})
				return
			}
		}

		if updateReq.Owner.ID != record.ID {
			if ok, err := server.DB.Exists(&models.User{ID: updateReq.Owner.ID}); err != nil || !ok {
				server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "owner does not exist"})
				return
			}
		}

		if err := copier.Copy(record, updateReq); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the record"})
			return
		}
		record.UpdatedAt = time.Now().Unix()
		if err := server.DB.Save(record); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the record"})
			log.Println(err)
			return
		}
		server.Render.JSON(w, http.StatusOK, record)
	}
}
