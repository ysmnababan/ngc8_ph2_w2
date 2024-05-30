package repo

import (
	"fmt"
	"ngc8/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PostgreRepo struct {
	DB *gorm.DB
}

func (r *PostgreRepo) GetAllProducts() (interface{}, error) {
	var products []model.ProductDB

	r.DB.Find(&products)

	if len(products) == 0 {
		return nil, ErrNoRows
	}

	return products, nil
}

func (r *PostgreRepo) GetProductById(id uint) (interface{}, error) {
	var p model.ProductDB
	result := r.DB.First(&p, id)
	fmt.Println(result.Error)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, ErrNoRows
	}

	return p, nil
}

func (r *PostgreRepo) CreateProduct(p model.ProductDB) (interface{}, error) {

	result := r.DB.Create(&p)
	if result.Error != nil {
		return nil, ErrQuery
	}

	return p, nil
}

func (r *PostgreRepo) UpdateProduct(id int, p model.ProductDB) error {
	var product model.ProductDB

	result := r.DB.First(&product, id)
	fmt.Println(result.Error)
	if result.Error == gorm.ErrRecordNotFound {
		return ErrNoRows
	}

	p.ID = product.ID
	r.DB.Save(&p)

	return nil
}

func (r *PostgreRepo) DeleteProduct(id int) error {
	var p model.ProductDB
	p.ID = uint(id)

	result := r.DB.First(&p)
	if result.Error == gorm.ErrRecordNotFound {
		return ErrNoRows
	}
	r.DB.Delete(&p)

	return nil
}

func (r *PostgreRepo) Register(u model.User) (model.User, error) {
	var exist model.User
	r.DB.Where("email= ?", u.Email).Find(&exist)

	if exist.Email == u.Email {
		return model.User{}, ErrUserExists
	}

	hashedpwd, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("error generating hashed")
	}
	u.Pwd = string(hashedpwd)
	result := r.DB.Create(&u)
	if result.Error != nil {
		return model.User{}, ErrQuery
	}
	return u, nil
}

func (r *PostgreRepo) Login(u model.User) (model.User, error) {
	var user model.User
	r.DB.Where("email = ?", u.Email).Find(&user)

	fmt.Println(u, user)
	// user not found
	if user.Email != u.Email {
		return model.User{}, ErrCredential
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(u.Pwd))
	if err != nil {
		return model.User{}, ErrCredential
	}
	return user, nil
}
