package post

import (
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
