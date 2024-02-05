package helper

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

func ToModelUserCredential(uc entities.UserCredentials) models.UserCredentials {
	return models.UserCredentials{
		ID:       uc.ID,
		UserId:   uc.UserId,
		Version:  uc.Version,
		Password: uc.Password,
		Login:    uc.Login,
	}

}

func ToEntitiesUserCredential(uc models.UserCredentials) entities.UserCredentials {
	return entities.UserCredentials{
		ID:       uc.ID,
		UserId:   uc.UserId,
		Version:  uc.Version,
		Password: uc.Password,
		Login:    uc.Login,
	}

}

func ToModelText(t entities.TextData) models.TextData {
	return models.TextData{
		ID:      t.ID,
		UserId:  t.UserId,
		Version: t.Version,
		Text:    t.Text,
	}

}

func ToEntitiesText(t models.TextData) entities.TextData {
	return entities.TextData{
		ID:      t.ID,
		UserId:  t.UserId,
		Version: t.Version,
		Text:    t.Text,
	}

}
