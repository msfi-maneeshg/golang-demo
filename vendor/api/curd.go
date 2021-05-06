package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Registration :
type Registration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

//UserRegistration :
func UserRegistration(w http.ResponseWriter, r *http.Request) {
	var objRegistration Registration
	var err error
	err = json.NewDecoder(r.Body).Decode(&objRegistration)
	if err != nil {
		fmt.Println("Error:" + err.Error())
	}

	finalOutput, err := json.Marshal(objRegistration)
	if err != nil {
		fmt.Println("Error:" + err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalOutput)

}
