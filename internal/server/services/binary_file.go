package services

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
	"io"
	"os"
	"path/filepath"
)

type binaryFile struct {
	log      logger.Logger
	repo     binaryFileRepo
	basePath string
}

//go:generate mockgen -package mocks -destination=./mocks/mock_binary_file_repo.go -source=binary_file.go -package=mock
type binaryFileRepo interface {
	Create(ctx context.Context, bf entities.BinaryFile) error
	FindFileByUserID(ctx context.Context, fileID, userID int) (entities.BinaryFile, error)
	FindAllByUserID(ctx context.Context, userID int) ([]entities.BinaryFile, error)
	Delete(ctx context.Context, userID, fileID int) error
}

func (b binaryFile) Upload(ctx context.Context, bf models.BinaryFile, reader <-chan []byte) error {
	userID := ctx.Value(types.UserIDKey).(int)
	path := fmt.Sprintf("%s/%d/%s", b.basePath, userID, bf.FileName)
	bf.Path = path
	bf.UserId = userID

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file:%v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			b.log.Errorf("close file err:%v", err)
		}
	}()
	for val := range reader {
		_, err = file.Write(val)
		if err != nil {
			return fmt.Errorf("save err %v", err)
		}
	}
	if err := b.repo.Create(ctx, helper.ToEntitiesBinaryFile(bf)); err != nil {
		return err
	}
	return nil
}

func (b binaryFile) Download(ctx context.Context, bf models.BinaryFile, write chan<- []byte) error {
	userID := ctx.Value(types.UserIDKey).(int)
	ef, err := b.repo.FindFileByUserID(ctx, bf.ID, userID)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(ef.Path, os.O_RDONLY, 066)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			b.log.Error("download: close file err: %v", err)
		}
	}()
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				b.log.Error("failed read file err: %v", err)
				return err
			}
			break
		}
		write <- buffer[:n]
	}
	return nil
}

func (b binaryFile) Files(ctx context.Context) ([]models.BinaryFile, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	ef, err := b.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	files := make([]models.BinaryFile, len(ef))
	for i, v := range ef {
		files[i] = helper.ToModelBinaryFile(v)
	}
	return files, nil
}
func (b binaryFile) Delete(ctx context.Context, fileID int) error {
	userID := ctx.Value(types.UserIDKey).(int)
	file, err := b.repo.FindFileByUserID(ctx, fileID, userID)
	if err != nil {
		return err
	}

	if err = os.Remove(file.Path); err != nil {
		return err
	}
	return b.repo.Delete(ctx, userID, fileID)
}
