package handlers

import (
	"go-fiber-crud/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type TodoHandler struct {
	TodoRepo *repository.TodoRepository
}

func NewTodoHandler(repository *repository.TodoRepository) *TodoHandler {
	return &TodoHandler{TodoRepo: repository}
}

func (h *TodoHandler) ShowTodos(c fiber.Ctx, store *session.Store) error {
	sess, _ := store.Get(c)
	userId := sess.Get("user_id").(int)
	todos, err := h.TodoRepo.GetTodos(strconv.Itoa(userId))
	if err != nil {
		return c.Status(400).Render("index", fiber.Map{
			"Error":     err.Error(),
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}

	return c.Status(200).Render("index", fiber.Map{
		"Todos":     todos,
		"CSRFToken": csrf.TokenFromContext(c),
	})
}

func (h *TodoHandler) CreateTodo(c fiber.Ctx, store *session.Store) error {
	sess, _ := store.Get(c)
	userId := sess.Get("user_id").(int)
	title := c.FormValue("title")

	err := h.TodoRepo.CreateTodo(strconv.Itoa(userId), title)
	if err != nil {
		return c.Status(400).Render("index", fiber.Map{
			"Error":     "Failed add todo :" + err.Error(),
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}

	return c.Redirect().To("/")
}

func (h *TodoHandler) DeleteTodo(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.TodoRepo.DeleteTodo(id)
	if err != nil {
		return c.Status(500).Render("index", fiber.Map{
			"Error":     "Failed delete todo : " + err.Error(),
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}
	return c.Redirect().To("/")
}
