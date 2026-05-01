module ecommitra-backend

go 1.19

require (
	github.com/joho/godotenv v1.5.1
	github.com/golang-jwt/jwt/v4 v4.5.0
	golang.org/x/crypto v0.14.0
	golang.org/x/text v0.13.0
	golang.org/x/sys v0.13.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.4
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.14.0
	golang.org/x/text => golang.org/x/text v0.13.0
	golang.org/x/sys => golang.org/x/sys v0.13.0
)
