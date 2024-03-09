package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/jmehdipour/gift-card/internal/domain"
)

type GiftCardRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *giftCardRepository
}

func (suite *GiftCardRepositoryTestSuite) SetupTest() {
	suite.db, suite.mock, _ = sqlmock.New()
	suite.repo = &giftCardRepository{
		db: suite.db,
	}
}

func (suite *GiftCardRepositoryTestSuite) TeardownTest() {
	_ = suite.db.Close()
}

func (suite *GiftCardRepositoryTestSuite) TestNewGiftCardRepository() {
	require := suite.Require()

	db, _, _ := sqlmock.New()
	repo := NewGiftCardRepository(db)

	require.NotNil(repo)
}

func (suite *GiftCardRepositoryTestSuite) TestCreate_Success() {
	require := suite.Require()
	id := uint(101)
	g := &domain.GiftCard{
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
	}

	suite.mock.ExpectExec("^INSERT INTO gift_cards").
		WithArgs(g.Amount, g.GifterID, g.GifteeID).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))

	err := suite.repo.Create(g)

	require.NoError(err)
	require.Equal(id, g.ID)
}

func (suite *GiftCardRepositoryTestSuite) TestCreate_Failure() {
	require := suite.Require()
	g := &domain.GiftCard{
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
	}
	expectedError := errors.New("error in inserting to gift_cards table")

	suite.mock.ExpectExec("^INSERT INTO gift_cards").
		WithArgs(g.Amount, g.GifterID, g.GifteeID).
		WillReturnError(expectedError)

	err := suite.repo.Create(g)

	require.EqualError(err, expectedError.Error())
}

func (suite *GiftCardRepositoryTestSuite) TestCreate_LastInsertIdError_Failure() {
	require := suite.Require()
	g := &domain.GiftCard{
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
	}
	expectedError := errors.New("LastInsertId error")

	suite.mock.ExpectExec("^INSERT INTO gift_cards").
		WithArgs(g.Amount, g.GifterID, g.GifteeID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId error")))

	err := suite.repo.Create(g)

	require.Error(err)
	require.EqualError(err, expectedError.Error())
}

func (suite *GiftCardRepositoryTestSuite) TestFindByID_DBError_Failure() {
	require := suite.Require()
	expectedError := "database failure"
	id := uint(10)

	suite.mock.ExpectQuery("^SELECT .+ FROM gift_cards").
		WithArgs(id).
		WillReturnError(errors.New("database failure"))

	result, err := suite.repo.FindByID(id)

	require.Error(err)
	require.EqualError(err, expectedError)
	require.Empty(result)
}

func (suite *GiftCardRepositoryTestSuite) TestFindByID_NotFound() {
	require := suite.Require()
	id := uint(10)

	suite.mock.ExpectQuery("^SELECT .+ FROM gift_cards").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	result, err := suite.repo.FindByID(id)
	require.NoError(err)
	require.Equal((*domain.GiftCard)(nil), result)
}

func (suite *GiftCardRepositoryTestSuite) TestFindByID_Success() {
	require := suite.Require()
	id := uint(10)
	expectedResult := &domain.GiftCard{
		ID:       10,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
		Status:   1,
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "status"}).AddRow(10, 10, 20, 100, 1)
	suite.mock.ExpectQuery("^SELECT .+ FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)

	result, err := suite.repo.FindByID(id)
	require.NoError(err)
	require.Equal(expectedResult, result)
}

func (suite *GiftCardRepositoryTestSuite) TestUpdateStatus_Success() {
	require := suite.Require()
	id := uint(101)
	status := domain.GCSAccepted

	suite.mock.ExpectExec("^UPDATE gift_cards SET status").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(101, 1))

	err := suite.repo.UpdateStatus(id, status)

	require.NoError(err)
}

func (suite *GiftCardRepositoryTestSuite) TestUpdateStatus_DBerror_Failure() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedError := errors.New("something went wrong")

	suite.mock.ExpectExec("^UPDATE gift_cards SET status").
		WithArgs(status, id).
		WillReturnError(expectedError)

	err := suite.repo.UpdateStatus(id, status)

	require.Equal(expectedError, err)
}

func (suite *GiftCardRepositoryTestSuite) TestFindReceivedGiftCardsByUserID_FindDBError_Failure() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedError := errors.New("something went wrong")

	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnError(expectedError)

	giftCards, total, err := suite.repo.FindReceivedGiftCardsByUserID(id, &status, 10, 1)

	require.Equal(expectedError, err)
	require.Equal(total, 0)
	require.Empty(giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindReceivedGiftCardsByUserID_CountDBError_Failure() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedError := errors.New("something went wrong")

	rows := sqlmock.NewRows([]string{"id", "status", "sender_id", "receiver_id", "amount"}).AddRow(10, 1, 10, 20, 100)
	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM gift_cards WHERE receiver_id = \? AND status = 0$`).
		WithArgs(id).
		WillReturnError(expectedError)

	giftCards, total, err := suite.repo.FindReceivedGiftCardsByUserID(id, &status, 10, 1)

	require.EqualError(expectedError, err.Error())
	require.Zero(total)
	require.Empty(giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindReceivedGiftCardsByUserID_Success() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedTotal := 1
	expectedResult := []domain.GiftCard{{
		ID:       10,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
		Status:   1,
	}}

	rows := sqlmock.NewRows([]string{"id", "status", "sender_id", "receiver_id", "amount"}).AddRow(10, 1, 10, 20, 100)
	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM gift_cards WHERE receiver_id = \? AND status = 0$`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

	giftCards, total, err := suite.repo.FindReceivedGiftCardsByUserID(id, &status, 10, 1)

	require.NoError(err)
	require.Equal(expectedTotal, total)
	require.Equal(expectedResult, giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindReceivedGiftCardsByUserID_WithoutStatus_Success() {
	require := suite.Require()
	id := uint(102)
	expectedTotal := 1
	expectedResult := []domain.GiftCard{{
		ID:       10,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
		Status:   1,
	}}

	rows := sqlmock.NewRows([]string{"id", "status", "sender_id", "receiver_id", "amount"}).AddRow(10, 1, 10, 20, 100)
	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM gift_cards WHERE receiver_id = \?$`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

	giftCards, total, err := suite.repo.FindReceivedGiftCardsByUserID(id, nil, 10, 1)

	require.NoError(err)
	require.Equal(expectedTotal, total)
	require.Equal(expectedResult, giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindSentGiftCardsByUserID_FindDBError_Failure() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedError := errors.New("something went wrong")

	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnError(expectedError)

	giftCards, total, err := suite.repo.FindSentGiftCardsByUserID(id, &status, 10, 1)

	require.Equal(expectedError, err)
	require.Equal(total, 0)
	require.Empty(giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindSentGiftCardsByUserID_CountDBError_Failure() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedError := errors.New("something went wrong")

	rows := sqlmock.NewRows([]string{"id", "status", "sender_id", "receiver_id", "amount"}).AddRow(10, 1, 10, 20, 100)
	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM gift_cards WHERE sender_id = \? AND status = 0$`).
		WithArgs(id).
		WillReturnError(expectedError)

	giftCards, total, err := suite.repo.FindSentGiftCardsByUserID(id, &status, 10, 1)

	require.EqualError(expectedError, err.Error())
	require.Zero(total)
	require.Empty(giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindSentGiftCardsByUserID_Success() {
	require := suite.Require()
	id := uint(102)
	status := domain.GCSAccepted
	expectedTotal := 1
	expectedResult := []domain.GiftCard{{
		ID:       10,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
		Status:   1,
	}}

	rows := sqlmock.NewRows([]string{"id", "status", "sender_id", "receiver_id", "amount"}).AddRow(10, 1, 10, 20, 100)
	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM gift_cards WHERE sender_id = \? AND status = 0$`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

	giftCards, total, err := suite.repo.FindSentGiftCardsByUserID(id, &status, 10, 1)

	require.NoError(err)
	require.Equal(expectedTotal, total)
	require.Equal(expectedResult, giftCards)
}

func (suite *GiftCardRepositoryTestSuite) TestFindSentGiftCardsByUserID_WithoutStatus_Success() {
	require := suite.Require()
	id := uint(102)
	expectedTotal := 1
	expectedResult := []domain.GiftCard{{
		ID:       10,
		GifterID: 10,
		GifteeID: 20,
		Amount:   100,
		Status:   1,
	}}

	rows := sqlmock.NewRows([]string{"id", "status", "sender_id", "receiver_id", "amount"}).AddRow(10, 1, 10, 20, 100)
	suite.mock.ExpectQuery("^SELECT .* FROM gift_cards").
		WithArgs(id).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`^SELECT COUNT\(\*\) FROM gift_cards WHERE sender_id = \?$`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

	giftCards, total, err := suite.repo.FindSentGiftCardsByUserID(id, nil, 10, 1)

	require.NoError(err)
	require.Equal(expectedTotal, total)
	require.Equal(expectedResult, giftCards)
}

func TestGiftCardRepository(t *testing.T) {
	suite.Run(t, new(GiftCardRepositoryTestSuite))
}
