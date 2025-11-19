package service

import (
    "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/repository"
)

type Service struct {
    repo    repository.Repo
    matcher domain.Matcher
}

func NewService(r repository.Repo, m domain.Matcher) *Service {
    return &Service{repo: r, matcher: m}
}

func (s *Service) Repo() repository.Repo { return s.repo }
