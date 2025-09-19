package post

import (
	"net/http"
	"strconv"
	"time"

	"news-feed/internal/middleware"
	"news-feed/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

type createPostRequest struct {
	Content string `json:"content" validate:"required,max=200"`
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	var req createPostRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.BadRequest("invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return err
	}

	userID := c.Locals("user_id").(float64) // jwt float64
	post, err := h.service.Create(uint(userID), req.Content)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(post)
}

func (h *Handler) GetFeed(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	// ambil pagination
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor", "")

	// ambil user yang di-follow + dirinya sendiri
	following := []uint{uint(userID)}
	// TODO: integrate dengan follow repository (sementara dummy dulu)

	posts, err := h.service.GetFeed(following, cursor, limit)
	if err != nil {
		return err
	}

	var nextCursor string
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		nextCursor = last.CreatedAt.Format(time.RFC3339)
	}

	return c.JSON(fiber.Map{
		"posts":       posts,
		"next_cursor": nextCursor,
	})
}
