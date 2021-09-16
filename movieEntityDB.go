package main

type MovieEntityDb struct{
	Id 				string	`json:"id"`
	Title 			string 	`json:"title"`
	Rating  		int     `json:"rating"`
	ReleasedYear	string  `json:"releasedYear"`
	Genres          string  `json:"genres"`
}