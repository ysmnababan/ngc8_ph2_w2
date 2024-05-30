deps:
	go get -u github.com/julienschmidt/httprouter
	go get -u golang.org/x/crypto/bcrypt
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/joho/godotenv
	go get -u github.com/golang-jwt/jwt/v4
	go get -u github.com/gin-gonic/gin

.PHONY: all
all: deps
