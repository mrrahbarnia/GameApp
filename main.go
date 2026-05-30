package main

import (
	"github.com/mrrahbarnia/GameApp/repository/postgresql"
)

func main() {

}

func test() {
	repo := postgresql.New()

	// repo.Register(
	// 	entity.User{
	// 		ID:          0,
	// 		Name:        "test",
	// 		PhoneNumber: "09131111111",
	// 	},
	// )
	exist, err := repo.IsPhoneNumberExist(
		"09131111113",
	)
	if err != nil {
		println(err)
	}
	println(exist)
}
