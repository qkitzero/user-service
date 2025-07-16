package user

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksidentity "github.com/qkitzero/user-service/mocks/domain/identity"
	mocksuser "github.com/qkitzero/user-service/mocks/domain/user"
	"github.com/qkitzero/user-service/testutil"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, user user.User)
	}{
		{
			name:    "success create user and identity",
			success: true,
			setup: func(mock sqlmock.Sqlmock, user user.User) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","display_name","birth_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(user.ID(), user.DisplayName(), user.BirthDate(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "identities" ("id","user_id") VALUES ($1,$2)`)).
					WithArgs(user.Identities()[0].ID(), user.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create user error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, user user.User) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","display_name","birth_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(user.ID(), user.DisplayName(), user.BirthDate(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnError(errors.New("create user error"))

				mock.ExpectRollback()
			},
		},
		{
			name:    "failure create identity error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, user user.User) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","display_name","birth_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(user.ID(), user.DisplayName(), user.BirthDate(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "identities" ("id","user_id") VALUES ($1,$2)`)).
					WithArgs(user.Identities()[0].ID(), user.ID()).
					WillReturnError(errors.New("create identity error"))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIdentity := mocksidentity.NewMockIdentity(ctrl)
			mockIdentity.EXPECT().ID().Return(identity.IdentityID("google-oauth2|000000000000000000000")).AnyTimes()

			mockUser := mocksuser.NewMockUser(ctrl)
			mockUser.EXPECT().ID().Return(user.UserID{UUID: uuid.New()}).AnyTimes()
			mockUser.EXPECT().Identities().Return([]identity.Identity{mockIdentity}).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(user.DisplayName("test user")).AnyTimes()
			mockUser.EXPECT().BirthDate().Return(user.BirthDate{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}).AnyTimes()
			mockUser.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockUser.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockUser)

			repo := NewUserRepository(gormDB)

			err = repo.Create(mockUser)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByIdentityID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		identityID identity.IdentityID
		setup      func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID)
	}{
		{
			name:       "success find user by identity id",
			success:    true,
			identityID: identity.IdentityID("google-oauth2|000000000000000000000"),
			setup: func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID) {
				identityRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(identityID, userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE id = $1 ORDER BY "identities"."id" LIMIT $2`)).
					WithArgs(identityID, 1).
					WillReturnRows(identityRows)

				userRows := sqlmock.NewRows([]string{"id", "display_name", "birth_date", "created_at", "updated_at"}).
					AddRow(userID, "test user", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnRows(userRows)

				identitiesRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(identityID, userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(identitiesRows)
			},
		},
		{
			name:       "failure identity not found",
			success:    false,
			identityID: identity.IdentityID("google-oauth2|000000000000000000000"),
			setup: func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE id = $1 ORDER BY "identities"."id" LIMIT $2`)).
					WithArgs(identityID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:       "failure find identity error",
			success:    false,
			identityID: identity.IdentityID("google-oauth2|000000000000000000000"),
			setup: func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE id = $1 ORDER BY "identities"."id" LIMIT $2`)).
					WithArgs(identityID, 1).
					WillReturnError(errors.New("find identity error"))
			},
		},
		{
			name:       "failure user not found",
			success:    false,
			identityID: identity.IdentityID("google-oauth2|000000000000000000000"),
			setup: func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID) {
				identityRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(identityID, userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE id = $1 ORDER BY "identities"."id" LIMIT $2`)).
					WithArgs(identityID, 1).
					WillReturnRows(identityRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:       "failure find user error",
			success:    false,
			identityID: identity.IdentityID("google-oauth2|000000000000000000000"),
			setup: func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID) {
				identityRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(identityID, userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE id = $1 ORDER BY "identities"."id" LIMIT $2`)).
					WithArgs(identityID, 1).
					WillReturnRows(identityRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnError(errors.New("find user error"))
			},
		},
		{
			name:       "failure find identities error",
			success:    false,
			identityID: identity.IdentityID("google-oauth2|000000000000000000000"),
			setup: func(mock sqlmock.Sqlmock, identityID identity.IdentityID, userID user.UserID) {
				identityRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(identityID, userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE id = $1 ORDER BY "identities"."id" LIMIT $2`)).
					WithArgs(identityID, 1).
					WillReturnRows(identityRows)

				userRows := sqlmock.NewRows([]string{"id", "display_name", "birth_date", "created_at", "updated_at"}).
					AddRow(userID, "test user", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnRows(userRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnError(errors.New("find identities error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			userID := user.UserID{UUID: uuid.New()}
			tt.setup(mock, tt.identityID, userID)

			repo := NewUserRepository(gormDB)

			_, err = repo.FindByIdentityID(tt.identityID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		userID  user.UserID
		setup   func(mock sqlmock.Sqlmock, userID user.UserID)
	}{
		{
			name:    "success find user by id",
			success: true,
			userID:  user.UserID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				userRows := sqlmock.NewRows([]string{"id", "display_name", "birth_date", "created_at", "updated_at"}).
					AddRow(userID, "test user", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnRows(userRows)

				identitiesRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(identity.IdentityID("google-oauth2|000000000000000000000"), userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(identitiesRows)
			},
		},
		{
			name:    "failure user not found",
			success: false,
			userID:  user.UserID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure find user error",
			success: false,
			userID:  user.UserID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnError(errors.New("find user error"))
			},
		},
		{
			name:    "failure find identities error",
			success: false,
			userID:  user.UserID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				userRows := sqlmock.NewRows([]string{"id", "display_name", "birth_date", "created_at", "updated_at"}).
					AddRow(userID, "test user", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(userID, 1).
					WillReturnRows(userRows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "identities" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnError(errors.New("find identities error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			tt.setup(mock, tt.userID)

			repo := NewUserRepository(gormDB)

			_, err = repo.FindByID(tt.userID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, user user.User)
	}{
		{
			name:    "success update user",
			success: true,
			setup: func(mock sqlmock.Sqlmock, user user.User) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "display_name"=$1,"birth_date"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
					WithArgs(user.DisplayName(), user.BirthDate(), testutil.AnyTime{}, testutil.AnyTime{}, user.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure update user error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, user user.User) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "display_name"=$1,"birth_date"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
					WithArgs(user.DisplayName(), user.BirthDate(), testutil.AnyTime{}, testutil.AnyTime{}, user.ID()).
					WillReturnError(errors.New("update user error"))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIdentity := mocksidentity.NewMockIdentity(ctrl)

			mockUser := mocksuser.NewMockUser(ctrl)
			mockUser.EXPECT().ID().Return(user.UserID{UUID: uuid.New()}).AnyTimes()
			mockUser.EXPECT().Identities().Return([]identity.Identity{mockIdentity}).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(user.DisplayName("updated user")).AnyTimes()
			mockUser.EXPECT().BirthDate().Return(user.BirthDate{Time: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}).AnyTimes()
			mockUser.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockUser.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockUser)

			repo := NewUserRepository(gormDB)

			err = repo.Update(mockUser)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
