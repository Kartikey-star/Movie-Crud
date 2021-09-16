package main

type Movie struct{
	Title 			string 		`json:"title"`
	ReleasedYear 	string			`json:"releasedYear"`
	Rating          int			`json:"rating"`
	Id              string        `json:"id"`
	Genres          []string    `json:"genres"`
}