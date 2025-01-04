package models

import (
    "time"
)

type Product struct {
    ProductID   string    `bson:"product_id,omitempty" json:"product_id"`
    Name        string    `bson:"name,omitempty" json:"name"`
    Description string    `bson:"description,omitempty" json:"description"`
    Price       float64   `bson:"price,omitempty" json:"price"`
    Stock       int       `bson:"stock,omitempty" json:"stock"`
    CreatedAt   time.Time `bson:"created_at,omitempty" json:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}