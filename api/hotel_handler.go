package api

import (
	"github.com/Jimbo8702/goreservation/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	hotelStore 	db.HotelStore
	roomStore 	db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore: rs,
	}
}

type HotelQueryParams struct {
	Rooms bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return nil
	}

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

// func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error 
// func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error 
// func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error 
// func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error 
