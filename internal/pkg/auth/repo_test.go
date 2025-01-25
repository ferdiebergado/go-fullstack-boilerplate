package auth_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/auth"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

const testEmail = "test@example.com"
const testPassword = "testpassword"
const testPasswordHash = "testpassword_hashed"

var mockOpts = sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
var dbConfig = config.DBConfig{}
var signupParams = auth.SignUpParams{
	Email:                testEmail,
	Password:             testPassword,
	PasswordConfirmation: testPassword,
}

func TestRepoSignUpHappyPath(t *testing.T) {
	mockDB, mock, err := sqlmock.New(mockOpts)
	if err != nil {
		t.Fatalf("failed to create mock database: %s", err)
	}
	defer mockDB.Close()

	now := time.Now()
	cols := []string{"id", "email", "auth_method", "created_at", "updated_at"}
	row := []driver.Value{"1", testEmail, user.BasicAuth, now, now}

	mock.ExpectQuery(auth.SignUpQuery).
		WithArgs(testEmail, testPassword, user.BasicAuth).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(row...))

	repo := auth.NewAuthRepo(&dbConfig, mockDB)

	got, err := repo.SignUp(context.Background(), signupParams)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := &user.User{
		Email:      testEmail,
		AuthMethod: user.BasicAuth,
	}

	if got.Email != want.Email {
		t.Errorf("want: %v but got: %v", want.Email, got.Email)
	}

	if string(got.AuthMethod) != string(want.AuthMethod) {
		t.Errorf("want: %v but got: %v", want.AuthMethod, got.AuthMethod)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestRepoSignUpDuplicateEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New(mockOpts)
	if err != nil {
		t.Fatalf("failed to create mock database: %s", err)
	}
	defer mockDB.Close()

	errEmail := auth.EmailExistsError{Email: testEmail}

	mock.ExpectQuery(auth.SignUpQuery).
		WithArgs(testEmail, testPassword, user.BasicAuth).WillReturnError(&errEmail)

	repo := auth.NewAuthRepo(&dbConfig, mockDB)

	_, err = repo.SignUp(context.Background(), signupParams)

	if err == nil {
		t.Error("expected an error but got nil")
	}

	var emailExistsErr *auth.EmailExistsError
	if !errors.As(err, &emailExistsErr) {
		t.Errorf("want: %v but got: %v", emailExistsErr, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestRepoSignUpDBError(t *testing.T) {
	mockDB, mock, err := sqlmock.New(mockOpts)
	if err != nil {
		t.Fatalf("failed to create mock database: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery(auth.SignUpQuery).
		WithArgs(testEmail, testPassword, user.BasicAuth).WillReturnError(sql.ErrNoRows)

	repo := auth.NewAuthRepo(&dbConfig, mockDB)

	_, err = repo.SignUp(context.Background(), signupParams)

	if err == nil {
		t.Error("expected an error but got nil")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("want: %v but got: %v", sql.ErrNoRows, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestRepoSignInHappyPath(t *testing.T) {
	mockDB, mock, err := sqlmock.New(mockOpts)
	if err != nil {
		t.Fatalf("failed to create mock database: %s", err)
	}
	defer mockDB.Close()

	cols := []string{"id", "password_hash"}
	row := []driver.Value{"1", testPasswordHash}

	mock.ExpectQuery(auth.SignInQuery).
		WithArgs(testEmail).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(row...))

	repo := auth.NewAuthRepo(&dbConfig, mockDB)

	got, err := repo.SignIn(context.Background(), testEmail)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := &auth.SignInResult{
		ID:   "1",
		Hash: testPasswordHash,
	}

	if got.ID != want.ID {
		t.Errorf("want: %v but got: %v", want.ID, got.ID)
	}

	if got.Hash != want.Hash {
		t.Errorf("want: %v but got: %v", want.Hash, got.Hash)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestRepoSignInEmailNotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New(mockOpts)
	if err != nil {
		t.Fatalf("failed to create mock database: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery(auth.SignInQuery).
		WithArgs(testEmail).
		WillReturnError(sql.ErrNoRows)

	repo := auth.NewAuthRepo(&dbConfig, mockDB)

	_, err = repo.SignIn(context.Background(), testEmail)

	if err == nil {
		t.Error("expected an error but got nil")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("want: %v but got: %v", sql.ErrNoRows, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestRepoSignInDBError(t *testing.T) {
	mockDB, mock, err := sqlmock.New(mockOpts)
	if err != nil {
		t.Fatalf("failed to create mock database: %s", err)
	}
	defer mockDB.Close()
	dbErr := errors.New("database error")

	mock.ExpectQuery(auth.SignInQuery).
		WithArgs(testEmail).
		WillReturnError(dbErr)

	repo := auth.NewAuthRepo(&dbConfig, mockDB)

	_, err = repo.SignIn(context.Background(), testEmail)

	if err == nil {
		t.Error("expected an error but got nil")
	}

	if !errors.Is(err, dbErr) {
		t.Errorf("want: %v but got: %v", dbErr, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
