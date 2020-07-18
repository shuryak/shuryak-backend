package  models

type ArticleDTO struct {
	Name        string      `json:"name"`
	ArticleData interface{} `json:"article_data"`
}
