package main

import (
	"encoding/json"
	"fmt"
	"golang-demo/api/data"
	"golang-demo/api/user"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	//-------connection database
	err := data.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	//-------setting up route
	router := mux.NewRouter()
	router.HandleFunc("/health", health)
	router.HandleFunc("/registration", user.UserRegistration).Methods("POST")
	router.HandleFunc("/user-list", user.UserList).Methods("GET")
	router.HandleFunc("/login", user.Login).Methods("POST")
	router.HandleFunc("/update-detail/{id}", user.UpdateDetail).Methods("UPDATE")
	router.HandleFunc("/update-password/{id}", user.UpdatePassword).Methods("UPDATE")
	router.HandleFunc("/image/{file-name}", GetImage)
	router.HandleFunc("/user-detail/{id}", user.UserDetail).Methods("GET")

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
