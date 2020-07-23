package articles

import (
	"context"
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var article models.ArticleDTO

	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	insertResult, err := collection.InsertOne(context.TODO(), article)
	if err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.InternalError,
			Message:   "Internal error",
		}
		json.NewEncoder(w).Encode(errorMessage)
	}

	result := struct {
		ArticleId primitive.ObjectID `json:"article_id"`
	}{
		insertResult.InsertedID.(primitive.ObjectID),
	}

	json.NewEncoder(w).Encode(result)
}

func FindOneHandler(w http.ResponseWriter, r *http.Request) {
	var query models.FindOneExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	var result models.MetaArticle

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		emptyResult := struct {
		}{}
		json.NewEncoder(w).Encode(emptyResult)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func FindManyHandler(w http.ResponseWriter, r *http.Request) {
	var query models.FindManyExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	options := options.Find()
	options.SetLimit(int64(query.Count))
	options.SetSkip(int64(query.Offset))
	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		emptyResult := struct {
		}{}
		json.NewEncoder(w).Encode(emptyResult)
		return
	}

	var results []*models.MetaArticle

	for cur.Next(context.TODO()) {
		var document models.MetaArticle
		err := cur.Decode(&document)
		if err != nil {
			errorMessage := models.ErrorDTO{
				ErrorCode: models.InternalError,
				Message:   "Internal error",
			}
			json.NewEncoder(w).Encode(errorMessage)
			return
		}

		results = append(results, &document)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}
