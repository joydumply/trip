package app

import (
	"fmt"
	"strconv"
	"trip-planner/internal/model"
	"trip-planner/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/trips", h.CreateTrip)
	app.Get("/trips", h.GetTrips)
	app.Get("/trips/:id", h.GetTripByID)
	app.Put("/trips/:id", h.UpdateTrip)
	app.Delete("/trips/:id", h.DeleteTrip)
}

func (h *Handler) CreateTrip(c *fiber.Ctx) error {
	var trip model.Trip
	if err := c.BodyParser(&trip); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	err := h.service.CreateTrip(c.Context(), &trip)
	if err != nil {
		fmt.Println("CreateTrip service error:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(trip)
}

func (h *Handler) GetTrips(c *fiber.Ctx) error {

	title := c.Query("title")
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")

	var leaderID *int64
	if id := c.Query("leader_id"); id != "" {
		parsed, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			leaderID = &parsed
		}
	}

	trips, err := h.service.GetTrips(c.Context(), title, leaderID, limit, offset, sortBy, order)
	if err != nil {
		fmt.Println("GetTrips handler error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get trips1"})
	}

	if trips == nil {
		trips = []model.Trip{}
	}

	return c.JSON(trips)
}

func (h *Handler) GetTripByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	trip, err := h.service.GetTripByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(trip)

}
func (h *Handler) UpdateTrip(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid trip ID"})
	}

	trip := new(model.Trip)
	if err := c.BodyParser(trip); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	trip.ID = id
	if err := h.service.UpdateTrip(c.Context(), trip); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Trip updated successfully"})
}

func (h *Handler) DeleteTrip(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid trip ID"})
	}

	err = h.service.DeleteTrip(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Trip deleted successfully"})
}
