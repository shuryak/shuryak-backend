package articles

import (
	"context"
	"encoding/json"
	"fmt"
	v "github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/shuryak/shuryak-backend/internal/models"
	"github.com/shuryak/shuryak-backend/internal/utils"
	"github.com/shuryak/shuryak-backend/internal/utils/http-result"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.ArticleDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if len(dto.CustomId) < int(models.ArticleIdMinLimit) || len(dto.CustomId) > int(models.ArticleIdMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("id length < ", models.ArticleIdMinLimit, " or > ", models.ArticleIdMaxLimit))
		return
	}

	if len(dto.Name) < int(models.ArticleNameMinLimit) || len(dto.Name) > int(models.ArticleNameMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("name length < ", models.ArticleNameMinLimit, " or > ", models.ArticleNameMaxLimit))
		return
	}

	if !v.IsURL(dto.Thumbnail) {
		http_result.WriteError(&w, models.BadRequest, "invalid thumbnail")
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	// Checking for the existence of an article with this name
	var dbArticle models.User
	findFilter := bson.D{{"custom_id", dto.CustomId}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbArticle); err == nil {
		http_result.WriteError(&w, models.NotUniqueData, "article with this id already exists")
		return
	}
	findFilter = bson.D{{"name", dto.Name}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbArticle); err == nil {
		http_result.WriteError(&w, models.NotUniqueData, "article with this name already exists")
		return
	}

	_, err := collection.InsertOne(context.TODO(), models.Article{
		Id:          primitive.NewObjectID(),
		CustomId:    dto.CustomId,
		Name:        dto.Name,
		Author:      r.Context().Value(models.JwtClaimsKey).(jwt.MapClaims)["nickname"].(string),
		IsDraft:     dto.IsDraft,
		Thumbnail:   dto.Thumbnail,
		ArticleData: dto.ArticleData,
	})
	if err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	result := models.MetaArticle{
		Id:        dto.CustomId,
		Name:      dto.Name,
		Author:    r.Context().Value(models.JwtClaimsKey).(jwt.MapClaims)["nickname"].(string),
		IsDraft:   dto.IsDraft,
		Thumbnail: dto.Thumbnail,
	}

	json.NewEncoder(w).Encode(result)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.ArticleDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if len(dto.CustomId) < int(models.ArticleIdMinLimit) || len(dto.CustomId) > int(models.ArticleIdMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("id length < ", models.ArticleIdMinLimit, " or > ", models.ArticleIdMaxLimit))
		return
	}

	if len(dto.Name) < int(models.ArticleNameMinLimit) || len(dto.Name) > int(models.ArticleNameMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("name length < ", models.ArticleNameMinLimit, " or > ", models.ArticleNameMaxLimit))
		return
	}

	if !v.IsURL(dto.Thumbnail) {
		http_result.WriteError(&w, models.BadRequest, "invalid thumbnail")
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	findFilter := bson.D{{"custom_id", dto.CustomId}}

	var dbArticle models.Article
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbArticle); err != nil {
		http_result.WriteError(&w, models.BadRequest, "article with this id doesn't exist")
		return
	}

	if dbArticle.Author != r.Context().Value(models.JwtClaimsKey).(jwt.MapClaims)["nickname"].(string) {
		http_result.WriteError(&w, models.BadRequest, "you're not the author of this article")
		return
	}

	articleUpdated := models.Article{
		Id:          dbArticle.Id,
		CustomId:    dbArticle.CustomId,
		Author:      dbArticle.Author,
		Name:        dto.Name,
		IsDraft:     dto.IsDraft,
		Thumbnail:   dto.Thumbnail,
		ArticleData: dto.ArticleData,
	}

	update := bson.D{{"$set", articleUpdated}}

	if _, err := collection.UpdateOne(context.TODO(), findFilter, update); err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	json.NewEncoder(w).Encode(models.ArticleDTO{
		CustomId:    articleUpdated.CustomId,
		Name:        articleUpdated.Name,
		Author:      articleUpdated.Author,
		IsDraft:     articleUpdated.IsDraft,
		Thumbnail:   articleUpdated.Thumbnail,
		ArticleData: articleUpdated.ArticleData,
	})
}

func FindOneHandler(w http.ResponseWriter, r *http.Request) {
	var query models.FindOneExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if query.Query == "" {
		http_result.WriteError(&w, models.InvalidFieldLength, "empty query string")
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	var result models.MetaArticle

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		http_result.WriteEmpty(&w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func FindManyHandler(w http.ResponseWriter, r *http.Request) {
	var query models.FindManyExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http_result.WriteError(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	// region Validation
	if query.Query == "" {
		http_result.WriteError(&w, models.InvalidFieldLength, "empty query string")
		return
	}

	if query.Count > uint(models.FindMaxLimit) {
		http_result.WriteError(&w, models.BadRequest, fmt.Sprint("count > ", models.FindMaxLimit))
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	options := options.Find()
	options.SetLimit(int64(query.Count))
	options.SetSkip(int64(query.Offset))
	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		http_result.WriteEmpty(&w)
		return
	}

	var results []*models.MetaArticle

	for cur.Next(context.TODO()) {
		var document models.MetaArticle
		err := cur.Decode(&document)
		if err != nil {
			http_result.WriteError(&w, models.InternalError, "internal error")
			return
		}

		results = append(results, &document)
	}

	if err := cur.Err(); err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}

func GetByCustomIdHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.ArticleCustomIdDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if len(dto.CustomId) < int(models.ArticleIdMinLimit) || len(dto.CustomId) > int(models.ArticleIdMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("id length < ", models.ArticleIdMinLimit, " or > ", models.ArticleIdMaxLimit))
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	filter := bson.D{{"custom_id", dto.CustomId}}

	var result models.ArticleDTO

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		http_result.WriteError(&w, models.BadRequest, "Article with this id doesn't exist")
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetDraftsListHandler(w http.ResponseWriter, r *http.Request) {
	var query models.GetListExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if query.Count > uint(models.FindMaxLimit) {
		http_result.WriteError(&w, models.BadRequest, fmt.Sprint("count > ", models.FindMaxLimit))
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	options := options.Find()
	options.SetLimit(int64(query.Count))
	options.SetSkip(int64(query.Offset))

	filter := bson.D{{"$and", []bson.D{
		bson.D{{"is_draft", true}},
		bson.D{{"author", r.Context().Value(models.JwtClaimsKey).(jwt.MapClaims)["nickname"].(string)}},
	}}}

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		http_result.WriteEmpty(&w)
		return
	}

	var results []*models.MetaArticle

	for cur.Next(context.TODO()) {
		var document models.MetaArticle
		err := cur.Decode(&document)
		if err != nil {
			http_result.WriteError(&w, models.InternalError, "internal error")
			return
		}

		results = append(results, &document)
	}

	if err := cur.Err(); err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	var query models.GetListExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if query.Count > uint(models.FindMaxLimit) {
		http_result.WriteError(&w, models.BadRequest, fmt.Sprint("count > ", models.FindMaxLimit))
		return
	}
	// endregion Validation

	collection := utils.Mongo.Database("shuryakDb").Collection("articles")

	options := options.Find()
	options.SetLimit(int64(query.Count))
	options.SetSkip(int64(query.Offset))

	filter := bson.D{{"is_draft", false}}

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		http_result.WriteEmpty(&w)
		return
	}

	var results []*models.MetaArticle

	for cur.Next(context.TODO()) {
		var document models.MetaArticle
		err := cur.Decode(&document)
		if err != nil {
			http_result.WriteError(&w, models.InternalError, "internal error")
			return
		}

		results = append(results, &document)
	}

	if err := cur.Err(); err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	cur.Close(context.TODO())

	json.NewEncoder(w).Encode(results)
}
