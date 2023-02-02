package repository

import (
	"database/sql"
	"github.com/aidanb22/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllPlayers() ([]*models.Player, error)
	AllUsers() ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	OnePlayer(id int) (*models.Player, error)
	OnePlayerForEdit(id int) (*models.Player, error)
	InsertPlayer(player models.Player) (int, error)
	UpdatePlayer(player models.Player) error
	DeletePlayer(id int) error
}
