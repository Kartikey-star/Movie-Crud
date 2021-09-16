package main

import(
    "net/http"
	"fmt"
    "github.com/gorilla/mux"
)

func main(){
	router:=mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/",handler).Methods("GET")
	router.HandleFunc("/api/InsertMovie",InsertMovie).Methods("POST")
	router.HandleFunc("/api/UpdateMovie",UpdateMovie).Methods("PUT")
	router.HandleFunc("/api/SearchMovie/{id}",SearchMovieById).Methods("GET")
	router.HandleFunc("/api/SearchMovieByRating/{rating}",SearchMovieByRating).Methods("GET")
	router.HandleFunc("/api/SearchMovieGenres",SearchMovieGenres).Methods("POST")
	router.HandleFunc("/api/SearchMovieTimeRange",SearchMovieTimeRange).Methods("POST")
	fmt.Println(http.ListenAndServe(":8080", router))
	http.ListenAndServe(":8080", router)
}

func handler(w http.ResponseWriter, r *http.Request){
	responseJSON(w,http.StatusOK,"listening")

}