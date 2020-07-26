package articles

import (
	"context"
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/models"
	"github.com/shuryak/shuryak-backend/internal/writers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.ArticleDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writers.ErrorWriter(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	// Checking for the existence of an article with this name
	var dbArticle models.User
	findFilter := bson.D{{"name", dto.Name}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbArticle); err == nil {
		writers.ErrorWriter(&w, models.NotUniqueData, "Article with this name already exists")
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), dto.ToArticle())
	if err != nil {
		writers.ErrorWriter(&w, models.InternalError, "Internal error")
		return
	}

	result := models.MetaArticle{
		Id:        insertResult.InsertedID.(primitive.ObjectID),
		Name:      dto.Name,
		IsDraft:   dto.IsDraft,
		Thumbnail: dto.Thumbnail,
	}

	json.NewEncoder(w).Encode(result)
}

func FindOneHandler(w http.ResponseWriter, r *http.Request) {
	var query models.FindOneExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writers.ErrorWriter(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	var result models.MetaArticle

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		writers.EmptyResultWriter(&w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func FindManyHandler(w http.ResponseWriter, r *http.Request) {
	var query models.FindManyExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writers.ErrorWriter(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	options := options.Find()
	options.SetLimit(int64(query.Count))
	options.SetSkip(int64(query.Offset))
	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		writers.EmptyResultWriter(&w)
		return
	}

	var results []*models.MetaArticle

	for cur.Next(context.TODO()) {
		var document models.MetaArticle
		err := cur.Decode(&document)
		if err != nil {
			writers.ErrorWriter(&w, models.InternalError, "Internal error")
			return
		}

		results = append(results, &document)
	}

	if err := cur.Err(); err != nil {
		writers.ErrorWriter(&w, models.InternalError, "Internal error")
		return
	}

	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}

func GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.ArticleIdDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writers.ErrorWriter(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	filter := bson.D{{"_id", dto.Id}}

	var result models.ArticleDTO

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		writers.ErrorWriter(&w, models.BadRequest, "Article with this id doesn't exist")
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	var query models.GetListExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writers.ErrorWriter(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	options := options.Find()
	options.SetLimit(int64(query.Count))
	options.SetSkip(int64(query.Offset))

	cur, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		writers.EmptyResultWriter(&w)
		return
	}

	var results []*models.MetaArticle

	for cur.Next(context.TODO()) {
		var document models.MetaArticle
		err := cur.Decode(&document)
		if err != nil {
			writers.ErrorWriter(&w, models.InternalError, "Internal error")
			return
		}

		results = append(results, &document)
	}

	if err := cur.Err(); err != nil {
		writers.ErrorWriter(&w, models.InternalError, "Internal error")
		return
	}

	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}
