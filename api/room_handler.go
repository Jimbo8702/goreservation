package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jimbo8702/goreservation/db"
	"github.com/Jimbo8702/goreservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if err := params.validate(); err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}
	
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return ErrInvalidID()
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return NewError(http.StatusInternalServerError, "internal server error")
	}

	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params) 
	if err != nil {
		return err
	}
	if !ok {
		return NewError(http.StatusBadRequest, "room already booked")
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceNotFound("rooms")
	}
	return c.JSON(rooms)
}

// func (h *RoomHandler) HandleGetBookingsForRoom(c *fiber.Ctx) error {}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		} ,
		"tillDate": bson.M{
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