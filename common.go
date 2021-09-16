package main

import(
	"encoding/json"
	"net/http"
	"strings"
)

func responseJSON(w http.ResponseWriter,status int,payload interface{}){
	response,err:=json.Marshal(payload)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func responseError(w http.ResponseWriter,status int,payload string){
	responseJSON(w,status,payload)
}

func ConvertDbEntityToResponse(movieEntityDb MovieEntityDb)(Movie){
	var movie Movie
	movie.Id=movieEntityDb.Id
	movie.Title=movieEntityDb.Title
	movie.ReleasedYear=movieEntityDb.ReleasedYear
	movie.Rating=movieEntityDb.Rating
	movie.Genres=strings.Split(movieEntityDb.Genres,",")
	return movie
}