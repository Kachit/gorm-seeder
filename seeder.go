package gorm_seeder

import (
	"gorm.io/gorm"
)

type SeederInterface interface {
	Seed(db *gorm.DB) error
	Clear(db *gorm.DB) error
}

type SeederConfiguration struct {
	Rows int
}

type SeederAbstract struct {
	Configuration SeederConfiguration
}

type SeedersStack struct {
	Seeders []SeederInterface
}

func (ss *SeedersStack) AddSeeder(seeder SeederInterface) *SeedersStack {
	ss.Seeders = append(ss.Seeders, seeder)
	return ss
}

func (ss *SeedersStack) Seed(db *gorm.DB) error {
	tx := db.Session(&gorm.Session{})
	for _, seeder := range ss.Seeders {
		err := seeder.Seed(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss *SeedersStack) Clear(db *gorm.DB) error {
	tx := db.Session(&gorm.Session{})
	for _, seeder := range ss.Seeders {
		err := seeder.Clear(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSeederAbstract(cfg SeederConfiguration) SeederAbstract {
	return SeederAbstract{Configuration: cfg}
}
