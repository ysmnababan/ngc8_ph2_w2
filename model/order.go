package model

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Img   string `json:"img"`
	Price int    `json:"price"`
	Store string `json:"store"`
}
