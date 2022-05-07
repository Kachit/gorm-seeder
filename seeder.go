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
	tx := ss.beginTransaction(db)
	for _, seeder := range ss.Seeders {
		err := seeder.Seed(tx)
		if err != nil {
			ss.rollbackTransaction(db)
			return err
		}
	}
	ss.commitTransaction(db)
	return nil
}

func (ss *SeedersStack) Clear() error {
	db := ss.getDb()
	tx := ss.beginTransaction(db)
	for _, seeder := range ss.Seeders {
		err := seeder.Clear(tx)
		if err != nil {
			ss.rollbackTransaction(db)
			return err
		}
	}
	ss.commitTransaction(db)
	return nil
}

func (ss *SeedersStack) getDb() *gorm.DB {
	if ss.db != nil {
		return ss.db
	}
	return ss.db.Session(&gorm.Session{SkipDefaultTransaction: true})
}

func (ss *SeedersStack) beginTransaction(db *gorm.DB) *gorm.DB {
	if db.SkipDefaultTransaction == true {
		db.Begin()
	}
	return db
}

func (ss *SeedersStack) commitTransaction(db *gorm.DB) *gorm.DB {
	if db.SkipDefaultTransaction == true {
		db.Commit()
	}
	return db
}

func (ss *SeedersStack) rollbackTransaction(db *gorm.DB) *gorm.DB {
	if db.SkipDefaultTransaction == true {
		db.Rollback()
	}
	return db
}

func NewSeedersStack(db *gorm.DB) *SeedersStack {
	return &SeedersStack{db: db}
}

func NewSeederAbstract(cfg SeederConfiguration) SeederAbstract {
	return SeederAbstract{Configuration: cfg}
}
