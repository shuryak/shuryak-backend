package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MetaArticle struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name"`
}

type ArticleDTO struct {
	Name        string      `json:"name"`
	ArticleData interface{} `json:"article_data"`
}
