# Simple Gorm seeder
[![Build Status](https://app.travis-ci.com/Kachit/gorm-seeder.svg?branch=master)](https://app.travis-ci.com/Kachit/gorm-seeder)
[![Codecov](https://codecov.io/gh/Kachit/gorm-seeder/branch/master/graph/badge.svg)](https://codecov.io/gh/Kachit/gorm-seeder)
[![Go Report Card](https://goreportcard.com/badge/github.com/kachit/gorm-seeder)](https://goreportcard.com/report/github.com/kachit/gorm-seeder)
[![Release](https://img.shields.io/github/v/release/Kachit/gorm-seeder.svg)](https://github.com/Kachit/gorm-seeder/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/kachit/gorm-seeder/blob/master/LICENSE)

## Description
Simple Gorm seeder package

## Download
```shell
go get -u github.com/kachit/gorm-seeder
```

## Usage
```go
package main

import (
    "fmt"
    "github.com/kachit/gorm-seeder"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "time"
)

//Init model
type User struct {
    Id          uint64 `json:"id" gorm:"primaryKey"`
    Name        string `json:"name"`
    Email        string `json:"email"`
    Password       string `json:"password"`
    CreatedAt  time.Time      `json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.CreatedAt = time.Now().UTC()
    return
}

// TableName overrides
func (User) TableName() string {
    return "users"
}

//Create new seeder instance of gorm_seeder.SeederInterface
type UsersSeeder struct {
    gorm_seeder.SeederAbstract
}

func NewUsersSeeder(cfg gorm_seeder.SeederConfiguration) UsersSeeder {
    return UsersSeeder{gorm_seeder.NewSeederAbstract(cfg)}
}

//Implement Seed method
func (s *UsersSeeder) Seed(db *gorm.DB) error {
    var users []User
    for i := 0; i < s.Configuration.Rows; i++ {
        indexStr := fmt.Sprint(i)
        user := User{
            Name: "Name LastName" + indexStr,
            Email: "foo" + indexStr + "@bar.gov",
            Password: "password-hash-string" + indexStr,
        }
        users = append(users, user)
    }
    return db.CreateInBatches(users, s.Configuration.Rows).Error
}

//Implement Clear method
func (s *UsersSeeder) Clear(db *gorm.DB) error {
    entity := User{}
    return s.SeederAbstract.Delete(db, entity.TableName())
}

func main(){
    //Init DB connection
    db, _ := gorm.Open(postgres.New(postgres.Config{
        DSN: "DSN connection string",
    }))

    //Build seeders stack
    usersSeeder := NewUsersSeeder(gorm_seeder.SeederConfiguration{Rows: 10})
    seedersStack := gorm_seeder.NewSeedersStack(db)
    seedersStack.AddSeeder(&usersSeeder)

    //Apply seed
    err := seedersStack.Seed()
    fmt.Println(err)

    //Apply clear
    err = seedersStack.Clear()
    fmt.Println(err)
}
```
