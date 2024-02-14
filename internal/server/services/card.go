package services

import (
	"fmt"

	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type creditCard struct {
	repo   cardRepo
	crypto crypto
}

//go:generate mockgen -package mocks -destination=./mocks/mock_card_repo.go -source=card.go -package=mock
type cardRepo interface {
	Create(ctx context.Context, card entities.Card) error
	FindAllByUserID(ctx context.Context, userID int) ([]entities.Card, error)
	FindByIDAndUserID(ctx context.Context, cardID, userID int) (entities.Card, error)
	Delete(ctx context.Context, cardID, userID int) error
	Update(ctx context.Context, card entities.Card) error
}

func (c creditCard) Create(ctx context.Context, card models.Card) error {
	userID := ctx.Value(types.UserIDKey).(int)
	card.UserId = userID
	entitiesCard, err := c.encryptToEntities(card)
	if err != nil {
		return err
	}
	if err := c.repo.Create(ctx, entitiesCard); err != nil {
		return fmt.Errorf("create card err: %w", err)
	}
	return nil
}

func (c creditCard) Cards(ctx context.Context) ([]models.Card, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	cards, err := c.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cards err: %w", err)
	}
	cardsModel := make([]models.Card, len(cards))
	for i, v := range cards {
		modelCard, err := c.decryptToModels(v)
		if err != nil {
			return nil, err
		}
		cardsModel[i] = modelCard
	}
	return cardsModel, err
}

func (c creditCard) Card(ctx context.Context, cardID int) (models.Card, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	card, err := c.repo.FindByIDAndUserID(ctx, cardID, userID)
	if err != nil {
		return models.Card{}, fmt.Errorf("get card err: %w", err)
	}
	modelCard, err := c.decryptToModels(card)
	if err != nil {
		return models.Card{}, err
	}
	return modelCard, nil
}

func (c creditCard) Delete(ctx context.Context, cardID int) error {
	userID := ctx.Value(types.UserIDKey).(int)

	if err := c.repo.Delete(ctx, cardID, userID); err != nil {
		return fmt.Errorf("delete card err: %w", err)
	}
	return nil
}

func (c creditCard) Update(ctx context.Context, card models.Card) error {
	userID := ctx.Value(types.UserIDKey).(int)
	card.UserId = userID
	entitiesCard, err := c.encryptToEntities(card)
	if err != nil {
		return err
	}
	if err := c.repo.Update(ctx, entitiesCard); err != nil {
		return fmt.Errorf("update card err:%w", err)
	}
	return nil
}
func (c creditCard) encryptToEntities(card models.Card) (entities.Card, error) {
	number, err := c.crypto.Encrypt([]byte(card.Number))
	if err != nil {
		return entities.Card{}, err
	}
	cvv, err := c.crypto.Encrypt([]byte(card.Cvv))
	if err != nil {
		return entities.Card{}, err
	}
	ex, err := c.crypto.Encrypt([]byte(card.ExpiredAt))
	if err != nil {
		return entities.Card{}, err
	}
	return entities.Card{
		Number:    number,
		Cvv:       cvv,
		ExpiredAt: ex,
		ID:        card.ID,
		UserId:    card.UserId,
		Version:   card.Version,
	}, nil
}

func (c creditCard) decryptToModels(card entities.Card) (models.Card, error) {
	number, err := c.crypto.Decrypt(card.Number)
	if err != nil {
		return models.Card{}, err
	}
	cvv, err := c.crypto.Decrypt(card.Cvv)
	if err != nil {
		return models.Card{}, err
	}
	ex, err := c.crypto.Decrypt(card.ExpiredAt)
	if err != nil {
		return models.Card{}, err
	}
	return models.Card{
		Number:    string(number),
		Cvv:       string(cvv),
		ExpiredAt: string(ex),
		ID:        card.ID,
		UserId:    card.UserId,
		Version:   card.Version,
	}, nil
}
