package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description string            `db:"description"`
	Status      domain.TaskStatus `db:"status"`
	Deadline    *time.Time        `db:"deadline"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	FindList(uId uint64, filter domain.TaskFilter) ([]domain.Task, error)
	Find(id uint64) (domain.Task, error)
	Update(t domain.Task) (domain.Task, error)
	Delete(id uint64) error
}

type taskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(session db.Session) TaskRepository {
	return taskRepository{
		coll: session.Collection(TasksTableName),
		sess: session,
	}
}
func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()

	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
}
func (r taskRepository) FindList(uId uint64, filter domain.TaskFilter) ([]domain.Task, error) {
	var tasks []task
	cond := db.Cond{"user_id": uId, "deleted_date": nil}
	if filter.Status != nil {
		cond["status"] = *filter.Status
	}
	if filter.DeadlineFrom != nil {
		cond["deadline >="] = *filter.DeadlineFrom
	}
	if filter.DeadlineTo != nil {
		cond["deadline <="] = *filter.DeadlineTo
	}

	// 3. Формуємо запит до бази
	query := r.coll.Find(cond)

	// 4. Додаємо сортування
	if filter.SortBy != "" {
		orderParam := filter.SortBy
		if filter.SortOrder == "desc" {
			orderParam = "-" + orderParam // Мінус для сортування за спаданням
		}
		query = query.OrderBy(orderParam)
	} else {
		// Сортування за замовчуванням
		query = query.OrderBy("-created_date")
	}

	// 5. Виконуємо запит
	err := query.All(&tasks)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(tasks), nil
}

func (r taskRepository) Find(id uint64) (domain.Task, error) {
	var t task

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&t)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(t), nil

}
func (r taskRepository) Update(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)

	tsk.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": t.Id, "deleted_date": nil}).Update(&tsk)

	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r taskRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	tasks := make([]domain.Task, len(ts))
	for i := range ts {
		tasks[i] = r.mapModelToDomain(ts[i])
	}
	return tasks
}
