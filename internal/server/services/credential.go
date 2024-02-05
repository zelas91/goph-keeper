package services

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type credential struct {
	repo credentialRepo
}

//go:generate mockgen -package mocks -destination=./mocks/mock_credential_repo.go -source=credential.go -package=mock
type credentialRepo interface {
	Create(ctx context.Context, uc entities.UserCredentials) error
	FindAllByUserID(ctx context.Context, userID int) ([]entities.UserCredentials, error)
	FindCredentialByUserID(ctx context.Context, ucID, userID int) (entities.UserCredentials, error)
	Delete(ctx context.Context, ucID, userID int) error
	Update(ctx context.Context, uc entities.UserCredentials) error
}

func (c credential) Create(ctx context.Context, uc models.UserCredentials) error {
	userID := ctx.Value(types.UserIDKey).(int)
	uc.UserId = userID
	if err := c.repo.Create(ctx, helper.ToEntitiesUserCredential(uc)); err != nil {
		return fmt.Errorf("create credential err: %v", err)
	}
	return nil
}

func (c credential) Credentials(ctx context.Context) ([]models.UserCredentials, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	ucs, err := c.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get credentials err: %v", err)
	}
	ucsModel := make([]models.UserCredentials, len(ucs))
	for i, v := range ucs {
		ucsModel[i] = helper.ToModelUserCredential(v)
	}
	return ucsModel, err
}

func (c credential) Credential(ctx context.Context, ucID int) (models.UserCredentials, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	uc, err := c.repo.FindCredentialByUserID(ctx, ucID, userID)
	if err != nil {
		return models.UserCredentials{}, fmt.Errorf("get Credentials err: %v", err)
	}
	return helper.ToModelUserCredential(uc), nil
}

func (c credential) Delete(ctx context.Context, ucID int) error {
	userID := ctx.Value(types.UserIDKey).(int)

	if err := c.repo.Delete(ctx, ucID, userID); err != nil {
		return fmt.Errorf("delete credential err: %v", err)
	}
	return nil
}

func (c credential) Update(ctx context.Context, uc models.UserCredentials) error {
	userID := ctx.Value(types.UserIDKey).(int)
	uc.UserId = userID
	if err := c.repo.Update(ctx, helper.ToEntitiesUserCredential(uc)); err != nil {
		return fmt.Errorf("update credentials err:%v", err)
	}
	return nil
}
