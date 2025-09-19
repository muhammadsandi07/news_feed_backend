package post

import (
	"net/http"
	"strconv"
	"time"

	"news-feed/internal/domain/follow"
	"news-feed/internal/middleware"
	"news-feed/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service       Service
	followService follow.Service
}

func NewHandler(s Service, f follow.Service) *Handler {
	return &Handler{s, f}
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

	userID := c.Locals("user_id").(string)
	post, err := h.service.Create(string(userID), req.Content)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(post)
}

func (h *Handler) GetFeed(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor := c.Query("cursor", "")

	following, err := h.followService.GetFollowingIDs(string(userID))
	if err != nil {
		return err
	}
	following = append(following, string(userID))
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
