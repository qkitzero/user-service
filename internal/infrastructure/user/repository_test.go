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

	"github.com/qkitzero/user/internal/domain/identity"
	"github.com/qkitzero/user/internal/domain/user"
	mocksidentity "github.com/qkitzero/user/mocks/domain/identity"
	mocksuser "github.com/qkitzero/user/mocks/domain/user"
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
					WithArgs(user.ID(), user.DisplayName(), user.BirthDate(), user.CreatedAt(), user.UpdatedAt()).
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
					WithArgs(user.ID(), user.DisplayName(), user.BirthDate(), user.CreatedAt(), user.UpdatedAt()).
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
					WithArgs(user.ID(), user.DisplayName(), user.BirthDate(), user.CreatedAt(), user.UpdatedAt()).
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
