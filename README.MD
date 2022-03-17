# Simple Bank API :bank:

Simple banking project. CRUD transaction between two users with different accounts.

## How to start a Go project?

1. Create container folder
2. Create `main.go` file
3. Initialize modules: `go mod init github.com/{username}/{projectdir}`
4. Install packages:
   - go get github.com/badoux/checkmail
   - go get github.com/jinzhu/gorm
   - go get golang.org/x/crypto/bcrypt
   - go get github.com/dgrijalva/jwt-go
   - go get github.com/gorilla/mux
   - go get github.com/jinzhu/gorm/dialects/mysql" //If using mysql
   - go get github.com/jinzhu/gorm/dialects/postgres //If using postgres
   - go get github.com/joho/godotenv
   - go get gopkg.in/go-playground/assert.v1

## Commands

- Run script (project in this case): `go run main.go`. main.go can be whatever file you want.
- Add module requirements and sums: `go mod tidy`.
