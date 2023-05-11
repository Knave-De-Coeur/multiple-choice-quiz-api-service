# User-API-Server
Golang Microservice/REST API

## Overview 

This project is a dedicated template project for a simple user CRUD and authorization.

- Sign up 
- Log in 
- Exit

**REST enpoints include;**

- **POST - /api/v1/login** - Generates JWT token used to authorize REST APIs
- **POST - /api/v1/logout** - Destroys Token 

**USER CRUD:** 
- **POST - /api/v1/user** - adds new user 
- **PUT - /api/v1/user/:uID** - modifies user data based on the json payload sent 
- **DELETE - /api/v1/user/:uID** - Either soft deletes or completely removes row from db
- **GET - /api/v1/user/:uID** - Gets specific user data (if authorized) 
- **GET - /api/v1/users** - Gets list of user data (just email and username)

--- 

### Set up 

To build project, start api and run migrations : 

``` 
cd multiple-choice-quiz-api-service

docker-compose up
```
---

### Running application

- Set up and run the project with all dependencies and local env:

```
docker-compose up -d
```

- Generate the executable locally:
```
go build -o $GOPATH/bin $GOPATH/src/github.com/knave-de-coeur/user-api/cmd/api/main.go
```

## Technical Description

Docker was used to create the mysql database and build an image of the api inside another container, attempting to simulate a real-life server environment on my local machine.

I'd just like to thank the following repos for their incredible libraries that helped make this personal project come to life.
Big shout out to:

- [Viper](https://github.com/spf13/viper) was used to manage configurations including fallbacks.
- [migrate](https://github.com/golang-migrate/migrate) for it also creates the schema and handles migrations.
- [gorm](https://github.com/go-gorm/gorm) for the sick ORM that made querying a breeze
- [gin](https://github.com/gin-gonic/gin) for easy set up with the api


