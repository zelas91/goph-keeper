package services

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
)

type binaryFile struct {
	log        logger.Logger
	repo       binaryFileRepo
	basePath   string
	compress   compress
	decompress decompress
}

type compress interface {
	Writer() *gzip.Writer
	Release(writer *gzip.Writer)
}

type decompress interface {
	Release(reader *gzip.Reader)
	Reader() *gzip.Reader
}

//go:generate mockgen -package mocks -destination=./mocks/mock_binary_file_repo.go -source=binary_file.go -package=mock
type binaryFileRepo interface {
	Create(ctx context.Context, bf entities.BinaryFile) error
	FindByIDAndUserID(ctx context.Context, fileID, userID int) (entities.BinaryFile, error)
	FindAllByUserID(ctx context.Context, userID int) ([]entities.BinaryFile, error)
	Delete(ctx context.Context, userID, fileID int) error
}

func (b *binaryFile) Upload(ctx context.Context, bf models.BinaryFile, reader <-chan []byte) error {
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
	gz := b.compress.Writer()
	defer b.compress.Release(gz)
	gz.Reset(file)
	count := 0
	for val := range reader {
		count = count + len(val)
		_, err = gz.Write(val)
		if err != nil {
			return fmt.Errorf("save err %v", err)
		}
	}
	if err = gz.Flush(); err != nil {
		return fmt.Errorf("gzip flush err %v", err)
	}

	if bf.Size != count {
		return errors.New("file does not match length")
	}

	if err := b.repo.Create(ctx, helper.ToEntitiesBinaryFile(bf)); err != nil {
		return err
	}
	return nil
}

func (b *binaryFile) Download(ctx context.Context, bf models.BinaryFile, write chan<- []byte) error {
	defer close(write)
	userID := ctx.Value(types.UserIDKey).(int)
	ef, err := b.repo.FindByIDAndUserID(ctx, bf.ID, userID)
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

	gz := b.decompress.Reader()
	defer b.decompress.Release(gz)
	if err = gz.Reset(file); err != nil {
		return fmt.Errorf("gzip reader reset err: %v", err)
	}

	for {
		buffer := make([]byte, 1024)
		n, err := gz.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("failed read file err: %v", err)
			}
			break
		}
		write <- buffer[:n]
	}
	return nil
}

func (b *binaryFile) Files(ctx context.Context) ([]models.BinaryFile, error) {
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

func (b *binaryFile) File(ctx context.Context, fileID int) (models.BinaryFile, error) {
	userID := ctx.Value(types.UserIDKey).(int)
	file, err := b.repo.FindByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return models.BinaryFile{}, fmt.Errorf("get file err: %v", err)
	}
	return helper.ToModelBinaryFile(file), nil
}
func (b *binaryFile) Delete(ctx context.Context, fileID int) error {
	userID := ctx.Value(types.UserIDKey).(int)
	file, err := b.repo.FindByIDAndUserID(ctx, fileID, userID)
	if err != nil {
		return err
	}

	if err = os.Remove(file.Path); err != nil {
		return err
	}
	return b.repo.Delete(ctx, userID, fileID)
}
