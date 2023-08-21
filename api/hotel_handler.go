package api

import (
	"github.com/Jimbo8702/goreservation/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := db.Map{"hotelID": id}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("room")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}


func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return ErrResourceNotFound("hotels")
	}
	return c.JSON(hotels)
}

// func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error 
// func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error 
// func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error 
// func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error 
