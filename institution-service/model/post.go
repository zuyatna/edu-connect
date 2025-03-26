package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	PostID        uuid.UUID      `json:"post_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InstitutionID uuid.UUID      `json:"institution_id" gorm:"type:uuid; not null"`
	Title         string         `json:"title" gorm:"type:varchar(255); not null"`
	Body          string         `json:"body" gorm:"type:text; not null"`
	DateStart     time.Time      `json:"date_start" gorm:"type:timestamp; not null"`
	DateEnd       time.Time      `json:"date_end" gorm:"type:timestamp; not null"`
	FundTarget    float64        `json:"fund_target" gorm:"type:float; not null"`
	FundAchieved  float64        `json:"fund_achieved" gorm:"type:float; default:0"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
	Institution   Institution    `json:"institution" gorm:"foreignKey:InstitutionID;references:InstitutionID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type PostRequest struct {
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	DateStart  time.Time `json:"date_start"`
	DateEnd    time.Time `json:"date_end"`
	FundTarget float64   `json:"fund_target"`
}

type PostResponse struct {
	PostID       string  `json:"post_id"`
	Title        string  `json:"title"`
	Body         string  `json:"body"`
	DateStart    string  `json:"date_start"`
	DateEnd      string  `json:"date_end"`
	FundTarget   float32 `json:"fund_target"`
	FundAchieved float32 `json:"fund_achieved"`
}

type PostFundAchievedResponse struct {
	PostID       string  `json:"post_id"`
	FundAchieved float32 `json:"fund_achieved"`
}

type PostDeleteResponse struct {
	Message string `json:"message"`
}
