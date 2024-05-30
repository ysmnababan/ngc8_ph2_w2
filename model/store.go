package model

type Store struct {
	StoreID    uint        `gorm:"primaryKey;autoIncrement"`
	StoreName  string      `gorm:"size:255"`
	StorePwd   string      `gorm:"size:255"`
	StoreEmail string      `gorm:"size:255;unique"`
	StoreType  string      `gorm:"size:255"`
	Products   []ProductDB `gorm:"foreignKey:StoreID"`
}

type ProductDB struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" gorm:"size:255"`
	Desc    string `json:"desc" gorm:"size:255"`
	Img     string `json:"img" gorm:"size:255"`
	Price   int    `json:"price"`
	StoreID uint   `json:"store_id"`
}
