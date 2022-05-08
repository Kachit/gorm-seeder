package gorm_seeder

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type StubUser struct {
	Id        uint64    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *StubUser) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now().UTC()
	return
}

// TableName overrides
func (StubUser) TableName() string {
	return "users"
}

type StubUsersSeederV1 struct {
	SeederAbstract
}

func NewUsersSeederV1(cfg SeederConfiguration) *StubUsersSeederV1 {
	return &StubUsersSeederV1{NewSeederAbstract(cfg)}
}

func (s *StubUsersSeederV1) Seed(db *gorm.DB) error {
	var users []StubUser
	for i := 0; i < s.Configuration.Rows; i++ {
		indexStr := fmt.Sprint(i)
		user := StubUser{
			Name:     "Name LastName" + indexStr,
			Email:    "foo" + indexStr + "@bar.gov",
			Password: "password-hash-string" + indexStr,
		}
		users = append(users, user)
	}
	return db.CreateInBatches(users, s.Configuration.Rows).Error
}

func (s *StubUsersSeederV1) Clear(db *gorm.DB) error {
	entity := StubUser{}
	return s.SeederAbstract.Delete(db, entity.TableName())
}

type StubUsersSeederV2 struct {
	SeederAbstract
}

func NewUsersSeederV2(cfg SeederConfiguration) *StubUsersSeederV2 {
	return &StubUsersSeederV2{NewSeederAbstract(cfg)}
}

func (s *StubUsersSeederV2) Seed(db *gorm.DB) error {
	for i := 0; i < s.Configuration.Rows; i++ {
		user := StubUser{
			Name:     "Name LastName",
			Email:    "foo@bar.gov",
			Password: "password-hash-string",
		}
		err := db.Create(&user).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StubUsersSeederV2) Clear(db *gorm.DB) error {
	entity := StubUser{}
	return s.SeederAbstract.Truncate(db, entity.TableName())
}
