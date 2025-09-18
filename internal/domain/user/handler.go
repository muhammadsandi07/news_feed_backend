package user

import (
	"net/http"
	"news-feed/internal/middleware"
	"news-feed/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.BadRequest("invalid request body")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return err
	}

	u, err := h.service.Register(req.Username, req.Password)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":       u.ID,
		"username": u.Username,
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.BadRequest("invalid request body")
	}

	access, refresh, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
