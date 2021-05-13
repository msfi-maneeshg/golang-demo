package user

//UserInformation :
type UserInformation struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	ProfileImage string `json:"profileImage"`
	Password     string `json:"password"`
}

//LoginDetails :
type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
