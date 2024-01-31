package utils

import (
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
)

func EntitiesUserInModelUser(user entities.User) models.User {
	return models.User{
		Login:    user.Login,
		Password: user.Password,
		Id:       user.Id,
	}

}

func ModelUserInEntitiesUser(user models.User) entities.User {
	return entities.User{
		Login:    user.Login,
		Password: user.Password,
		Id:       user.Id,
	}

}
