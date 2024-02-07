package services

import (
	"github.com/patrickmn/go-cache"
	"github.com/zelas91/goph-keeper/internal/logger"
	compress2 "github.com/zelas91/goph-keeper/internal/utils/compress"
	"time"
)

type Service struct {
	Auth       *auth
	CreditCard *creditCard
	Credential *credential
	TextData   *textData
	BinaryFile *binaryFile
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

func WithCardUseRepository(cr cardRepo) func(s *Service) {
	return func(s *Service) {
		s.CreditCard = &creditCard{repo: cr}
	}
}

func WithCredentialUseRepository(cr credentialRepo) func(s *Service) {
	return func(s *Service) {
		s.Credential = &credential{repo: cr}
	}
}
func WithTextUseRepository(tr textDataRepo) func(s *Service) {
	return func(s *Service) {
		s.TextData = &textData{repo: tr}
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
