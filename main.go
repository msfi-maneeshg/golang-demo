package main

import (
	"api"
	"database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	//-------connection database
	err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	//-------setting up route
	router := mux.NewRouter()
	router.HandleFunc("/health", health)
	router.HandleFunc("/registration", api.UserRegistration).Methods("POST")
	router.HandleFunc("/user-list", api.UserList).Methods("GET")
	router.HandleFunc("/login", api.Login).Methods("POST")
	router.HandleFunc("/update-detail/{id}", api.UpdateDetail).Methods("UPDATE")
	router.HandleFunc("/update-password/{id}", api.UpdatePassword).Methods("UPDATE")
	router.HandleFunc("/image/{file-name}", GetImage)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "UPDATE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	fmt.Println("Server is started...")
	log.Fatal(http.ListenAndServe(":8000", c.Handler(router)))

}

func health(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

//GetImage :
func GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var fileName = vars["file-name"]
	data, _ := ioutil.ReadFile("images/" + fileName)
	// w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
	r.Body.Close()
}
