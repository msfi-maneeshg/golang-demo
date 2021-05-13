package user

import (
	"database/sql"
	"fmt"
	"golang-demo/api/data"
	"strconv"
)

//---------data functions
func isEmailExist(email string) (bool, error) {
	var emailID string
	sqlStr := "SELECT email FROM users WHERE email = ?"
	err := data.DemoDB.QueryRow(sqlStr, email).Scan(&emailID)
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

	err := data.DemoDB.QueryRow(sqlStr, id).Scan(&emailID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func insertNewUser(objRegistration UserInformation) error {
	noProfileImage := "no-profile.jpg"
	sqlStr := fmt.Sprintf("INSERT INTO users (email,phone,name,password,profile_image) VALUES ('%v','%v','%v','%v','%v')", objRegistration.Email, objRegistration.Phone, objRegistration.Name, objRegistration.Password, noProfileImage)
	stmt, err := data.DemoDB.Prepare(sqlStr)
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
	var profileImageStr string
	if objUserDetail.ProfileImage != "" {
		profileImageStr = fmt.Sprintf(", `profile_image` = '%v'", objUserDetail.ProfileImage)
	}
	sqlStr := fmt.Sprintf("Update users SET `name` = '%v', `phone` = '%v'"+profileImageStr+" where id = '%v'", objUserDetail.Name, objUserDetail.Phone, id)

	stmt, err := data.DemoDB.Prepare(sqlStr)
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

func updateUserPassword(newPassword, id string) error {
	sqlStr := fmt.Sprintf("Update users SET `password` = '%v'where id = '%v'", newPassword, id)
	stmt, err := data.DemoDB.Prepare(sqlStr)
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
	sqlStr := "SELECT id,name,email,phone,password,profile_image FROM users "
	if findEmail != "" {
		whereStr = " WHERE email = '" + findEmail + "'"
	}
	allRows, err := data.DemoDB.Query(sqlStr + whereStr)
	if err != nil {
		return allUsers, err
	}
	for allRows.Next() {
		var userDetails UserInformation
		var name, email, password, phone, profileImage sql.NullString
		var id sql.NullInt64
		allRows.Scan(
			&id,
			&name,
			&email,
			&phone,
			&password,
			&profileImage,
		)
		userDetails.Name = name.String
		userDetails.Email = email.String
		userDetails.Phone = phone.String
		userDetails.ProfileImage = profileImage.String
		userDetails.ID = strconv.Itoa(int(id.Int64))
		if findEmail != "" {
			userDetails.Password = password.String
		}
		allUsers = append(allUsers, userDetails)
	}
	return allUsers, nil
}

func GetUserInfo(id string) (objUserInfo UserInformation, err error) {
	sqlStr := "SELECT name,email,phone,password,profile_image FROM users WHERE id = ?"
	var name, email, password, phone, profileImage sql.NullString
	err = data.DemoDB.QueryRow(sqlStr, id).Scan(
		&name,
		&email,
		&phone,
		&password,
		&profileImage,
	)
	if err != nil && err != sql.ErrNoRows {
		return objUserInfo, err
	}
	if err != sql.ErrNoRows {
		objUserInfo.Name = name.String
		objUserInfo.Email = email.String
		objUserInfo.Phone = phone.String
		objUserInfo.ProfileImage = profileImage.String
		objUserInfo.ID = id
	}
	return objUserInfo, nil
}
