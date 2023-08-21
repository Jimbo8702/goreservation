package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jimbo8702/goreservation/db"
	"github.com/Jimbo8702/goreservation/types"
	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	FromDate time.Time 	`json:"fromDate"`
	TillDate time.Time 	`json:"tillDate"`
	NumPersons int		`json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var (
		params BookRoomParams
		roomID = c.Params("id")
	)
		if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if err := params.validate(); err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return NewError(http.StatusInternalServerError, "internal server error")
	}

	ok, err := h.isRoomAvailableForBooking(c.Context(), roomID, params) 
	if err != nil {
		return err
	}
	if !ok {
		return NewError(http.StatusBadRequest, "room already booked")
	}
	booking, err := types.NewBookingFromParams(types.CreateBookingParms{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPersons: params.NumPersons,
	})
	if err != nil {
		return err
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), db.Map{})
	if err != nil {
		return ErrResourceNotFound("rooms")
	}
	return c.JSON(rooms)
}

// func (h *RoomHandler) HandleGetBookingsForRoom(c *fiber.Ctx) error {}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID string, params BookRoomParams) (bool, error) {
	where := db.Map{
		"roomID": roomID,
		"fromDate": db.Map{
			"$gte": params.FromDate,
		} ,
		"tillDate": db.Map{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}