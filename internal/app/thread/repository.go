package thread

import "TechnoParkDBProject/internal/app/thread/models"

type Repository interface {
	CreateThread(thread *models.Thread) error
}
