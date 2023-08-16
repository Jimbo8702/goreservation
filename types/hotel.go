package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID 			primitive.ObjectID 		`bson:"_id,omitempty" json:"id,omitempty"`
	Name 		string 					`bson:"name" json:"name"`
	Location 	string 					`bson:"location" json:"location"`
	Rooms		[]primitive.ObjectID 	`bson:"rooms" json:"rooms"`
	Rating 		int 					`bson:"rating" json:"rating"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType 
	DoubleRoomType 
	SeaSideRoomType 
	DeluxeRoomType
)

type Room struct {
	ID 			primitive.ObjectID 	`bson:"_id,omitempty" json:"id,omitempty"`
	//small, normal, kingsize
	Size 		string 				`bson:"size" json:"size"`
	Seaside 	bool 				`bson:"type" json:"type"`
	Price 		float64	 			`bson:"price" json:"price"`		
	HotelID		primitive.ObjectID 	`bson:"hotelID" json:"hotelID"`
}