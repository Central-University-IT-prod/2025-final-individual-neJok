package model

import "github.com/google/uuid"

type Client struct {
	ClientID uuid.UUID `bson:"_id" json:"client_id" binding:"required" format:"uuid"`
	Login    string    `bson:"login" json:"login" binding:"required"`
	Age      int       `bson:"age" json:"age" binding:"required,min=0,max=100"`
	Location string    `bson:"location" json:"location" binding:"required"`
	Gender   string    `bson:"gender" json:"gender" binding:"required,oneof=MALE FEMALE"`
}

type ClientCreateOrUpdateRequest []Client
