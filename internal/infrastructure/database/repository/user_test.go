package repository_test

import (
	"context"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/Beriw98/user-management/ent"
	"github.com/Beriw98/user-management/internal/app/domain"
	"github.com/Beriw98/user-management/internal/infrastructure/database/repository"
)

func mockDbClient() (*ent.Client, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	drv := sql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	return client, mock
}

func TestNewUserRepository(t *testing.T) {
	t.Run("NewUserRepository", func(t *testing.T) {
		c, _ := mockDbClient()
		got := repository.NewUserRepository(c)

		assert.NotNil(t, got)
	})
}

func TestUser_Create(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		user := domain.User{
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "test",
		}

		mock.ExpectExec("INSERT INTO \"users\"").
			WithArgs(user.Name, user.Surname, user.Email, user.Password, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := userRepo.Create(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("Create error", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		user := domain.User{
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "test",
		}

		mock.ExpectExec("INSERT INTO \"users\"").
			WithArgs(user.Name, user.Surname, user.Email, user.Password, sqlmock.AnyArg()).
			WillReturnError(assert.AnError)

		err := userRepo.Create(ctx, user)
		assert.Error(t, err)
	})
}

func TestUser_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		id := "1"

		mock.ExpectExec("DELETE FROM \"users\"").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := userRepo.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Delete error", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		id := "1"

		mock.ExpectExec("DELETE FROM \"users\"").
			WithArgs(id).
			WillReturnError(assert.AnError)

		err := userRepo.Delete(ctx, id)
		assert.Error(t, err)
	})
}

func TestUser_GetByID(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		id := "1"
		user := domain.User{
			ID:      id,
			Name:    "Test",
			Surname: "Test",
			Email:   "test@test.pl",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "email"}).
			AddRow(user.ID, user.Name, user.Surname, user.Email)

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WithArgs(id).
			WillReturnRows(rows)

		got, err := userRepo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, &user, got)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		id := "1"

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "email"})

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WithArgs(id).
			WillReturnRows(rows)

		got, err := userRepo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetByID error", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		id := "1"

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WithArgs(id).
			WillReturnError(assert.AnError)

		got, err := userRepo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

}

func TestUser_GetMany(t *testing.T) {
	t.Run("GetMany", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		users := []domain.User{
			{
				ID:      "1",
				Name:    "Test",
				Surname: "Test",
				Email:   "test@test.pl",
			},
		}

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "email"}).
			AddRow(users[0].ID, users[0].Name, users[0].Surname, users[0].Email)

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WillReturnRows(rows)

		got, err := userRepo.GetMany(ctx, 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, users, got)
	})

	t.Run("GetMany error", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WillReturnError(assert.AnError)

		got, err := userRepo.GetMany(ctx, 10, 0)

		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestUser_Update(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		user := domain.User{
			ID:       "1",
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "testPassword",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE \"users\"").
			WithArgs(user.Name, user.Surname, user.Email, user.Password, user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT \"id\", \"name\", \"surname\", \"email\", \"password\" FROM \"users\"").
			WithArgs(user.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "email", "password"}).
				AddRow(user.ID, user.Name, user.Surname, user.Email, user.Password))

		mock.ExpectCommit()

		err := userRepo.Update(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("Update error", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		user := domain.User{
			ID:       "1",
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "test",
		}

		mock.ExpectExec("UPDATE \"users\"").
			WithArgs(user.Name, user.Surname, user.Email, user.ID).
			WillReturnError(assert.AnError)

		err := userRepo.Update(ctx, user)
		assert.Error(t, err)
	})
}

func TestUser_GetByEmail(t *testing.T) {
	t.Run("GetByEmail", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		email := "tet@test.pl"
		user := domain.User{
			ID:      "1",
			Name:    "Test",
			Surname: "Test",
			Email:   email,
		}

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "email"}).
			AddRow(user.ID, user.Name, user.Surname, user.Email)

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WithArgs(email).
			WillReturnRows(rows)

		got, err := userRepo.GetByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, &user, got)

	})

	t.Run("GetByEmail not found", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		email := "test@test.pl"

		rows := sqlmock.NewRows([]string{"id", "name", "surname", "email"})

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WithArgs(email).
			WillReturnRows(rows)

		got, err := userRepo.GetByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetByEmail error", func(t *testing.T) {
		client, mock := mockDbClient()
		userRepo := repository.NewUserRepository(client)
		ctx := context.Background()

		email := "test@test.pl"

		mock.ExpectQuery("SELECT \"users\".\"id\", \"users\".\"name\", \"users\".\"surname\", \"users\".\"email\", \"users\".\"password\" FROM \"users\"").
			WithArgs(email).
			WillReturnError(assert.AnError)

		got, err := userRepo.GetByEmail(ctx, email)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
