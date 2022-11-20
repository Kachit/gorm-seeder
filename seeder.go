package gorm_seeder

import (
	"fmt"
	"gorm.io/gorm"
)

//SeederInterface interface
type SeederInterface interface {
	Seed(db *gorm.DB) error
	Clear(db *gorm.DB) error
}

//SeederConfiguration struct
type SeederConfiguration struct {
	Rows int
}

//SeederAbstract struct
type SeederAbstract struct {
	Configuration SeederConfiguration
}

//Delete method
func (sa *SeederAbstract) Delete(db *gorm.DB, table string) error {
	sql := fmt.Sprintf("DELETE FROM %v", table)
	return db.Exec(sql).Error
}

//Truncate method
func (sa *SeederAbstract) Truncate(db *gorm.DB, table string) error {
	sql := fmt.Sprintf("TRUNCATE %v", table)
	return db.Exec(sql).Error
}

//SeedersStack struct
type SeedersStack struct {
	db      *gorm.DB
	seeders []SeederInterface
}

//AddSeeder method
func (ss *SeedersStack) AddSeeder(seeder SeederInterface) *SeedersStack {
	ss.seeders = append(ss.seeders, seeder)
	return ss
}

//Seed method
func (ss *SeedersStack) Seed() error {
	tx := ss.beginTransaction(ss.db)
	for _, seeder := range ss.seeders {
		err := seeder.Seed(tx)
		if err != nil {
			ss.rollbackTransaction(tx)
			return err
		}
	}
	ss.commitTransaction(tx)
	return nil
}

//Clear method
func (ss *SeedersStack) Clear() error {
	tx := ss.beginTransaction(ss.db)
	for _, seeder := range ss.seeders {
		err := seeder.Clear(tx)
		if err != nil {
			ss.rollbackTransaction(tx)
			return err
		}
	}
	ss.commitTransaction(tx)
	return nil
}

//beginTransaction method
func (ss *SeedersStack) beginTransaction(db *gorm.DB) *gorm.DB {
	var tx *gorm.DB
	if db.SkipDefaultTransaction {
		tx = db.Begin()
	} else {
		tx = db
	}
	return tx
}

//commitTransaction method
func (ss *SeedersStack) commitTransaction(db *gorm.DB) *gorm.DB {
	if db.SkipDefaultTransaction {
		db.Commit()
	}
	return db
}

//rollbackTransaction method
func (ss *SeedersStack) rollbackTransaction(db *gorm.DB) *gorm.DB {
	if db.SkipDefaultTransaction {
		db.Rollback()
	}
	return db
}

//NewSeedersStack function
func NewSeedersStack(db *gorm.DB) *SeedersStack {
	return &SeedersStack{db: db}
}

//NewSeederAbstract function
func NewSeederAbstract(cfg SeederConfiguration) SeederAbstract {
	return SeederAbstract{Configuration: cfg}
}
