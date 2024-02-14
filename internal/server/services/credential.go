package services

import (
	"fmt"

	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type credential struct {
	repo   credentialRepo
	crypto crypto
}

//go:generate mockgen -package mocks -destination=./mocks/mock_credential_repo.go -source=credential.go -package=mock
type credentialRepo interface {
	Create(ctx context.Context, uc entities.UserCredentials) error
	FindAllByUserID(ctx context.Context, userID int) ([]entities.UserCredentials, error)
	FindByIDAndUserID(ctx context.Context, ucID, userID int) (entities.UserCredentials, error)
	Delete(ctx context.Context, ucID, userID int) error
	Update(ctx context.Context, uc entities.UserCredentials) error
}

func (c credential) Create(ctx context.Context, uc models.UserCredentials) error {
	userID := ctx.Value(types.UserIDKey).(int)
	uc.UserId = userID
	ucEntities, err := c.encryptToEntities(uc)
	if err != nil {
		return err
	}
	if err := c.repo.Create(ctx, ucEntities); err != nil {
		return fmt.Errorf("create credential err: %w", err)
	}
	return nil
}

func (c credential) Credentials(ctx context.Context) ([]models.UserCredentials, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	ucs, err := c.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get credentials err: %w", err)
	}
	ucsModel := make([]models.UserCredentials, len(ucs))
	for i, v := range ucs {
		ucModel, err := c.decryptToModels(v)
		if err != nil {
			return nil, err
		}
		ucsModel[i] = ucModel
	}
	return ucsModel, err
}

func (c credential) Credential(ctx context.Context, ucID int) (models.UserCredentials, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	uc, err := c.repo.FindByIDAndUserID(ctx, ucID, userID)

	if err != nil {
		return models.UserCredentials{}, fmt.Errorf("get Credentials err: %w", err)
	}
	ucModel, err := c.decryptToModels(uc)
	if err != nil {
		return models.UserCredentials{}, err
	}
	return ucModel, nil
}

func (c credential) Delete(ctx context.Context, ucID int) error {
	userID := ctx.Value(types.UserIDKey).(int)

	if err := c.repo.Delete(ctx, ucID, userID); err != nil {
		return fmt.Errorf("delete credential err: %w", err)
	}
	return nil
}

func (c credential) Update(ctx context.Context, uc models.UserCredentials) error {
	userID := ctx.Value(types.UserIDKey).(int)
	uc.UserId = userID
	ucEntities, err := c.encryptToEntities(uc)
	if err != nil {
		return err
	}
	if err := c.repo.Update(ctx, ucEntities); err != nil {
		return fmt.Errorf("update credentials err:%w", err)
	}
	return nil
}

func (c credential) encryptToEntities(us models.UserCredentials) (entities.UserCredentials, error) {
	password, err := c.crypto.Encrypt([]byte(us.Password))
	if err != nil {
		return entities.UserCredentials{}, err
	}
	login, err := c.crypto.Encrypt([]byte(us.Login))
	if err != nil {
		return entities.UserCredentials{}, err
	}
	return entities.UserCredentials{
		Login:    login,
		Password: password,
		ID:       us.ID,
		UserId:   us.UserId,
		Version:  us.Version,
	}, nil
}

func (c credential) decryptToModels(us entities.UserCredentials) (models.UserCredentials, error) {
	login, err := c.crypto.Decrypt(us.Login)
	if err != nil {
		return models.UserCredentials{}, err
	}
	password, err := c.crypto.Decrypt(us.Password)
	if err != nil {
		return models.UserCredentials{}, err
	}
	return models.UserCredentials{
		ID:       us.ID,
		UserId:   us.UserId,
		Version:  us.Version,
		Password: string(password),
		Login:    string(login),
	}, nil
}
