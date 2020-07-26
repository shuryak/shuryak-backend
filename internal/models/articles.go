package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetaArticle struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name"`
	IsDraft   bool               `json:"is_draft"`
	Thumbnail string             `json:"thumbnail"`
}

type ArticleIdDTO struct {
	Id primitive.ObjectID `json:"article_id"`
}

type ArticleDTO struct {
	Name        string      `json:"name"`
	IsDraft     bool        `json:"is_draft"`
	Thumbnail   string      `json:"thumbnail"`
	ArticleData interface{} `json:"article_data"`
}

type Article struct {
	Name        string      `bson:"name"`
	IsDraft     bool        `bson:"is_draft"`
	Thumbnail   string      `bson:"thumbnail"`
	ArticleData interface{} `bson:"article_data"`
}

func (dto ArticleDTO) ToArticle() *Article {
	return &Article{
		Name:        dto.Name,
		IsDraft:     dto.IsDraft,
		Thumbnail:   dto.Thumbnail,
		ArticleData: dto.ArticleData,
	}
}
