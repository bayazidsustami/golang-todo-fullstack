package main

import (
	"context"
	"database/sql"
	"go-fiber-crud/internal/config"
	"go-fiber-crud/internal/database"
	"go-fiber-crud/internal/handlers"
	"go-fiber-crud/internal/repository"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func main() {

	loggerHandler := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(loggerHandler)

	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		AppName:     "ToDo App",
		ReadTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New())

	store := session.NewStore(session.Config{
		IdleTimeout:    24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSecure:   false,
	})

	csrfMiddleware := csrf.New(csrf.Config{
		Session:   store,
		Extractor: extractors.FromForm("_csrf"),
	})

	app.Use(csrfMiddleware)

	setupRoutes(app, db, store)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("sever is starting", "port", cfg.ServerPort)
		if err := app.Listen(cfg.ServerPort); err != nil {
			slog.Error("server shutdown with error", "error", err)
		}
	}()

	<-stop
	slog.Info("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server cleanup finished")
}

func setupRoutes(app *fiber.App, db *sql.DB, store *session.Store) {

	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo)
	todoHandler := handlers.NewTodoHandler(todoRepo)

	authMiddleware := func(c fiber.Ctx) error {
		sess, _ := store.Get(c)
		userId := sess.Get("user_id")
		if userId == nil {
			return c.Redirect().To("/login")
		}
		return c.Next()
	}

	guestMiddleware := func(c fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Next()
		}
		if sess.Get("user_id") != nil {
			return c.Redirect().To("/")
		}

		return c.Next()
	}

	app.Get("/login", guestMiddleware, authHandler.ShowLogin)
	app.Get("/register", guestMiddleware, authHandler.ShowRegister)
	app.Post("/register", authHandler.RegisterUser)
	app.Post("/login", func(c fiber.Ctx) error { return authHandler.LoginUser(c, store) })
	app.Post("/logout", func(c fiber.Ctx) error { return authHandler.LogoutUser(c, store) })

	app.Get("/", authMiddleware, func(c fiber.Ctx) error { return todoHandler.ShowTodos(c, store) })
	app.Post("/todos", authMiddleware, func(c fiber.Ctx) error { return todoHandler.CreateTodo(c, store) })
	app.Post("/todos/delete/:id", authMiddleware, todoHandler.DeleteTodo)
}
