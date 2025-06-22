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
		- CreatedAt: 作成日時
		- UpdatedAt: 更新日時
		--------------------------------
		内容
	*/
	return fmt.Sprintf("# %s\n\n- CreatedAt: %s\n- UpdatedAt: %s\n\n%s\n\n", memo.Title, memo.CreatedAt, memo.UpdatedAt, memo.Content)
}

func generateJsonContent(memo *model.Memo) string {
	/*
		Sample
		{
			"title": "タイトル",
			"content": "内容",
			"createdAt": "作成日時",
			"updatedAt": "更新日時"
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
	CreatedAt: 作成日時
	UpdatedAt: 更新日時
	--------------------------------
	内容
	*/
	return fmt.Sprintf("Title: %s\nCreatedAt: %s\nUpdatedAt: %s\n--------------------------------\n%s", memo.Title, memo.CreatedAt, memo.UpdatedAt, memo.Content)
}
