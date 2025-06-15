package main

import (
	"github.com/dancankarani/safa/database"
	"github.com/dancankarani/safa/endpoints"
	"github.com/dancankarani/safa/models"
	//"github.com/dancankarani/safa/services"
)

func main (){
	database.ConnectDB()
	//go services.SendEmail("karanidancan120@gmail.com","dsa","dsad")
	// Initialize your application here
	models.MigrateDb()
	endpoints.RegisterEndpoint()
}