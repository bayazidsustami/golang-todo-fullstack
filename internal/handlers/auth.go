package handlers

import (
	"database/sql"
	"go-fiber-crud/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

func NewAuthHandler(repository *repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: repository}
}

func (h *AuthHandler) ShowRegister(c fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"CSRFToken": csrf.TokenFromContext(c),
	})
}

func (h *AuthHandler) ShowLogin(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"CSRFToken": csrf.TokenFromContext(c),
	})
}

func (h *AuthHandler) RegisterUser(c fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := h.UserRepo.Create(username, string(hashedPassword))
	if err != nil {
		return c.Status(400).Render("login", fiber.Map{
			"Error":     err.Error(),
			"Username":  username,
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}

	return c.Redirect().To("/login")
}

func (h *AuthHandler) LoginUser(c fiber.Ctx, store *session.Store) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.UserRepo.FindByUsername(username)

	if err == sql.ErrNoRows {
		return c.Status(401).Render("login", fiber.Map{
			"Error":     "user not found",
			"Username":  username,
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(400).Render("login", fiber.Map{
			"Error":     "wrong password",
			"Username":  username,
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}

	sess, _ := store.Get(c)
	sess.Set("user_id", user.ID)
	if err := sess.Save(); err != nil {
		return c.Status(500).Render("login", fiber.Map{
			"Error":     "failed to set session",
			"Username":  username,
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}

	return c.Redirect().To("/")
}

func (h *AuthHandler) LogoutUser(c fiber.Ctx, store *session.Store) error {
	sess, _ := store.Get(c)
	if err := sess.Destroy(); err != nil {
		return c.Status(500).Render("index", fiber.Map{
			"Error":     "failed to logout ",
			"CSRFToken": csrf.TokenFromContext(c),
		})
	}
	return c.Redirect().To("/login")
}
