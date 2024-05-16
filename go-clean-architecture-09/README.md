# clean-architecture-with-go
## 🗄 About this repository
* This is a sample API built by Go(Echo) and SQLBoiler according to Clean Architecture.

## 👟 How to run
### install
#### this repository
```shell
git clone git@github.com:yagikota/clean-architecture-with-go.git

cd clean_architecture-with-go
```
#### [SQLBoiler](https://github.com/volatiletech/sqlboiler#getting-started)
```shell
$ go install github.com/volatiletech/sqlboiler/v4@latest
# install driver for MySQL
$ go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
```

#### [golang-migrate](https://formulae.brew.sh/formula/golang-migrate)
```shell
brew install golang-migrate
```

### set environment variables
```shell
echo "PORT=8080

MYSQL_USER=root
MYSQL_ROOT_PASSWORD=
MYSQL_ALLOW_EMPTY_PASSWORD=true
MYSQL_PORT=3306
MYSQL_ADDR=mysql:3306
MYSQL_DATABASE=sample" > .env
```

### initialize local DB
```shell
make run-db
make migrate 
make seed
```
`make migrate` command may fail, but try some time.

### build API
```shell
make run-go
```
and access`localhost:8080/health_check`.
You will get the following response.
```json
{
    "message": "Hello! you've requested: /health_check"
}
```
Then, you have fully prepared for running API.

## 📄 API Document
* Please copy and paste [this file](https://github.com/yagikota/clean_architecture_with_go/blob/main/api_doc.yml) on https://editor.swagger.io/ .

🐶 I hope this repository helps you studying Clean Architecture with Go.

## Docuemtn

https://zenn.dev/88888888_kota/articles/1a5caac8d743b8
git clone https://github.com/yagikota/clean-architecture-with-go.git go-clean-architecture-09