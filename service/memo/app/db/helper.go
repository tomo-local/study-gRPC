package db

import (
	"encoding/json"
	"fmt"
	"memo/db/model"
)

func generateMarkdownContent(memo *model.Memo) string {
	/*
		Sample
		# タイトル
		- ModifiedAt: 更新日時
		--------------------------------
		内容
	*/
	return fmt.Sprintf("## %s\n\n- ModifiedAt: %s\n\n--------------------------------\n%s\n\n", memo.Title, memo.ModifiedAt, memo.Content)
}

func generateJsonContent(memo *model.Memo) string {
	/*
		Sample
		{
			"title": "タイトル",
			"content": "内容",
			"modified_at": "更新日時"
		}
	*/
	json, err := json.Marshal(memo)
	if err != nil {
		return "{}"
	}
	return string(json)
}

func generateTextContent(memo *model.Memo) string {
	/*　Sample
	Title: タイトル
	ModifiedAt: 更新日時
	--------------------------------
	内容
	*/
	return fmt.Sprintf("Title: %s\nModifiedAt: %s\n--------------------------------\n%s", memo.Title, memo.ModifiedAt, memo.Content)
}
