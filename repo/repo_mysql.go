package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"ngc8/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoRows        = errors.New("no rows in result set")
	ErrQuery         = errors.New("query execution failed")
	ErrScan          = errors.New("row scanning failed")
	ErrInvalidId     = errors.New("invalid id")
	ErrUserExists    = errors.New("user already exist")
	ErrRowsAffected  = errors.New("unable to get affected row")
	ErrNoAffectedRow = errors.New("rows affected is 0")
	ErrLastInsertId  = errors.New("unable to get last insert id")
	ErrNoUpdate      = errors.New("data already exists")
	ErrBindJSON      = errors.New("unable to bind json")
	ErrParam         = errors.New("error or missing parameter")
	ErrCredential    = errors.New("password or email doesn't match")
)

type ProductRepo interface {
	GetAllProducts() (interface{}, error)
	GetProductById(id uint) (interface{}, error)
	CreateProduct(p model.ProductDB) (interface{}, error)
	UpdateProduct(id int, p model.ProductDB) error
	DeleteProduct(id int) error
}

type UserRepo interface {
	Login(u model.User) (model.User, error)
	Register(u model.User) (model.User, error)
}

func (r *MysqlRepo) Register(u model.User) (model.User, error) {
	var isExist bool

	err := r.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)", u.Email).Scan(&isExist)
	if err != nil {
		return model.User{}, ErrQuery
	}

	if isExist {
		return model.User{}, ErrUserExists
	}

	hashedpwd, _ := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.DefaultCost)
	result, err := r.DB.Exec("INSERT INTO users (name, email, pwd) VALUES (?,?,?)", u.Name, u.Email, hashedpwd)
	if err != nil {
		return model.User{}, ErrQuery
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, ErrLastInsertId
	}

	u.ID = uint(id)
	return u, nil
}

func (r *MysqlRepo) Login(u model.User) (model.User, error) {
	var isExist bool
	var user model.User
	err := r.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)", u.Email).Scan(&isExist)
	if err != nil {
		return model.User{}, ErrQuery
	}

	if !isExist {
		return model.User{}, ErrNoRows
	}

	err = r.DB.QueryRow("SELECT id, name, pwd FROM users WHERE email = ?", u.Email).Scan(&user.ID, &user.Name, &user.Pwd)
	if err != nil {
		return model.User{}, ErrQuery
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(u.Pwd))
	if err != nil {
		return model.User{}, ErrCredential
	}

	return user, nil
}

type MysqlRepo struct {
	DB *sql.DB
}

func (r *MysqlRepo) IsIDExist(id int) (bool, error) {
	var isExist bool
	err := r.DB.QueryRow("SELECT EXISTS (SELECT 1 from products WHERE id = ?)", id).Scan(&isExist)
	if err != nil {
		log.Println("error querying", err)
		return false, err
	}

	return isExist, nil
}

func (r *MysqlRepo) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	query := "SELECT id, name, description, img, price, store_name FROM products JOIN stores ON products.store_id = stores.store_id"
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Println("error query", err)
		return nil, ErrQuery
	}

	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Desc, &p.Img, &p.Price, &p.Store)
		if err != nil {
			log.Println("error scan row")
			return nil, ErrScan
		}

		products = append(products, p)
	}

	if len(products) == 0 {
		log.Println("empty table")
		return nil, ErrNoRows
	}

	return products, nil
}

func (r *MysqlRepo) GetProductById(id int) (model.Product, error) {
	var p model.Product
	isExist, err := r.IsIDExist(id)
	if err != nil {
		log.Println("error query", err)
		return model.Product{}, ErrQuery
	}

	if !isExist {
		log.Println("id not found")
		return model.Product{}, ErrNoRows
	}
	fmt.Println("here")
	query := "SELECT name, description, img, price, store_name FROM products JOIN stores ON products.store_id = stores.store_id WHERE id = ?"
	err = r.DB.QueryRow(query, id).Scan(&p.Name, &p.Desc, &p.Img, &p.Price, &p.Store)
	if err != nil {
		log.Println("error query")
		return model.Product{}, ErrQuery
	}

	p.Id = int(id)
	return p, nil
}

func (r *MysqlRepo) CreateProduct(p model.Product) (model.Product, error) {

	query := "INSERT INTO products (name, description, img, price, store_id) VALUES (?,?,?,?,(SELECT store_id FROM stores WHERE store_name = ?))"

	fmt.Println(p)
	result, err := r.DB.Exec(query, p.Name, p.Desc, p.Img, p.Price, p.Store)
	if err != nil {
		log.Println("error query", err)
		return model.Product{}, ErrQuery
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println("error getting last inserted id")
		return model.Product{}, ErrLastInsertId
	}

	p.Id = int(lastId)
	return p, nil
}

func (r *MysqlRepo) UpdateProduct(id int, p model.Product) error {
	isExist, err := r.IsIDExist(id)
	if err != nil {
		log.Println("error query", err)
		return ErrQuery
	}

	if !isExist {
		log.Println("id not found")
		return ErrNoRows
	}

	query := `
	UPDATE products
	SET name =?,
	description = ?, 
	img = ?,
	price = ?,
	store_id = (SELECT store_id FROM stores WHERE store_name = ?)
	WHERE id = ?
	`

	result, err := r.DB.Exec(query, p.Name, p.Desc, p.Img, p.Price, p.Store, id)
	if err != nil {
		log.Println("error query", err)
		return ErrQuery
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Println("error getting num of rows affected")
		return ErrRowsAffected
	}

	if affectedRow == 0 {
		log.Println("")
		return ErrNoUpdate
	}
	return nil
}

func (r *MysqlRepo) DeleteProduct(id int) error {
	isExist, err := r.IsIDExist(id)
	if err != nil {
		log.Println("error query", err)
		return ErrQuery
	}

	if !isExist {
		log.Println("id not found")
		return ErrNoRows
	}

	query := "DELETE FROM products WHERE id = ? "
	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Println("error query", err)
		return ErrQuery
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Println("error getting num of rows affected")
		return ErrRowsAffected
	}

	if affectedRow == 0 {
		log.Println("")
		return ErrNoUpdate
	}

	return nil
}
