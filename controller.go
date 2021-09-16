package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"github.com/eefret/gomdb"
	"strconv"
	"encoding/json"
	"strings"
	"fmt"
	"os"
)

type GenresRequest struct{
	Genres []string `json:"genres"`
}
type TimeRangeRequest struct{
	RangeReq 	bool 	`json:"rangeReq"`
	ReqYear 	string	`json:"reqYear"`
	StartYear   string  `json:"startYear"`
	EndYear     string   `json:"endYear"`
}
type UpdateMovieRequest struct {
	Title 	string 	`json:"title"`
	Rating 	int 	`json:"rating"`
	Genres  []string  `json:"genres"`
}

type InsertMovieRequest struct{
	Title 	string 	`json:"title"`
}
var  YOUR_API_KEY = os.Getenv("API_KEY")

//Insert movie in database if not present call imdb api and then insert movie
func InsertMovie(w http.ResponseWriter, r *http.Request){
	var req InsertMovieRequest
	json.NewDecoder(r.Body).Decode(&req)
	movie,err:=SearchMovieByTitleInDb(req.Title) 
	if err!=nil{
		fmt.Println(err)
	}
	if movie.Id != "" {
		responseJSON(w,http.StatusOK,movie.Id)
		return
	}
	api := gomdb.Init(YOUR_API_KEY)
	queryData:=&gomdb.QueryData{
		Title: req.Title,
	}
	movieresult,err:=api.MovieByTitle(queryData)
	if err!=nil{
		responseError(w,http.StatusInternalServerError,"Some problem occured")
	}
	movieEntity:=MovieEntityDb{
		Title: movieresult.Title,
		Id: movieresult.ImdbID,
		ReleasedYear:movieresult.Year,
		Genres:movieresult.Type,
	}
	_,err=InsertMovieInDb(movieEntity)
	if err==nil{
		movie,err=SearchMovieByTitleInDb(req.Title) 
		if err!=nil{
			responseError(w,http.StatusInternalServerError,"Some problem occured")
			return
		}
		if movie.Id != "" {
			responseJSON(w,http.StatusOK,movie.Id)
			return
		}
	}else{
		responseError(w,http.StatusInternalServerError,"Some problem occured")
		return
	}
	
}

//update the movie rating and genre
func UpdateMovie(w http.ResponseWriter, r *http.Request){
	var req UpdateMovieRequest
	json.NewDecoder(r.Body).Decode(&req)
	genresString:=strings.Join(req.Genres,",")
	result,err:=UpdateByIdInDb(req.Title,req.Rating,genresString)
	if result{
		responseJSON(w,http.StatusCreated,"Record Succesfully updated")
		return
	}else{
		responseError(w,http.StatusInternalServerError,err.Error())
		return
	}
}

func SearchMovieById(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id:= vars["id"]
	movieEntityDb,err:=SearchMovieByIdInDb(id)
	if err!=nil{
		responseError(w,http.StatusInternalServerError,err.Error())
		return
	}
	movieResponse:=ConvertDbEntityToResponse(movieEntityDb)
	responseJSON(w,http.StatusOK,movieResponse)
}


func SearchMovieByRating(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	rating,err:= strconv.Atoi(vars["rating"])
	if err!=nil{
		responseError(w,http.StatusInternalServerError,err.Error())
	}
	var movieResponse []Movie
	moviesDbEntity,err:=SearchMovieNotEqualRatingInDb(rating)
	if err!=nil{
		responseError(w,http.StatusInternalServerError,err.Error())
	}
	for _,movieDb:=range moviesDbEntity{
		var movie Movie
		movie=ConvertDbEntityToResponse(movieDb)
		movieResponse=append(movieResponse,movie)
	}
	responseJSON(w,http.StatusOK,movieResponse)
}


func SearchMovieGenres(w http.ResponseWriter, r *http.Request){
	var genresRequest GenresRequest
	json.NewDecoder(r.Body).Decode(&genresRequest)
	
	moviesDbEntity,err:=GetAllMovies()
	if err!=nil{
		responseError(w,http.StatusInternalServerError,err.Error())
	}
	var movieResponse []Movie
	for _,movieDb:=range moviesDbEntity{
		var movie Movie
		movie=ConvertDbEntityToResponse(movieDb)
		movieResponse=append(movieResponse,movie)
	}
	var finalMovieResult []Movie
	for _,movie:=range movieResponse{
		if CheckIfGenresInEntity(movie.Genres, genresRequest.Genres){
			finalMovieResult=append(finalMovieResult,movie)
		}
	}
	responseJSON(w,http.StatusOK,finalMovieResult)
}


func SearchMovieTimeRange(w http.ResponseWriter, r *http.Request){
	var timeRangeRequest TimeRangeRequest
	json.NewDecoder(r.Body).Decode(&timeRangeRequest)
	movies,err:=SearchMovieForReleasedYearsInDb(timeRangeRequest)
	if err!=nil{
		responseError(w,http.StatusInternalServerError,err.Error())
	}
	var movieResponse []Movie
	for _,movieDb:=range movies{
		var movie Movie
		movie=ConvertDbEntityToResponse(movieDb)
		movieResponse=append(movieResponse,movie)
	}
	responseJSON(w,http.StatusOK,movieResponse)
}

func CheckIfGenresInEntity(genresInDb []string,req []string)bool{
	for _,genres:=range genresInDb{
		for _,reqType:=range req{
			if reqType==genres{
				return true
			}
		}
	}
	return false
}