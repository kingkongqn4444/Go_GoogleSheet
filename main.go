package main

import (
	"awesomeProject/src/mongoConnect"
	tripplan "awesomeProject/src/tripPlan"

	"github.com/gin-gonic/gin"
)

// gin-swagger middleware
// swagger embed files    // swagger embed files

func main() {
	//router := mux.NewRouter().StrictSlash(true)
	r := gin.Default()

	mongoConnect.ConnectMongo()
	//router.HandleFunc("/getListCompany, ListCompany).Methods("POST")
	//router.HandleFunc("/getListLocation, utilities.ListLocation).Methods("GET")
	//router.HandleFunc("/getPersion/{id}", mongoConnect.GetPersonByID).Methods("GET")
	//router.HandleFunc("/updatePersion/{id}", mongoConnect.UpdatePerson).Methods("PUT")
	//router.HandleFunc("/deletePersion/{id}", mongoConnect.DeletePerson).Methods("DELETE")
	//router.HandleFunc("/login", mongoConnect.LoginHandler).Methods("POST")
	// router.HandleFunc("/register", mongoConnect.RegisterHandler).Methods("POST")
	//router.HandleFunc("/uploadAvatar", mongoConnect.UploadFile)

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/getListCompany",utilities.ListCompany)
	// r.GET("/getListLocation", utilities.ListLocation)
	// r.POST("/getListCompany", utilities.ListCompany)
	// r.POST("/createPersion", mongoConnect.CreatePerson)
	// r.POST("/loginPersion", mongoConnect.LoginHandler)
	// r.POST("/createCategory", mongoConnect.CreateCategory)
	r.POST("/registerPersion", mongoConnect.RegisterHandler)
	// r.GET("/profilePersion", mongoConnect.ProfileHandler)
	// r.GET("/getListCategory", mongoConnect.GetListAllCategory)

	// four square
	r.GET("/autocomplete/:keyword", tripplan.AutoCompletePlace)

	r.Run()

	//log.Fatal(http.ListenAndServe(":90", router))
}
