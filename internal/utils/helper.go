package utils

import (
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
)

func ToModelUser(user entities.User) models.User {
	return models.User{
		Login:    user.Login,
		Password: user.Password,
		ID:       user.ID,
	}

}

func ToEntitiesUser(user models.User) entities.User {
	return entities.User{
		Login:    user.Login,
		Password: user.Password,
		ID:       user.ID,
	}

}

func ToModelCard(card entities.Card) models.Card {
	return models.Card{
		ID:        card.ID,
		Number:    card.Number,
		Cvv:       card.Cvv,
		ExpiredAt: card.ExpiredAt,
		UserId:    card.UserId,
		Version:   card.Version,
	}

}

func ToEntitiesCard(card models.Card) entities.Card {
	return entities.Card{
		ID:        card.ID,
		UserId:    card.UserId,
		Version:   card.Version,
		Cvv:       card.Cvv,
		ExpiredAt: card.ExpiredAt,
		Number:    card.Number,
	}

}
