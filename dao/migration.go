package dao

import (
	"aresumei/model"
	"fmt"
)

func Migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&model.User{}, &model.Resume{})
	if err != nil {
		fmt.Println(err)
	}
	return
}
