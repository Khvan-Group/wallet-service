package model

type UserView struct {
	Login      string  `json:"login" database:"login"`
	Email      string  `json:"email" database:"email"`
	FirstName  string  `json:"firstName" database:"first_name"`
	MiddleName *string `json:"middleName" database:"middle_name"`
	LastName   string  `json:"lastName" database:"last_name"`
	Birthdate  string  `json:"birthdate" database:"birthdate"`
	Role       Role    `json:"role" database:"role"`
}

type Role struct {
	Code string `json:"code" database:"code"`
	Name string `json:"name" database:"name"`
}

type JwtUser struct {
	Login string
	Role  string
}
