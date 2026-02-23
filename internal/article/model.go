package article

import "time"

type TagType string

const (
	TagTypeOpening   TagType = "opening"
	TagTypePlayer    TagType = "player"
	TagTypeVariation TagType = "variation"
)

type Article struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	ImageURL  string       `json:"image_url"`
	Tags      []ArticleTag `json:"tags"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type ArticleTag struct {
	ID                    int64    `json:"id"`
	ArticleID             int64    `json:"article_id"`
	TagType               TagType  `json:"tag_type"`
	TagValue              string   `json:"tag_value"`
	PlayerImageURL        *string  `json:"player_image_url,omitempty"`
	PlayerFIDEElo         *int     `json:"player_fide_elo,omitempty"`
	OpeningBoardImageURL  *string  `json:"opening_board_image_url,omitempty"`
	OpeningVariationCount *int     `json:"opening_variation_count,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

type CreateArticleInput struct {
	Title    string              `json:"title" validate:"required,min=1,max=255"`
	Content  string              `json:"content" validate:"required,min=1"`
	ImageURL string              `json:"image_url" validate:"required,min=1,max=2000"`
	Tags     []CreateArticleTagInput `json:"tags" validate:"required,min=1,dive"`
}

type CreateArticleTagInput struct {
	TagType               TagType  `json:"tag_type" validate:"required"`
	TagValue              string   `json:"tag_value" validate:"required,min=1,max=255"`
	PlayerImageURL        *string  `json:"player_image_url,omitempty"`
	PlayerFIDEElo         *int     `json:"player_fide_elo,omitempty"`
	OpeningBoardImageURL  *string  `json:"opening_board_image_url,omitempty"`
	OpeningVariationCount *int     `json:"opening_variation_count,omitempty"`
}
