package gorm_seeder

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type SeederAbstract_TestSuite struct {
	suite.Suite
	db       *gorm.DB
	mock     sqlmock.Sqlmock
	testable SeederAbstract
}

func (suite *SeederAbstract_TestSuite) SetupTest() {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	suite.db = db
	suite.mock = mock
	suite.testable = NewSeederAbstract(SeederConfiguration{})
}

func (suite *SeederAbstract_TestSuite) TestDelete() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM foo`)).WillReturnResult(sqlmock.NewResult(0, 1))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Delete(suite.db, "foo")
	assert.NoError(suite.T(), err)
}

func (suite *SeederAbstract_TestSuite) TestTruncate() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`TRUNCATE foo`)).WillReturnResult(sqlmock.NewResult(0, 1))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Truncate(suite.db, "foo")
	assert.NoError(suite.T(), err)
}

func Test_SeederAbstract_TestSuite(t *testing.T) {
	suite.Run(t, new(SeederAbstract_TestSuite))
}

type SeedersStack_TestSuite struct {
	suite.Suite
	db       *gorm.DB
	mock     sqlmock.Sqlmock
	testable *SeedersStack
}

func (suite *SeedersStack_TestSuite) SetupTest() {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)
	suite.db = db
	suite.mock = mock
	suite.testable = NewSeedersStack(db)
}

func (suite *SeedersStack_TestSuite) TestAddSeeder() {
	assert.Len(suite.T(), suite.testable.seeders, 0)
	suite.testable.AddSeeder(NewUsersSeederV1(SeederConfiguration{}))
	assert.Len(suite.T(), suite.testable.seeders, 1)
}

func (suite *SeedersStack_TestSuite) TestBeginTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.testable.beginTransaction(suite.db)
}

func (suite *SeedersStack_TestSuite) TestCommitTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.testable.commitTransaction(suite.db)
}

func (suite *SeedersStack_TestSuite) TestRollbackTransaction() {
	suite.db.SkipDefaultTransaction = true
	suite.testable.rollbackTransaction(suite.db)
}

func (suite *SeedersStack_TestSuite) TestSeedInBatches() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","created_at") VALUES ($1,$2,$3,$4),($5,$6,$7,$8),($9,$10,$11,$12),($13,$14,$15,$16),($17,$18,$19,$20) RETURNING "id"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	suite.mock.ExpectCommit()
	suite.testable.AddSeeder(NewUsersSeederV1(SeederConfiguration{Rows: 5}))
	err := suite.testable.Seed()
	assert.NoError(suite.T(), err)
}

func (suite *SeedersStack_TestSuite) TestSeedInLoopWithTransaction() {
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

func (suite *SeedersStack_TestSuite) TestClearDelete() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM users`)).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.testable.AddSeeder(NewUsersSeederV1(SeederConfiguration{Rows: 5}))
	err := suite.testable.Clear()
	assert.NoError(suite.T(), err)
}

func (suite *SeedersStack_TestSuite) TestClearTruncate() {
	suite.mock.ExpectExec(regexp.QuoteMeta(
		`TRUNCATE users`)).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.testable.AddSeeder(NewUsersSeederV2(SeederConfiguration{Rows: 5}))
	err := suite.testable.Clear()
	assert.NoError(suite.T(), err)
}

func Test_SeedersStack_TestSuite(t *testing.T) {
	suite.Run(t, new(SeedersStack_TestSuite))
}
