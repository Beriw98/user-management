package repository

import (
	"context"

	"github.com/Beriw98/user-management/ent"
	entuser "github.com/Beriw98/user-management/ent/user"
	"github.com/Beriw98/user-management/internal/app/domain"
)

type User struct {
	Client *ent.UserClient
}

func NewUserRepository(client *ent.Client) *User {
	return &User{
		Client: client.User,
	}
}

func (u *User) Create(ctx context.Context, user domain.User) error {
	_, err := u.Client.Create().
		SetID(user.ID).
		SetName(user.Name).
		SetSurname(user.Surname).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(ctx)

	return err
}

func (u *User) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := u.Client.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.Client.Query().Where(entuser.Email(email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (u *User) Update(ctx context.Context, user domain.User) error {
	_, err := u.Client.UpdateOneID(user.ID).
		SetName(user.Name).
		SetSurname(user.Surname).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(ctx)
	return err
}

func (u *User) Delete(ctx context.Context, id string) error {
	return u.Client.DeleteOneID(id).Exec(ctx)
}

func (u *User) GetMany(ctx context.Context, limit, offset int) ([]domain.User, error) {
	users, err := u.Client.Query().Limit(limit).Offset(offset).Order(entuser.ByID()).All(ctx)
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range users {
		domainUsers = append(domainUsers, domain.User{
			ID:       user.ID,
			Name:     user.Name,
			Surname:  user.Surname,
			Email:    user.Email,
			Password: user.Password,
		})
	}

	return domainUsers, nil
}
