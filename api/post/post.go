package post

import (
	"golang-demo/api/common"
	"golang-demo/api/user"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func AddPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var userID = vars["id"]

	var postImageName string
	postContent := r.FormValue("postContent")
	file, handler, _ := r.FormFile("postImage")

	if handler == nil && postContent == "" {
		common.APIResponse(w, http.StatusBadRequest, "Post data can not be empty.")
		return
	}

	//-------check userID
	userInfo, err := user.GetUserInfo(userID)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while checking userID."+err.Error())
		return
	}
	if userInfo == (user.UserInformation{}) {
		common.APIResponse(w, http.StatusUnauthorized, "Invalid userID.")
		return
	}

	if handler != nil && handler.Filename != "" {
		defer file.Close()
		postImageName = time.Now().Format("200601021504050700") + "-" + handler.Filename
		dst, err := os.Create("images/post/" + postImageName)
		defer dst.Close()
		if err != nil {
			common.APIResponse(w, http.StatusInternalServerError, "Error while creating new image file."+err.Error())
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			common.APIResponse(w, http.StatusInternalServerError, "Error while saving image file."+err.Error())
			return
		}
	}

	var objNewPost PostInfo
	objNewPost.Content = postContent
	objNewPost.Image = postImageName
	objNewPost.UserID, _ = strconv.Atoi(userID)
	err = InsertNewPost(objNewPost)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while inserting new post."+err.Error())
		return
	}
	common.APIResponse(w, http.StatusOK, "Post added successfully.")
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	allPosts, err := GetAllPostDetails()
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Error while getting user post"+err.Error())
		return
	}
	if len(allPosts) == 0 {
		common.APIResponse(w, http.StatusNoContent, "There is no post!")
		return
	}
	common.APIResponse(w, http.StatusOK, allPosts)
}
