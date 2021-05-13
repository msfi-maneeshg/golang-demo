package user

import (
	"encoding/base64"
	"encoding/json"
	"golang-demo/api/common"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

}

//UserList :
func UserList(w http.ResponseWriter, r *http.Request) {
	userList, err := getUserList("")
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while getting user list")
		return
	}
	common.APIResponse(w, http.StatusOK, userList)
}

//UserDetail :
func UserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	userProfile, err := GetUserInfo(id)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while getting user list")
		return
	}
	if userProfile == (UserInformation{}) {
		common.APIResponse(w, http.StatusNotFound, "User profile not found")
		return
	}

	common.APIResponse(w, http.StatusOK, userProfile)
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

}

//UpdateDetail :
func UpdateDetail(w http.ResponseWriter, r *http.Request) {
	var objUserInformation UserInformation
	var err error
	vars := mux.Vars(r)
	var id = vars["id"]

	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objUserInformation)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	//-------validate data

	if objUserInformation.Name == "" {
		common.APIResponse(w, http.StatusBadRequest, "Name can not be empty")
		return
	}

	if objUserInformation.Phone == "" {
		common.APIResponse(w, http.StatusBadRequest, "Phone number can not be empty")
		return
	}
	_, err = strconv.Atoi(objUserInformation.Phone)
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

	//-------
	if objUserInformation.ProfileImage != "" {
		fileName := time.Now().Format("200601021504050700") + "-" + objUserInformation.Name + ".jpg"
		dec, err := base64.StdEncoding.DecodeString(objUserInformation.ProfileImage)
		if err != nil {
			common.APIResponse(w, http.StatusBadRequest, "Invalid data of image."+err.Error())
			return
		}

		f, err := os.Create("images/profile/" + fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}
		objUserInformation.ProfileImage = fileName
	}

	err = updateUserDetails(objUserInformation, id)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while updating user info. Error:"+err.Error())
		return
	}

	common.APIResponse(w, http.StatusOK, "User detail has been updated.")

}

//UpdatePassword :
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	var id = vars["id"]

	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if password == "" || confirmPassword == "" {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank.")
		return
	}

	if password != confirmPassword {
		common.APIResponse(w, http.StatusBadRequest, "Password is not matching.")
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

	err = updateUserPassword(password, id)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while updating user password. Error:"+err.Error())
		return
	}

	common.APIResponse(w, http.StatusOK, "Password has been changed.")
}
