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

type ResourceResp struct {
	Results int 	`json:"results"`
	Page 	int		`json:"page"`
	Data 	any 	`json:"data"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int 
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{}
	if params.Rating != 0 {
		filter["rating"] = params.Rating
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("hotels")
	}
	resp := ResourceResp{
		Results: len(hotels),
		Page: int(params.Page),
		Data: hotels,
	}
	return c.JSON(resp)
}

