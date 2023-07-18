package jsonapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"log"
)

func setJsonHeader(rw http.ResponseWriter){
	rw.Header().Set("content-type", "application/json; charset:utf-8")
}

func fromJson[T any](body io.Reader, target T){

	buffer := new(bytes.Buffer)

	buffer.ReadFrom(body)

	json.Unmarshal(buffer.Bytes(), &target)
}

func returnJson[T any](rw http.ResponseWriter,withData func()(T,error)){
	
	setJsonHeader(rw)

	data, serverErr := withData()

	if serverErr != nil{

		rw.WriteHeader(500)

		serverErrorJson, err := json.Marshal(&serverErr)

		if err != nil{
			
			log.Println(err)
			return
		}

		rw.Write(serverErrorJson)
		return
	}

	dataJson, err := json.Marshal(&data)

	if err != nil{
		rw.WriteHeader(500)
		log.Println(err)

		return
	}

	rw.Write(dataJson)

}

func returnError(rw http.ResponseWriter, err error, code int){
	
	returnJson(rw, func() (interface{},error) {
		errorMessage := struct{
			msg string
		}{
			msg: err.Error(),
		}

		rw.WriteHeader(code)

		return errorMessage, nil
	})
}