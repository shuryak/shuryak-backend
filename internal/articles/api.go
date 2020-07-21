package articles

import (
	"context"
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	w.Header().Set("Content-Type", "application/json")

	var query models.FindExpression

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

	var result models.ArticleDTO

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		emptyResult := struct {
		}{}
		json.NewEncoder(w).Encode(emptyResult)
		return
	}

	json.NewEncoder(w).Encode(result)
}
