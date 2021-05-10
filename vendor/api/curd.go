package api

import (
	"common"
	"database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//UserInformation :
type UserInformation struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

//LoginDetails :
type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//UserRegistration :
func UserRegistration(w http.ResponseWriter, r *http.Request) {
	var objRegistration UserInformation
	var err error

	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objRegistration)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	//-------validate data

	if objRegistration.Name == "" {
		common.APIResponse(w, http.StatusBadRequest, "Name can not be empty")
		return
	}

	if objRegistration.Email == "" {
		common.APIResponse(w, http.StatusBadRequest, "Email can not be empty")
		return
	}

	if objRegistration.Phone == "" {
		common.APIResponse(w, http.StatusBadRequest, "Phone number can not be empty")
		return
	}
	_, err = strconv.Atoi(objRegistration.Phone)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Invalid data of phone number")
		return
	}

	if objRegistration.Password == "" {
		common.APIResponse(w, http.StatusBadRequest, "Password can not be empty")
		return
	}

	isExist, err := isEmailExist(objRegistration.Email)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error when checking availability of new email. Error:"+err.Error())
		return
	}
	if isExist {
		common.APIResponse(w, http.StatusBadRequest, "This email is already registered.")
		return
	}

	err = insertNewUser(objRegistration)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error inserting new user info. Error:"+err.Error())
		return
	}
	common.APIResponse(w, http.StatusOK, "User registered successfully.")
	return

}

//UserList :
func UserList(w http.ResponseWriter, r *http.Request) {
	userList, err := getUserList("")
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while getting user list")
		return
	}
	common.APIResponse(w, http.StatusOK, userList)
	return
}

//Login :
func Login(w http.ResponseWriter, r *http.Request) {
	var objLoginDetails LoginDetails
	var err error

	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objLoginDetails)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	//-------validate data
	if objLoginDetails.Email == "" {
		common.APIResponse(w, http.StatusBadRequest, "Email can not be empty")
		return
	}

	if objLoginDetails.Password == "" {
		common.APIResponse(w, http.StatusBadRequest, "Password can not be empty")
		return
	}

	users, err := getUserList(objLoginDetails.Email)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error when checking login details. Error:"+err.Error())
		return
	}
	if len(users) == 0 {
		common.APIResponse(w, http.StatusNotFound, "This email is not registered.")
		return
	}

	if users[0].Password != objLoginDetails.Password {
		common.APIResponse(w, http.StatusBadRequest, "Invalid login details.")
		return
	}
	users[0].Password = ""
	common.APIResponse(w, http.StatusOK, users[0])
	return

}

//UpdateDetail :
func UpdateDetail(w http.ResponseWriter, r *http.Request) {
	var objRegistration UserInformation
	var err error
	vars := mux.Vars(r)
	var id = vars["id"]

	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objRegistration)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	//-------validate data

	if objRegistration.Name == "" {
		common.APIResponse(w, http.StatusBadRequest, "Name can not be empty")
		return
	}

	if objRegistration.Phone == "" {
		common.APIResponse(w, http.StatusBadRequest, "Phone number can not be empty")
		return
	}
	_, err = strconv.Atoi(objRegistration.Phone)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Invalid data of phone number")
		return
	}

	isExist, err := isIDExist(id)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while checking is exist or not."+err.Error())
		return
	}

	if !isExist {
		common.APIResponse(w, http.StatusBadRequest, "Invalid data.")
		return
	}

	err = updateUserDetails(objRegistration, id)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while updating user info. Error:"+err.Error())
		return
	}

	common.APIResponse(w, http.StatusOK, "User detail has been updated.")
	return
}

//---------data functions
func isEmailExist(email string) (bool, error) {
	var emailID string
	sqlStr := "SELECT email FROM users WHERE email = ?"

	err := database.DemoDB.QueryRow(sqlStr, email).Scan(&emailID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func isIDExist(id string) (bool, error) {
	var emailID string
	sqlStr := "SELECT email FROM users WHERE id = ?"

	err := database.DemoDB.QueryRow(sqlStr, id).Scan(&emailID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func insertNewUser(objRegistration UserInformation) error {
	sqlStr := fmt.Sprintf("INSERT INTO users (email,phone,name,password) VALUES ('%v','%v','%v','%v')", objRegistration.Email, objRegistration.Phone, objRegistration.Name, objRegistration.Password)
	stmt, err := database.DemoDB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func updateUserDetails(objUserDetail UserInformation, id string) error {
	sqlStr := fmt.Sprintf("Update users SET `name` = '%v', `phone` = '%v' where id = '%v'", objUserDetail.Name, objUserDetail.Phone, id)
	stmt, err := database.DemoDB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func getUserList(findEmail string) ([]UserInformation, error) {
	var allUsers []UserInformation
	var whereStr string
	sqlStr := "SELECT id,name,email,phone,password FROM users "
	if findEmail != "" {
		whereStr = " WHERE email = '" + findEmail + "'"
	}
	allRows, err := database.DemoDB.Query(sqlStr + whereStr)
	if err != nil {
		return allUsers, err
	}
	for allRows.Next() {
		var userDetails UserInformation
		var name, email, password, phone sql.NullString
		var id sql.NullInt64
		allRows.Scan(
			&id,
			&name,
			&email,
			&phone,
			&password,
		)
		userDetails.Name = name.String
		userDetails.Email = email.String
		userDetails.Phone = phone.String
		userDetails.ID = strconv.Itoa(int(id.Int64))
		if findEmail != "" {
			userDetails.Password = password.String
		}
		allUsers = append(allUsers, userDetails)
	}
	return allUsers, nil
}
