package gorm_seeder

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
)

func Test_SeederAbstract_Delete(t *testing.T) {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM foo`)).WillReturnResult(sqlmock.NewResult(0, 1))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Delete(db, "foo")
	assert.NoError(t, err)
}

func Test_SeederAbstract_Truncate(t *testing.T) {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectExec(regexp.QuoteMeta(
		`TRUNCATE foo`)).WillReturnResult(sqlmock.NewResult(0, 1))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Truncate(db, "foo")
	assert.NoError(t, err)
}

func Test_SeedersStack_AddSeeder(t *testing.T) {
	ss := NewSeedersStack(nil)
	assert.Len(t, ss.seeders, 0)
	ss.AddSeeder(NewUsersSeederV1(SeederConfiguration{}))
	assert.Len(t, ss.seeders, 1)
}

func Test_SeedersStack_BeginTransaction(t *testing.T) {
	db, mock := getDatabaseMock()
	db.SkipDefaultTransaction = true
	mock.MatchExpectationsInOrder(false)
	ss := NewSeedersStack(db)
	ss.beginTransaction(db)
}

func Test_SeedersStack_CommitTransaction(t *testing.T) {
	db, mock := getDatabaseMock()
	db.SkipDefaultTransaction = true
	mock.MatchExpectationsInOrder(false)
	ss := NewSeedersStack(db)
	ss.commitTransaction(db)
}

func Test_SeedersStack_RollbackTransaction(t *testing.T) {
	db, mock := getDatabaseMock()
	db.SkipDefaultTransaction = true
	mock.MatchExpectationsInOrder(false)
	ss := NewSeedersStack(db)
	ss.rollbackTransaction(db)
}

func Test_SeedersStack_SeedInBatches(t *testing.T) {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","created_at") VALUES ($1,$2,$3,$4),($5,$6,$7,$8),($9,$10,$11,$12),($13,$14,$15,$16),($17,$18,$19,$20) RETURNING "id"`)).
		//WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	ss := NewSeedersStack(db)
	ss.AddSeeder(NewUsersSeederV1(SeederConfiguration{Rows: 5}))
	err := ss.Seed()
	assert.NoError(t, err)
}

func Test_SeedersStack_ClearDelete(t *testing.T) {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM users`)).WillReturnResult(sqlmock.NewResult(0, 1))
	ss := NewSeedersStack(db)
	ss.AddSeeder(NewUsersSeederV1(SeederConfiguration{Rows: 5}))
	err := ss.Clear()
	assert.NoError(t, err)
}

func Test_SeedersStack_SeedInLoopWithTransaction(t *testing.T) {
	db, mock := getDatabaseMock()
	db.SkipDefaultTransaction = true
	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","created_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		//WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	ss := NewSeedersStack(db)
	ss.AddSeeder(NewUsersSeederV2(SeederConfiguration{Rows: 1}))
	err := ss.Seed()
	assert.NoError(t, err)
}

func Test_SeedersStack_ClearTruncate(t *testing.T) {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec(regexp.QuoteMeta(
		`TRUNCATE users`)).WillReturnResult(sqlmock.NewResult(0, 1))
	ss := NewSeedersStack(db)
	ss.AddSeeder(NewUsersSeederV2(SeederConfiguration{Rows: 5}))
	err := ss.Clear()
	assert.NoError(t, err)
}
