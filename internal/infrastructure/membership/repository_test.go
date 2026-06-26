package membership

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
	"github.com/qkitzero/user-service/testutil"
)

func newTestGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to new sqlmock: %s", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}
	return gormDB, mock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, m membership.Membership)
	}{
		{
			name:    "success create membership",
			success: true,
			setup: func(mock sqlmock.Sqlmock, m membership.Membership) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "memberships" ("id","user_id","group_id","role","joined_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(m.ID(), m.UserID(), m.GroupID(), m.Role(), testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create membership error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, m membership.Membership) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "memberships" ("id","user_id","group_id","role","joined_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(m.ID(), m.UserID(), m.GroupID(), m.Role(), testutil.AnyTime{}).
					WillReturnError(errors.New("create membership error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			m := membership.NewMembership(membership.NewMembershipID(), user.NewUserID(), group.NewGroupID(), membership.RoleOwner, time.Now())
			tt.setup(mock, m)

			repo := NewMembershipRepository(gormDB)

			err := repo.Create(context.Background(), m)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByUserAndGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID)
	}{
		{
			name:    "success find membership",
			success: true,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "group_id", "role", "joined_at"}).
					AddRow(membership.NewMembershipID(), userID, groupID, "owner", time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE user_id = $1 AND group_id = $2 ORDER BY "memberships"."id" LIMIT $3`)).
					WithArgs(userID, groupID, 1).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure membership not found",
			success: false,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE user_id = $1 AND group_id = $2 ORDER BY "memberships"."id" LIMIT $3`)).
					WithArgs(userID, groupID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure find error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE user_id = $1 AND group_id = $2 ORDER BY "memberships"."id" LIMIT $3`)).
					WithArgs(userID, groupID, 1).
					WillReturnError(errors.New("find error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			userID := user.NewUserID()
			groupID := group.NewGroupID()
			tt.setup(mock, userID, groupID)

			repo := NewMembershipRepository(gormDB)

			_, err := repo.FindByUserAndGroup(context.Background(), userID, groupID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByGroupID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		wantLen int
		setup   func(mock sqlmock.Sqlmock, groupID group.GroupID)
	}{
		{
			name:    "success list",
			success: true,
			wantLen: 2,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "group_id", "role", "joined_at"}).
					AddRow(membership.NewMembershipID(), user.NewUserID(), groupID, "owner", time.Now()).
					AddRow(membership.NewMembershipID(), user.NewUserID(), groupID, "member", time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE group_id = $1`)).
					WithArgs(groupID).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure list error",
			success: false,
			wantLen: 0,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE group_id = $1`)).
					WithArgs(groupID).
					WillReturnError(errors.New("list error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			groupID := group.NewGroupID()
			tt.setup(mock, groupID)

			repo := NewMembershipRepository(gormDB)

			memberships, err := repo.ListByGroupID(context.Background(), groupID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && len(memberships) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(memberships), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListByUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		wantLen int
		setup   func(mock sqlmock.Sqlmock, userID user.UserID)
	}{
		{
			name:    "success list",
			success: true,
			wantLen: 1,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "group_id", "role", "joined_at"}).
					AddRow(membership.NewMembershipID(), userID, group.NewGroupID(), "owner", time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure list error",
			success: false,
			wantLen: 0,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "memberships" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnError(errors.New("list error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			userID := user.NewUserID()
			tt.setup(mock, userID)

			repo := NewMembershipRepository(gormDB)

			memberships, err := repo.ListByUserID(context.Background(), userID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && len(memberships) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(memberships), tt.wantLen)
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
		setup   func(mock sqlmock.Sqlmock, m membership.Membership)
	}{
		{
			name:    "success update membership",
			success: true,
			setup: func(mock sqlmock.Sqlmock, m membership.Membership) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "memberships" SET "user_id"=$1,"group_id"=$2,"role"=$3,"joined_at"=$4 WHERE "id" = $5`)).
					WithArgs(m.UserID(), m.GroupID(), m.Role(), testutil.AnyTime{}, m.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure update error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, m membership.Membership) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "memberships" SET "user_id"=$1,"group_id"=$2,"role"=$3,"joined_at"=$4 WHERE "id" = $5`)).
					WithArgs(m.UserID(), m.GroupID(), m.Role(), testutil.AnyTime{}, m.ID()).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			m := membership.NewMembership(membership.NewMembershipID(), user.NewUserID(), group.NewGroupID(), membership.RoleAdmin, time.Now())
			tt.setup(mock, m)

			repo := NewMembershipRepository(gormDB)

			err := repo.Update(context.Background(), m)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID)
	}{
		{
			name:    "success delete membership",
			success: true,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "memberships" WHERE user_id = $1 AND group_id = $2`)).
					WithArgs(userID, groupID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure delete error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, userID user.UserID, groupID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "memberships" WHERE user_id = $1 AND group_id = $2`)).
					WithArgs(userID, groupID).
					WillReturnError(errors.New("delete error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			userID := user.NewUserID()
			groupID := group.NewGroupID()
			tt.setup(mock, userID, groupID)

			repo := NewMembershipRepository(gormDB)

			err := repo.Delete(context.Background(), userID, groupID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCountOwners(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		wantCount int
		setup     func(mock sqlmock.Sqlmock, groupID group.GroupID)
	}{
		{
			name:      "success count owners",
			success:   true,
			wantCount: 2,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "memberships" WHERE group_id = $1 AND role = $2`)).
					WithArgs(groupID, membership.RoleOwner).
					WillReturnRows(rows)
			},
		},
		{
			name:      "failure count error",
			success:   false,
			wantCount: 0,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "memberships" WHERE group_id = $1 AND role = $2`)).
					WithArgs(groupID, membership.RoleOwner).
					WillReturnError(errors.New("count error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			groupID := group.NewGroupID()
			tt.setup(mock, groupID)

			repo := NewMembershipRepository(gormDB)

			count, err := repo.CountOwners(context.Background(), groupID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && count != tt.wantCount {
				t.Errorf("count = %d, want %d", count, tt.wantCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
