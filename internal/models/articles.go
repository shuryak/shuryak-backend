package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetaArticle struct {
	Id        string `json:"id" bson:"custom_id"`
	Name      string `json:"name"`
	IsDraft   bool   `json:"is_draft"`
	Thumbnail string `json:"thumbnail"`
}

type ArticleCustomIdDTO struct {
	CustomId string `json:"id"`
}

type ArticleDTO struct {
	CustomId    string                 `json:"id" bson:"custom_id"`
	Name        string                 `json:"name" bson:"name"`
	IsDraft     bool                   `json:"is_draft" bson:"is_draft"`
	Thumbnail   string                 `json:"thumbnail" bson:"thumbnail"`
	ArticleData map[string]interface{} `json:"article_data" bson:"article_data"`
	// https://medium.com/rungo/working-with-json-in-go-7e3a37c5a07b
}

type Article struct {
	Id          primitive.ObjectID `bson:"_id"`
	CustomId    string             `bson:"custom_id"`
	Name        string             `bson:"name"`
	IsDraft     bool               `bson:"is_draft"`
	Thumbnail   string             `bson:"thumbnail"`
	ArticleData interface{}        `bson:"article_data"`
}

func (dto ArticleDTO) ToArticle() *Article {
	return &Article{
		Id:          primitive.NewObjectID(),
		CustomId:    dto.CustomId,
		Name:        dto.Name,
		IsDraft:     dto.IsDraft,
		Thumbnail:   dto.Thumbnail,
		ArticleData: dto.ArticleData,
	}
}
