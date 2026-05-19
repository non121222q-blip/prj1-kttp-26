package requests

import (
	"net/http"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Deadline    *int64 `json:"deadline"`
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {
	var timeUnix int64
	if r.Deadline != nil {
		timeUnix = *r.Deadline
	}
	var deadline *time.Time
	if timeUnix != 0 {
		d1 := time.Unix(timeUnix, 0)
		deadline = &d1
	}
	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Deadline:    deadline,
	}, nil
}

// Форма для зміни тільки статусу
type ChangeTaskStatusRequest struct {
	Status string `json:"status"`
}

// Bind потрібен, щоб Роутер зміг прочитати цей JSON
func (r *ChangeTaskStatusRequest) Bind(req *http.Request) error {
	return nil
}
func (r *ChangeTaskStatusRequest) ToDomainModel() (interface{}, error) {
	return nil, nil
}
