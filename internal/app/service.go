package app

import "database/sql"

type service struct {
	repo Repo
}

type Service interface {
	Stats() sql.DBStats
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Stats() sql.DBStats {
	return s.repo.Stats()
}
