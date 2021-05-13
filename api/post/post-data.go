package post

import (
	"database/sql"
	"fmt"
	"golang-demo/api/data"
	"strings"
	"time"
)

func InsertNewPost(objPostInfo PostInfo) error {
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := fmt.Sprintf("INSERT INTO post (user_id,content,images,dateadded,status) VALUES ('%v','%v','%v','%v','%v')", objPostInfo.UserID, strings.ReplaceAll(objPostInfo.Content, `'`, `\'`), objPostInfo.Image, currentDateTime, 1)
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

func GetAllPostDetails() ([]PostInfo, error) {
	var allPostInfo []PostInfo
	var whereStr string
	sqlStr := `SELECT p.id, p.content, p.images, p.dateadded, u.id, u.name, u.profile_image 
		FROM post AS p 
		LEFT JOIN users AS u ON  p.user_id = u.id 
		order by p.dateadded DESC`

	allRows, err := data.DemoDB.Query(sqlStr + whereStr)
	if err != nil {
		return allPostInfo, err
	}
	for allRows.Next() {
		var objPostInfo PostInfo
		var content, images, dateadded, name, profileImage sql.NullString
		var id, userID sql.NullInt64
		allRows.Scan(
			&id,
			&content,
			&images,
			&dateadded,
			&userID,
			&name,
			&profileImage,
		)
		objPostInfo.ID = int(id.Int64)
		objPostInfo.UserID = int(userID.Int64)
		objPostInfo.Content = content.String
		objPostInfo.Image = images.String
		objPostInfo.Name = name.String
		objPostInfo.ProfileImage = profileImage.String
		objPostInfo.Datetime = dateadded.String

		allPostInfo = append(allPostInfo, objPostInfo)
	}
	return allPostInfo, nil
}
