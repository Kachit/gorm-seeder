package gorm_seeder

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_SeederAbstract_Delete(t *testing.T) {
	db, mock := getDatabaseMock()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(regexp.QuoteMeta(
		`DELETE FROM foo`))

	sa := NewSeederAbstract(SeederConfiguration{})
	err := sa.Delete(db, "foo")
	assert.NoError(t, err)
}
