package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductID   string             `bson:"product_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"descriptsion,omitempty"`
	Price       float64            `bson:"price,omitempty"`
	Stock       int                `bson:"stock,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}
