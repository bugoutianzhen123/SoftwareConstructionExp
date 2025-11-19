package dao

import "github.com/bugoutianzhen123/SoftwareConstructionExp/domain"

type DAO interface {
    AddUser(u *domain.User) (*domain.User, error)
    ListUsers() []*domain.User
    GetUser(id int64) *domain.User
    GetUserByEmail(email string) *domain.User
    AddProject(p *domain.Project) (*domain.Project, error)
    ListProjects() []*domain.Project
    GetProject(id int64) *domain.Project
    AddApplication(a *domain.Application) (*domain.Application, error)
    ListApplications() []*domain.Application
    GetApplication(id int64) *domain.Application
    UpdateApplicationStatus(id int64, status string) error
    AddTracking(t *domain.Tracking) (*domain.Tracking, error)
    AddFeedback(f *domain.Feedback) (*domain.Feedback, error)
    ListFeedbacks() []*domain.Feedback
    AddDocument(d *domain.Document) (*domain.Document, error)
    ListDocumentsByApplication(appID int64) []*domain.Document
}