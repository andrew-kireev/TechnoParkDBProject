package thread

import "TechnoParkDBProject/internal/app/thread/models"

type Usecase interface {
	CreateThread(thread *models.Thread) error
}
