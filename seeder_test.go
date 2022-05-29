package gorm_seeder

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type SeederAbstractTestSuite struct {
	suite.Suite
	db       *gorm.DB
	mock     sqlmock.Sqlmock
	testable SeederAbstract
}

func (suite *SeederAbstractTestSuite) SetupTest() {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	suite.db = db
	suite.mock = mock
	suite.testable = NewSeederAbstract(SeederConfiguration{})
}

func (suite *SeederAbstractTestSuite) TestDelete() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM foo`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Delete(suite.db, "foo")
	assert.NoError(suite.T(), err)
}

func (suite *SeederAbstractTestSuite) TestTruncate() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`TRUNCATE foo`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Truncate(suite.db, "foo")
	assert.NoError(suite.T(), err)
}

func TestSeederAbstractTestSuite(t *testing.T) {
	suite.Run(t, new(SeederAbstractTestSuite))
}

type SeedersStackTestSuite struct {
	suite.Suite
	db       *gorm.DB
	mock     sqlmock.Sqlmock
	testable *SeedersStack
}

func (suite *SeedersStackTestSuite) SetupTest() {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	suite.db = db
	suite.mock = mock
	suite.testable = NewSeedersStack(db)
}

func (suite *SeedersStackTestSuite) TestAddSeeder() {
	assert.Len(suite.T(), suite.testable.seeders, 0)
	suite.testable.AddSeeder(NewUsersSeederV1(SeederConfiguration{}))
	assert.Len(suite.T(), suite.testable.seeders, 1)
}

func (suite *SeedersStackTestSuite) TestBeginTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.testable.beginTransaction(suite.db)
}

func (suite *SeedersStackTestSuite) TestCommitTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.testable.commitTransaction(suite.db)
}

func (suite *SeedersStackTestSuite) TestRollbackTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.testable.rollbackTransaction(suite.db)
}

func (suite *SeedersStackTestSuite) TestSeedInBatches() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","created_at") VALUES ($1,$2,$3,$4),($5,$6,$7,$8),($9,$10,$11,$12),($13,$14,$15,$16),($17,$18,$19,$20) RETURNING "id"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	suite.mock.ExpectCommit()
	suite.testable.AddSeeder(NewUsersSeederV1(SeederConfiguration{Rows: 5}))
	err := suite.testable.Seed()
	assert.NoError(suite.T(), err)
}

func (suite *SeedersStackTestSuite) TestSeedInLoopWithTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","created_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	suite.mock.ExpectCommit()
	suite.testable.AddSeeder(NewUsersSeederV2(SeederConfiguration{Rows: 1}))
	err := suite.testable.Seed()
	assert.NoError(suite.T(), err)
}

func (suite *SeedersStackTestSuite) TestClearDelete() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM users`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	suite.testable.AddSeeder(NewUsersSeederV1(SeederConfiguration{Rows: 5}))
	err := suite.testable.Clear()
	assert.NoError(suite.T(), err)
}

func (suite *SeedersStackTestSuite) TestClearTruncate() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`TRUNCATE users`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	suite.testable.AddSeeder(NewUsersSeederV2(SeederConfiguration{Rows: 5}))
	err := suite.testable.Clear()
	assert.NoError(suite.T(), err)
}

func TestSeedersStackTestSuite(t *testing.T) {
	suite.Run(t, new(SeedersStackTestSuite))
}
