package services

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type creditCard struct {
	repo cardRepo
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
	if err := c.repo.Create(ctx, helper.ToEntitiesCard(card)); err != nil {
		return fmt.Errorf("create card err: %v", err)
	}
	return nil
}

func (c creditCard) Cards(ctx context.Context) ([]models.Card, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	cards, err := c.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cards err: %v", err)
	}
	cardsModel := make([]models.Card, len(cards))
	for i, v := range cards {
		cardsModel[i] = helper.ToModelCard(v)
	}
	return cardsModel, err
}

func (c creditCard) Card(ctx context.Context, cardID int) (models.Card, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	card, err := c.repo.FindByIDAndUserID(ctx, cardID, userID)
	if err != nil {
		return models.Card{}, fmt.Errorf("get card err: %v", err)
	}
	return helper.ToModelCard(card), nil
}

func (c creditCard) Delete(ctx context.Context, cardID int) error {
	userID := ctx.Value(types.UserIDKey).(int)

	if err := c.repo.Delete(ctx, cardID, userID); err != nil {
		return fmt.Errorf("delete card err: %v", err)
	}
	return nil
}

func (c creditCard) Update(ctx context.Context, card models.Card) error {
	userID := ctx.Value(types.UserIDKey).(int)
	card.UserId = userID

	if err := c.repo.Update(ctx, helper.ToEntitiesCard(card)); err != nil {
		return fmt.Errorf("update card err:%v", err)
	}
	return nil
}
