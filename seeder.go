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
	db      *gorm.DB
	Seeders []SeederInterface
}

func (ss *SeedersStack) AddSeeder(seeder SeederInterface) *SeedersStack {
	ss.Seeders = append(ss.Seeders, seeder)
	return ss
}

func (ss *SeedersStack) Seed() error {
	db := ss.getDb()
	tx := db.Begin()
	for _, seeder := range ss.Seeders {
		err := seeder.Seed(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (ss *SeedersStack) Clear() error {
	db := ss.getDb()
	tx := db.Begin()
	for _, seeder := range ss.Seeders {
		err := seeder.Clear(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (ss *SeedersStack) getDb() *gorm.DB {
	if ss.db != nil {
		return ss.db
	}
	return ss.db.Session(&gorm.Session{SkipDefaultTransaction: true})
}

func NewSeedersStack(db *gorm.DB) *SeedersStack {
	return &SeedersStack{db: db}
}

func NewSeederAbstract(cfg SeederConfiguration) SeederAbstract {
	return SeederAbstract{Configuration: cfg}
}
