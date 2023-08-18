package db

import (
	"context"

	"github.com/Jimbo8702/goreservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
}

type MongoBookingStore struct {
	client 	*mongo.Client
	coll 	*mongo.Collection
	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore{
	return &MongoBookingStore{
		client: client,
		coll: client.Database(DBNAME).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := resp.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
