package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jimbo8702/goreservation/db"
	"github.com/Jimbo8702/goreservation/types"
	"github.com/gofiber/fiber/v2"
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

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg: "internal server error",
		})
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPersons: params.NumPersons,
	}

	fmt.Print(booking)

	return nil
}