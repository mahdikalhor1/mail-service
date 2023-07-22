package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"mailService/mailDB"
	"net/http"
	"errors"
)

func setJsonHeader(rw http.ResponseWriter) {
	rw.Header().Set("content-type", "application/json; charset:utf-8")
}

func fromJson[T any](body io.Reader, target T) {

	buffer := new(bytes.Buffer)

	buffer.ReadFrom(body)

	json.Unmarshal(buffer.Bytes(), &target)
}

func returnJson[T any](rw http.ResponseWriter, withData func() (T, error)) {

	setJsonHeader(rw)

	data, serverErr := withData()

	if serverErr != nil {

		rw.WriteHeader(500)

		serverErrorJson, err := json.Marshal(&serverErr)

		if err != nil {

			log.Println(err)
			return
		}

		rw.Write(serverErrorJson)
		return
	}

	dataJson, err := json.Marshal(&data)

	if err != nil {
		rw.WriteHeader(500)
		log.Println(err)

		return
	}

	rw.Write(dataJson)

}

func returnError(rw http.ResponseWriter, err error, code int) {

	returnJson(rw, func() (interface{}, error) {
		errorMessage := struct {
			msg string
		}{
			msg: err.Error(),
		}

		rw.WriteHeader(code)

		return errorMessage, nil
	})
}


func GetEmail(db *sql.DB) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){

		if req.Method != "Get"{
			return
		}

		emailEntry := mailDB.EmailEntry{}

		fromJson(req.Body, &emailEntry)


		returnJson(rw, func()(interface{}, error){
			log.Printf("Json Get Email: %v\n", emailEntry.Email)

			return mailDB.GetEmail(db, emailEntry.Email)
		})
	})
}


func InsertEmail(db *sql.DB) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){

		if req.Method != "POST"{
			return
		}

		emailEntry := mailDB.EmailEntry{}

		fromJson(req.Body, &emailEntry)

		err := mailDB.InsertEmail(db, emailEntry.Email)

		if err != nil{
			returnError(rw, err, 400)
			return
		}

		returnJson(rw, func()(interface{}, error){
			log.Printf("Json Insert Email: %v\n", emailEntry.Email)

			return mailDB.GetEmail(db, emailEntry.Email)
		})
	})
}

func UpdateEmail(db *sql.DB) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){

		if req.Method != "PUT"{
			return
		}

		emailEntry := mailDB.EmailEntry{}

		fromJson(req.Body, &emailEntry)

		if err := mailDB.UpdateEmail(db, &emailEntry); err != nil{
			returnError(rw, err, 400)
			return
		}

		returnJson(rw, func()(interface{}, error){

			log.Printf("Json Update Email: %v\n", emailEntry.Email)

			return mailDB.GetEmail(db, emailEntry.Email)
		})
	})
}


func DeleteEmail(db *sql.DB) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){

		if req.Method != "POST"{
			return
		}

		emailEntry := mailDB.EmailEntry{}

		fromJson(req.Body, &emailEntry)

		err := mailDB.DeleteEmail(db, emailEntry.Email)

		if err != nil{
			returnError(rw, err, 400)
			return
		}

		returnJson(rw, func()(interface{}, error){
			log.Printf("Json Delete Email: %v\n", emailEntry.Email)

			return mailDB.GetEmail(db, emailEntry.Email)
		})
	})
}

func GetEmailBatch(db *sql.DB) http.Handler{

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){

		if req.Method != "GET"{
			return
		}

		batchParams := mailDB.GetEmailBatchParams{}

		fromJson(req.Body, &batchParams)

		if batchParams.Count <= 0 || batchParams.Page <= 0{
			returnError(rw, errors.New("Email count and page are requiered and must be grater than 0."), 400)
			return
		}

		returnJson(rw, func()(interface{}, error){
			log.Printf("Json GetEmailBatch count: %v page: %v\n", batchParams.Count, batchParams.Page)

			return mailDB.GetEmailBatch(db, batchParams)
		})
	})
}


func Serve(db *sql.DB, bind string){

	http.Handle("email/insert", InsertEmail(db))
	http.Handle("email/get", GetEmail(db))
	http.Handle("email/update", UpdateEmail(db))
	http.Handle("email/get-batch", GetEmailBatch(db))
	http.Handle("email/delete", DeleteEmail(db))

	err := http.ListenAndServe(bind, nil)

	if err != nil{
		log.Fatalf("Server Error : %v", err)
	}
}