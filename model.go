package main

import(
	"database/sql"
	"fmt"
	_"github.com/denisenkom/go-mssqldb"
	"log"
	"context"
	"errors"
	"os"
)

var db *sql.DB

var server=os.Getenv("DBSERVER")
var port = os.Getenv("DBPORT")
var user = os.Getenv("DBUSER")
var password = os.Getenv("DBPASSWORD")
var database = os.Getenv("DBNAME")
var connectionString string

func OpenSqlConnection()(*sql.DB,error){
	var err error
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database) 
	db, err = sql.Open("sqlserver", connString) 
	if err != nil { log.Fatal("Error creating connection pool: ", err.Error()) } 
	ctx := context.Background() 
	err = db.PingContext(ctx) 
	if err != nil { log.Fatal(err.Error()) } 
	fmt.Printf("Connected")
	return db,nil
}

func CloseSqlConnection(db *sql.DB){
	db.Close()
}

func SearchMovieByTitleInDb(title string)(MovieEntityDb,error){
	db,err := OpenSqlConnection()
	if err!=nil{
		return MovieEntityDb{},err
	}
	results,err:=db.Query("Select * from movies where title=@title",sql.Named("title",title))
	if err!=nil{
		fmt.Println("Error in retrieving data",err.Error())
		return MovieEntityDb{},err
	}
	defer CloseSqlConnection(db)
	var movieResult MovieEntityDb
	for results.Next(){
		results.Scan(&movieResult.Id,&movieResult.Title,&movieResult.Rating,&movieResult.ReleasedYear,&movieResult.Genres)
	}
	return movieResult,nil
}

func UpdateByIdInDb(title string,rating int,genres string) (bool,error){
	db,err := OpenSqlConnection()
	if err!=nil{
		return false,err
	}	
	row:=db.QueryRow("UPDATE Movies SET rating=@rating,genres=@genres WHERE title=@title", sql.Named("rating",rating),sql.Named("genres",genres),sql.Named("title",title))
	if row==nil{
		return false,errors.New("Some Problem Occured")
	}
	return true,nil
}

func SearchMovieByIdInDb(id string)(MovieEntityDb,error){
	db,err := OpenSqlConnection()
	if err!=nil{
		return MovieEntityDb{},err
	}
	results,err:=db.Query("Select * from movies where id=@id",sql.Named("id",id))
	if err!=nil{
		return MovieEntityDb{},err
	}
	defer CloseSqlConnection(db)
	var movieResult MovieEntityDb
	for results.Next(){
		results.Scan(&movieResult.Id,&movieResult.Title,&movieResult.Rating,&movieResult.ReleasedYear,&movieResult.Genres)
	}
	return movieResult,nil
}

func SearchMovieNotEqualRatingInDb(rating int)([]MovieEntityDb,error){
	var movies []MovieEntityDb
	db,err := OpenSqlConnection()
	if err!=nil{
		return nil,err
	}
	results,err:=db.Query("Select * from movies where rating!=@rating",sql.Named("rating",rating))
	if err!=nil{
		return nil,err
	}
	defer CloseSqlConnection(db)
	
	for results.Next(){
		var movieResult MovieEntityDb
		results.Scan(&movieResult.Id,&movieResult.Title,&movieResult.Rating,&movieResult.ReleasedYear,&movieResult.Genres)
		movies=append(movies,movieResult)
	}
	return movies,nil
}

func InsertMovieInDb(movie MovieEntityDb)(int,error) {
	db,err := OpenSqlConnection()
	ctx := context.Background() 
	err = db.PingContext(ctx) 
	if err != nil { log.Fatal(err.Error()) } 
	if err!=nil{
		return -1,err
	}
	insert, err := db.Query("INSERT INTO MOVIES(title,rating,genres,releasedyear,Id) VALUES(@title,@rating,@genres,@releasedyear,@id)",sql.Named("title",movie.Title),sql.Named("rating",movie.Rating),sql.Named("genres",movie.Genres),sql.Named("releasedyear",movie.ReleasedYear),sql.Named("id",movie.Id))
	if err != nil {
		return -1,err
	}
	defer insert.Close()
	defer db.Close()
	return 1,nil
}

func SearchMovieForReleasedYearsInDb(req TimeRangeRequest)([]MovieEntityDb,error){
	var movies []MovieEntityDb
	var err error
	var results *sql.Rows
	db,err := OpenSqlConnection()
	if err!=nil{
		return nil,err
	}
	if req.RangeReq==false{
		results,err=db.Query("Select * from movies where releasedYear=@releasedYear",sql.Named("releasedYear",req.ReqYear))
	}else{
		results,err=db.Query("Select * from movies where releasedYear>=@startYear and releasedYear<=@endYear",sql.Named("startYear",req.StartYear),sql.Named("endYear",req.EndYear))
	}
	if err!=nil{
		return nil,err
	}
	defer CloseSqlConnection(db)
	for results.Next(){
		var movieResult MovieEntityDb
		results.Scan(&movieResult.Id,&movieResult.Title,&movieResult.Rating,&movieResult.ReleasedYear,&movieResult.Genres)
		movies=append(movies,movieResult)
	}
	return movies,nil
}


func SearchMovieByGenresInDb(genres string)([]MovieEntityDb,error){
	var movies []MovieEntityDb
	db,err := OpenSqlConnection()
	if err!=nil{
		return nil,err
	}
	results,err:=db.Query("Select * from movies where genres in",genres)
	if err!=nil{
		return nil,err
	}
	defer CloseSqlConnection(db)
	
	for results.Next(){
		var movieResult MovieEntityDb
		results.Scan(&movieResult.Id,&movieResult.Title,&movieResult.Rating,&movieResult.ReleasedYear,&movieResult.Genres)
		movies=append(movies,movieResult)
	}
	return movies,nil
}

func GetAllMovies()([]MovieEntityDb,error){
	db,err := OpenSqlConnection()
	if err!=nil{
		return nil,err
	}
	results,err:=db.Query("Select * from movies")
	if err!=nil{
		return nil,err
	}
	defer CloseSqlConnection(db)
	var movies []MovieEntityDb
	for results.Next(){
		var movieResult MovieEntityDb
		results.Scan(&movieResult.Id,&movieResult.Title,&movieResult.Rating,&movieResult.ReleasedYear,&movieResult.Genres)
		movies=append(movies,movieResult)
	}
	return movies,nil
}