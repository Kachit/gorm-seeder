package gorm_seeder

import (
	"fmt"
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

func (sa *SeederAbstract) Delete(db *gorm.DB, table string) error {
	sql := fmt.Sprintf("DELETE FROM %v", table)
	return db.Exec(sql).Error
}

func (sa *SeederAbstract) Truncate(db *gorm.DB, table string) error {
	sql := fmt.Sprintf("TRUNCATE %v", table)
	return db.Exec(sql).Error
}

type SeedersStack struct {
	db      *gorm.DB
	seeders []SeederInterface
}

func (ss *SeedersStack) AddSeeder(seeder SeederInterface) *SeedersStack {
	ss.seeders = append(ss.seeders, seeder)
	return ss
}

func (ss *SeedersStack) Seed() error {
	db := ss.getDb()
	tx := ss.beginTransaction(db)
	for _, seeder := range ss.seeders {
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
	for _, seeder := range ss.seeders {
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
	return ss.db
}

func (ss *SeedersStack) beginTransaction(db *gorm.DB) *gorm.DB {
	if db.SkipDefaultTransaction == true {
		tx := db.Begin()
		db = tx
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
