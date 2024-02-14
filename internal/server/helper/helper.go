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

func ToModelBinaryFile(b entities.BinaryFile) models.BinaryFile {
	return models.BinaryFile{
		ID:       b.ID,
		UserId:   b.UserId,
		FileName: b.FileName,
		Path:     b.Path,
		Size:     b.Size,
	}

}

func ToEntitiesBinaryFile(b models.BinaryFile) entities.BinaryFile {
	return entities.BinaryFile{
		ID:       b.ID,
		UserId:   b.UserId,
		Path:     b.Path,
		FileName: b.FileName,
		Size:     b.Size,
	}

}
