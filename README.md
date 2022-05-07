# Simple Gorm seeder

## Description
Gorm seeder package

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

type UsersSeeder struct {
	gorm_seeder.SeederAbstract
}

func NewUsersSeeder(cfg gorm_seeder.SeederConfiguration) UsersSeeder {
	return UsersSeeder{gorm_seeder.NewSeederAbstract(cfg)}
}

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

func (s *UsersSeeder) Seed(db *gorm.DB) error {
    var users []User
    for i := 0; i < s.Configuration.Rows; i++ {
    
        box := User{
            Name:        "Name LastName",
            Email:        "foo@bar.gov",
            Password:        "password-hash-string",
        }
        users = append(users, box)
    }
    return db.CreateInBatches(users, s.Configuration.Rows).Error
}

func (s *UsersSeeder) Clear(db *gorm.DB) error {
    entity := User{}
    sql := fmt.Sprintf("DELETE FROM %v", entity.TableName())
    return db.Exec(sql).Error
}

func main(){
    db, _ := gorm.Open(postgres.New(postgres.Config{
        DSN: "DSN connection string",
    }))

    usersSeeder := NewUsersSeeder(gorm_seeder.SeederConfiguration{Rows: 10})
    seedersStack := gorm_seeder.NewSeedersStack(db)
    seedersStack.AddSeeder(&usersSeeder)

    err := seedersStack.Seed()
    fmt.Println(err)
}
```
