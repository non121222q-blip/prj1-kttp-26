package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}
func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Save(request.Bind): %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		task.UserId = user.Id
		task.Status = domain.NewTaskStatus

		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController.Save(c.taskService.Save): %s", err)
			InternalServerError(w, err)
			return

		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)

		Success(w, taskDto)
	}
}
func (c TaskController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		// 1. Створюємо пусту "коробку" для фільтрів
		var filter domain.TaskFilter

		// 2. Читаємо побажання клієнта з посилання
		if statusStr := r.URL.Query().Get("status"); statusStr != "" {
			status := domain.TaskStatus(statusStr)
			filter.Status = &status
		}

		filter.SortBy = r.URL.Query().Get("sortBy")
		filter.SortOrder = r.URL.Query().Get("sortOrder")

		if fromStr := r.URL.Query().Get("deadlineFrom"); fromStr != "" {
			if parsedFrom, err := time.Parse(time.RFC3339, fromStr); err == nil {
				filter.DeadlineFrom = &parsedFrom
			}
		}
		if toStr := r.URL.Query().Get("deadlineTo"); toStr != "" {
			if parsedTo, err := time.Parse(time.RFC3339, toStr); err == nil {
				filter.DeadlineTo = &parsedTo
			}
		}

		// 3. Віддаємо зібраний фільтр Менеджеру!
		tasks, err := c.taskService.FindList(user.Id, filter)
		if err != nil {
			log.Printf("TaskController.FindList(c.taskService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		tasksDto := taskDto.DomainToDtoCollection(tasks)

		Success(w, tasksDto)
	}
}
func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)

		Success(w, taskDto)
	}
}

//todo: add method for chang (update) Task status

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		updTask, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Save(requests.Bind): %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		task.Title = updTask.Title
		task.Description = updTask.Description
		task.Deadline = updTask.Deadline

		task, err = c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController.Update(c.taskService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)

		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))

			return
		}
		err := c.taskService.Delete(task.Id)
		if err != nil {
			log.Printf("TaskController.Delete(c.taskService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (c TaskController) ChangeStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Дістаємо таску, яку хочемо змінити (помічник tpom вже поклав її сюди)
		task := r.Context().Value(TaskKey).(domain.Task)

		// 2. Читаємо новий статус, який прислав клієнт у JSON
		var req requests.ChangeTaskStatusRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("TaskController.ChangeStatus(requests.Bind): %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		// 3. Змінюємо статус у нашій тасці
		task.Status = domain.TaskStatus(req.Status)

		// 4. Віддаємо Менеджеру, щоб він зберіг зміни в базу
		updatedTask, err := c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController.ChangeStatus(c.taskService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		// 5. Віддаємо успішну відповідь
		taskDto := resources.TaskDto{}
		Success(w, taskDto.DomainToDto(updatedTask))
	}
}
