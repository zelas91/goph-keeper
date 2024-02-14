package services

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/zelas91/goph-keeper/internal/logger"
	compress2 "github.com/zelas91/goph-keeper/internal/utils/compress"
)

type Service struct {
	Auth       *auth
	CreditCard *creditCard
	Credential *credential
	TextData   *textData
	BinaryFile *binaryFile
}

type crypto interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

func New(options ...func(s *Service)) *Service {
	sv := &Service{}
	for _, opt := range options {
		opt(sv)
	}
	return sv
}

func WithAuthUseRepository(up userRepo) func(s *Service) {
	return func(s *Service) {
		s.Auth = &auth{repo: up, cache: cache.New(time.Minute*10, time.Minute*10)}
	}
}

func WithCardUseRepository(cr cardRepo, crypto crypto) func(s *Service) {
	return func(s *Service) {
		s.CreditCard = &creditCard{repo: cr, crypto: crypto}
	}
}

func WithCredentialUseRepository(cr credentialRepo, crypto crypto) func(s *Service) {
	return func(s *Service) {
		s.Credential = &credential{repo: cr, crypto: crypto}
	}
}
func WithTextUseRepository(tr textDataRepo, crypto crypto) func(s *Service) {
	return func(s *Service) {
		s.TextData = &textData{repo: tr, crypto: crypto}
	}
}

func WithBinaryFileUseRepository(bf binaryFileRepo, log logger.Logger, basePathSaveFile string) func(s *Service) {
	return func(s *Service) {
		s.BinaryFile = &binaryFile{
			repo:       bf,
			log:        log,
			basePath:   basePathSaveFile,
			compress:   compress2.NewCompress(log),
			decompress: compress2.NewDecompress(),
		}
	}
}
