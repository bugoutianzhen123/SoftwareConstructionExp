package service

import (
    "errors"
    "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"
)

func (s *Service) MatchForStudent(studentID int64) ([]domain.MatchResult, error) {
    stu := s.repo.GetUser(studentID)
    if stu == nil || stu.Role != domain.RoleStudent { return nil, errors.New("学生不存在") }
    ps := s.repo.ListProjects()
    return s.matcher.Match(stu, ps), nil
}