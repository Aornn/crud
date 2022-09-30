package pkg

import "encoding/json"

type (
	Friend struct {
		ID   int    `form:"id" json:"id" xml:"id"  binding:"required"`
		Name string `form:"name" json:"name" xml:"name"  binding:"required"`
	}
	User struct {
		ID        *string      `form:"id" json:"id" xml:"id"`
		Password  *string      `form:"password" json:"password" xml:"password"`
		IsActive  *bool        `form:"isActive" json:"isActive" xml:"isActive"`
		Balance   *string      `form:"balance" json:"balance" xml:"balance"`
		Age       *json.Number `form:"age" json:"age" xml:"age"`
		Gender    *string      `form:"gender" json:"gender" xml:"gender"`
		Company   *string      `form:"company" json:"company" xml:"company"`
		Email     *string      `form:"email" json:"email" xml:"email"`
		Phone     *string      `form:"phone" json:"phone" xml:"phone"`
		Address   *string      `form:"address" json:"address" xml:"address"`
		About     *string      `form:"about" json:"about" xml:"about"`
		Registred *string      `form:"registered" json:"registered" xml:"registered"`
		Latitude  *float64     `form:"latitude" json:"latitude" xml:"latitude"`
		Longitude *float64     `form:"longitude" json:"longitude" xml:"longitude"`
		Tags      *[]string    `form:"tags" json:"tags" xml:"tags"`
		Friends   *[]Friend    `form:"friends" json:"friends" xml:"friends"`
		Data      *string      `form:"data" json:"data" xml:"data"`
	}

	UserLogin struct {
		ID       string `form:"id" json:"id" xml:"id"  binding:"required"`
		Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	}
)
