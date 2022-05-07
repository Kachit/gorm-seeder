package gorm_seeder

import (
	"gorm.io/gorm"
	"time"
)

type StubUsersSeeder struct {
	SeederAbstract
}

func NewUsersSeeder(cfg SeederConfiguration) StubUsersSeeder {
	return StubUsersSeeder{NewSeederAbstract(cfg)}
}

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

func (s *StubUsersSeeder) Seed1(db *gorm.DB) error {
	var users []StubUser
	for i := 0; i < s.Configuration.Rows; i++ {
		user := StubUser{
			Name:     "Name LastName",
			Email:    "foo@bar.gov",
			Password: "password-hash-string",
		}
		users = append(users, user)
	}
	return db.CreateInBatches(users, s.Configuration.Rows).Error
}

func (s *StubUsersSeeder) Seed(db *gorm.DB) error {
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

func (s *StubUsersSeeder) Clear(db *gorm.DB) error {
	entity := StubUser{}
	return s.SeederAbstract.Delete(db, entity.TableName())
}
