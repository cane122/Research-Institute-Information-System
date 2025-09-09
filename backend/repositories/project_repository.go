package repositories

import (
"database/sql"
"github.com/cane/research-institute-system/backend/models"
)

type ProjectRepository struct {
db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
return nil
}

func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
return &models.Project{}, nil
}

func (r *ProjectRepository) GetByUserID(userID int) ([]models.Project, error) {
return []models.Project{}, nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
return nil  
}

func (r *ProjectRepository) GetMembers(projectID int) ([]models.User, error) {
return []models.User{}, nil
}
