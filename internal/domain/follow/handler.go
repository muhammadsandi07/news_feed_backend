package follow

import (
	"news-feed/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Follow(c *fiber.Ctx) error {
	userID := uint(c.Locals("user_id").(float64))
	targetID, err := c.ParamsInt("id")
	if err != nil {
		return middleware.BadRequest("invalid user id")
	}

	if err := h.service.FollowUser(userID, uint(targetID)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "you are now following user " + c.Params("id")})
}

func (h *Handler) Unfollow(c *fiber.Ctx) error {
	userID := uint(c.Locals("user_id").(float64))
	targetID, err := c.ParamsInt("id")
	if err != nil {
		return middleware.BadRequest("invalid user id")
	}

	if err := h.service.UnfollowUser(userID, uint(targetID)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "you unfollowed user " + c.Params("id")})
}
