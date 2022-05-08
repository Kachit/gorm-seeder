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
