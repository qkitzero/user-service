package group

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksgroup "github.com/qkitzero/user-service/mocks/domain/group"
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
		setup   func(mock sqlmock.Sqlmock, g group.Group, ownerID user.UserID)
	}{
		{
			name:    "success create group with owner",
			success: true,
			setup: func(mock sqlmock.Sqlmock, g group.Group, ownerID user.UserID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "groups" ("id","name","created_at","updated_at") VALUES ($1,$2,$3,$4)`)).
					WithArgs(g.ID(), g.Name(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "memberships" ("id","user_id","group_id","role","joined_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(sqlmock.AnyArg(), ownerID, g.ID(), membership.RoleOwner, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create group error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, g group.Group, ownerID user.UserID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "groups" ("id","name","created_at","updated_at") VALUES ($1,$2,$3,$4)`)).
					WithArgs(g.ID(), g.Name(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnError(errors.New("create group error"))
				mock.ExpectRollback()
			},
		},
		{
			name:    "failure create owner membership error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, g group.Group, ownerID user.UserID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "groups" ("id","name","created_at","updated_at") VALUES ($1,$2,$3,$4)`)).
					WithArgs(g.ID(), g.Name(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "memberships" ("id","user_id","group_id","role","joined_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(sqlmock.AnyArg(), ownerID, g.ID(), membership.RoleOwner, testutil.AnyTime{}).
					WillReturnError(errors.New("create owner membership error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGroup := mocksgroup.NewMockGroup(ctrl)
			mockGroup.EXPECT().ID().Return(group.NewGroupID()).AnyTimes()
			mockGroup.EXPECT().Name().Return(group.GroupName("test group")).AnyTimes()
			mockGroup.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockGroup.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

			ownerID := user.NewUserID()
			tt.setup(mock, mockGroup, ownerID)

			repo := NewGroupRepository(gormDB)

			err := repo.Create(context.Background(), mockGroup, ownerID)
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

func TestFindByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, groupID group.GroupID)
	}{
		{
			name:    "success find group by id",
			success: true,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(groupID, "test group", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE id = $1 ORDER BY "groups"."id" LIMIT $2`)).
					WithArgs(groupID, 1).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure group not found",
			success: false,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE id = $1 ORDER BY "groups"."id" LIMIT $2`)).
					WithArgs(groupID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure find group error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE id = $1 ORDER BY "groups"."id" LIMIT $2`)).
					WithArgs(groupID, 1).
					WillReturnError(errors.New("find group error"))
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

			repo := NewGroupRepository(gormDB)

			_, err := repo.FindByID(context.Background(), groupID)
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

func TestUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, g group.Group)
	}{
		{
			name:    "success update group",
			success: true,
			setup: func(mock sqlmock.Sqlmock, g group.Group) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "groups" SET "name"=$1,"created_at"=$2,"updated_at"=$3 WHERE "id" = $4`)).
					WithArgs(g.Name(), testutil.AnyTime{}, testutil.AnyTime{}, g.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure update group error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, g group.Group) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "groups" SET "name"=$1,"created_at"=$2,"updated_at"=$3 WHERE "id" = $4`)).
					WithArgs(g.Name(), testutil.AnyTime{}, testutil.AnyTime{}, g.ID()).
					WillReturnError(errors.New("update group error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGroup := mocksgroup.NewMockGroup(ctrl)
			mockGroup.EXPECT().ID().Return(group.NewGroupID()).AnyTimes()
			mockGroup.EXPECT().Name().Return(group.GroupName("updated group")).AnyTimes()
			mockGroup.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockGroup.EXPECT().UpdatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockGroup)

			repo := NewGroupRepository(gormDB)

			err := repo.Update(context.Background(), mockGroup)
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
		setup   func(mock sqlmock.Sqlmock, groupID group.GroupID)
	}{
		{
			name:    "success delete group",
			success: true,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "groups" WHERE id = $1`)).
					WithArgs(groupID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure delete group error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, groupID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "groups" WHERE id = $1`)).
					WithArgs(groupID).
					WillReturnError(errors.New("delete group error"))
				mock.ExpectRollback()
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

			repo := NewGroupRepository(gormDB)

			err := repo.Delete(context.Background(), groupID)
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

func TestAddChild(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, parentID, childID group.GroupID)
	}{
		{
			name:    "success add child",
			success: true,
			setup: func(mock sqlmock.Sqlmock, parentID, childID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "group_relations" ("parent_id","child_id") VALUES ($1,$2)`)).
					WithArgs(parentID, childID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure add child error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, parentID, childID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "group_relations" ("parent_id","child_id") VALUES ($1,$2)`)).
					WithArgs(parentID, childID).
					WillReturnError(errors.New("add child error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			parentID := group.NewGroupID()
			childID := group.NewGroupID()
			tt.setup(mock, parentID, childID)

			repo := NewGroupRepository(gormDB)

			err := repo.AddChild(context.Background(), parentID, childID)
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

func TestRemoveChild(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, parentID, childID group.GroupID)
	}{
		{
			name:    "success remove child",
			success: true,
			setup: func(mock sqlmock.Sqlmock, parentID, childID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "group_relations" WHERE parent_id = $1 AND child_id = $2`)).
					WithArgs(parentID, childID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure remove child error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, parentID, childID group.GroupID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "group_relations" WHERE parent_id = $1 AND child_id = $2`)).
					WithArgs(parentID, childID).
					WillReturnError(errors.New("remove child error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			parentID := group.NewGroupID()
			childID := group.NewGroupID()
			tt.setup(mock, parentID, childID)

			repo := NewGroupRepository(gormDB)

			err := repo.RemoveChild(context.Background(), parentID, childID)
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

func TestFindChildren(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		wantLen int
		setup   func(mock sqlmock.Sqlmock, parentID group.GroupID)
	}{
		{
			name:    "success find children",
			success: true,
			wantLen: 1,
			setup: func(mock sqlmock.Sqlmock, parentID group.GroupID) {
				childID := group.NewGroupID()
				relationRows := sqlmock.NewRows([]string{"parent_id", "child_id"}).AddRow(parentID, childID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "group_relations" WHERE parent_id = $1`)).
					WithArgs(parentID).
					WillReturnRows(relationRows)
				groupRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(childID, "child group", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE id IN ($1)`)).
					WithArgs(childID).
					WillReturnRows(groupRows)
			},
		},
		{
			name:    "success no children",
			success: true,
			wantLen: 0,
			setup: func(mock sqlmock.Sqlmock, parentID group.GroupID) {
				relationRows := sqlmock.NewRows([]string{"parent_id", "child_id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "group_relations" WHERE parent_id = $1`)).
					WithArgs(parentID).
					WillReturnRows(relationRows)
			},
		},
		{
			name:    "failure find relations error",
			success: false,
			wantLen: 0,
			setup: func(mock sqlmock.Sqlmock, parentID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "group_relations" WHERE parent_id = $1`)).
					WithArgs(parentID).
					WillReturnError(errors.New("find relations error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			parentID := group.NewGroupID()
			tt.setup(mock, parentID)

			repo := NewGroupRepository(gormDB)

			groups, err := repo.FindChildren(context.Background(), parentID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && len(groups) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(groups), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindParents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		wantLen int
		setup   func(mock sqlmock.Sqlmock, childID group.GroupID)
	}{
		{
			name:    "success find parents",
			success: true,
			wantLen: 1,
			setup: func(mock sqlmock.Sqlmock, childID group.GroupID) {
				parentID := group.NewGroupID()
				relationRows := sqlmock.NewRows([]string{"parent_id", "child_id"}).AddRow(parentID, childID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "group_relations" WHERE child_id = $1`)).
					WithArgs(childID).
					WillReturnRows(relationRows)
				groupRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(parentID, "parent group", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE id IN ($1)`)).
					WithArgs(parentID).
					WillReturnRows(groupRows)
			},
		},
		{
			name:    "success no parents",
			success: true,
			wantLen: 0,
			setup: func(mock sqlmock.Sqlmock, childID group.GroupID) {
				relationRows := sqlmock.NewRows([]string{"parent_id", "child_id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "group_relations" WHERE child_id = $1`)).
					WithArgs(childID).
					WillReturnRows(relationRows)
			},
		},
		{
			name:    "failure find relations error",
			success: false,
			wantLen: 0,
			setup: func(mock sqlmock.Sqlmock, childID group.GroupID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "group_relations" WHERE child_id = $1`)).
					WithArgs(childID).
					WillReturnError(errors.New("find relations error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := newTestGormDB(t)

			childID := group.NewGroupID()
			tt.setup(mock, childID)

			repo := NewGroupRepository(gormDB)

			groups, err := repo.FindParents(context.Background(), childID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && len(groups) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(groups), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
