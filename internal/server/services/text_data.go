package services

import (
	"fmt"

	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type textData struct {
	repo textDataRepo
}

//go:generate mockgen -package mocks -destination=./mocks/mock_text_data_repo.go -source=text_data.go -package=mock
type textDataRepo interface {
	Create(ctx context.Context, text entities.TextData) error
	FindAllByUserID(ctx context.Context, textID int) ([]entities.TextData, error)
	FindByIDAndUserID(ctx context.Context, textID, userID int) (entities.TextData, error)
	Delete(ctx context.Context, textID, userID int) error
	Update(ctx context.Context, uc entities.TextData) error
}

func (t textData) Create(ctx context.Context, text models.TextData) error {
	userID := ctx.Value(types.UserIDKey).(int)
	text.UserId = userID
	if err := t.repo.Create(ctx, helper.ToEntitiesText(text)); err != nil {
		return fmt.Errorf("create text err: %v", err)
	}
	return nil
}

func (t textData) Texts(ctx context.Context) ([]models.TextData, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	texts, err := t.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get texts err: %v", err)
	}
	textsModel := make([]models.TextData, len(texts))
	for i, v := range texts {
		textsModel[i] = helper.ToModelText(v)
	}
	return textsModel, err
}

func (t textData) Text(ctx context.Context, ucID int) (models.TextData, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	text, err := t.repo.FindByIDAndUserID(ctx, ucID, userID)
	if err != nil {
		return models.TextData{}, fmt.Errorf("get text err: %v", err)
	}
	return helper.ToModelText(text), nil
}

func (t textData) Delete(ctx context.Context, ucID int) error {
	userID := ctx.Value(types.UserIDKey).(int)

	if err := t.repo.Delete(ctx, ucID, userID); err != nil {
		return fmt.Errorf("delete text err: %v", err)
	}
	return nil
}

func (t textData) Update(ctx context.Context, text models.TextData) error {
	userID := ctx.Value(types.UserIDKey).(int)
	text.UserId = userID
	if err := t.repo.Update(ctx, helper.ToEntitiesText(text)); err != nil {
		return fmt.Errorf("update credentials err:%v", err)
	}
	return nil
}
