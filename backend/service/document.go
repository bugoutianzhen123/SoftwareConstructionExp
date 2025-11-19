package service

import "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"

func (s *Service) AddDocument(d *domain.Document) (*domain.Document, error) { return s.repo.AddDocument(d) }
func (s *Service) ListDocumentsByApplication(appID int64) []*domain.Document { return s.repo.ListDocumentsByApplication(appID) }