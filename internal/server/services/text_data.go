package services

import (
	"fmt"

	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type textData struct {
	repo   textDataRepo
	crypto crypto
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
	textEntities, err := t.encryptToEntities(text)
	if err != nil {
		return err
	}
	if err := t.repo.Create(ctx, textEntities); err != nil {
		return fmt.Errorf("create text err: %w", err)
	}
	return nil
}

func (t textData) Texts(ctx context.Context) ([]models.TextData, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	texts, err := t.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get texts err: %w", err)
	}
	textsModel := make([]models.TextData, len(texts))
	for i, v := range texts {
		textModel, err := t.decryptToModels(v)
		if err != nil {
			return nil, err
		}
		textsModel[i] = textModel
	}
	return textsModel, err
}

func (t textData) Text(ctx context.Context, ucID int) (models.TextData, error) {
	userID := ctx.Value(types.UserIDKey).(int)

	text, err := t.repo.FindByIDAndUserID(ctx, ucID, userID)
	if err != nil {
		return models.TextData{}, fmt.Errorf("get text err: %w", err)
	}
	textModel, err := t.decryptToModels(text)
	if err != nil {
		return models.TextData{}, err
	}
	return textModel, nil
}

func (t textData) Delete(ctx context.Context, ucID int) error {
	userID := ctx.Value(types.UserIDKey).(int)

	if err := t.repo.Delete(ctx, ucID, userID); err != nil {
		return fmt.Errorf("delete text err: %w", err)
	}
	return nil
}

func (t textData) Update(ctx context.Context, text models.TextData) error {
	userID := ctx.Value(types.UserIDKey).(int)
	text.UserId = userID
	textEntities, err := t.encryptToEntities(text)
	if err != nil {
		return err
	}
	if err := t.repo.Update(ctx, textEntities); err != nil {
		return fmt.Errorf("update credentials err:%w", err)
	}
	return nil
}

func (t textData) encryptToEntities(td models.TextData) (entities.TextData, error) {
	text, err := t.crypto.Encrypt([]byte(td.Text))
	if err != nil {
		return entities.TextData{}, err
	}
	return entities.TextData{
		Text:    text,
		ID:      td.ID,
		UserId:  td.UserId,
		Version: td.Version,
	}, nil
}

func (t textData) decryptToModels(td entities.TextData) (models.TextData, error) {
	text, err := t.crypto.Decrypt(td.Text)
	if err != nil {
		return models.TextData{}, err
	}
	return models.TextData{
		ID:      td.ID,
		UserId:  td.UserId,
		Version: td.Version,
		Text:    string(text),
	}, nil
}
