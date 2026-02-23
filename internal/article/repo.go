package article

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, input CreateArticleInput) (Article, error)
	List(ctx context.Context) ([]Article, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, input CreateArticleInput) (Article, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return Article{}, err
	}
	defer tx.Rollback(ctx)

	const articleQ = `
		INSERT INTO articles (title, content, image_url)
		VALUES ($1, $2, $3)
		RETURNING id, title, content, image_url, created_at, updated_at
	`

	var article Article
	err = tx.QueryRow(ctx, articleQ, input.Title, input.Content, input.ImageURL).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.ImageURL,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return Article{}, err
	}

	const tagQ = `
		INSERT INTO article_tags (
			article_id,
			tag_type,
			tag_value,
			player_image_url,
			player_fide_elo,
			opening_board_image_url,
			opening_variation_count
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, article_id, tag_type, tag_value, player_image_url, player_fide_elo, opening_board_image_url, opening_variation_count, created_at
	`

	article.Tags = make([]ArticleTag, 0, len(input.Tags))
	for _, t := range input.Tags {
		var tag ArticleTag
		err = tx.QueryRow(
			ctx,
			tagQ,
			article.ID,
			t.TagType,
			t.TagValue,
			t.PlayerImageURL,
			t.PlayerFIDEElo,
			t.OpeningBoardImageURL,
			t.OpeningVariationCount,
		).Scan(
			&tag.ID,
			&tag.ArticleID,
			&tag.TagType,
			&tag.TagValue,
			&tag.PlayerImageURL,
			&tag.PlayerFIDEElo,
			&tag.OpeningBoardImageURL,
			&tag.OpeningVariationCount,
			&tag.CreatedAt,
		)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	if err := tx.Commit(ctx); err != nil {
		return Article{}, err
	}

	return article, nil
}

func (r *repository) List(ctx context.Context) ([]Article, error) {
	const articleQ = `
		SELECT id, title, content, image_url, created_at, updated_at
		FROM articles
		ORDER BY id DESC
		LIMIT 100
	`

	rows, err := r.db.Query(ctx, articleQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.ImageURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}

		tags, err := r.listTagsByArticleID(ctx, a.ID)
		if err != nil {
			return nil, err
		}
		a.Tags = tags
		articles = append(articles, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *repository) listTagsByArticleID(ctx context.Context, articleID int64) ([]ArticleTag, error) {
	const tagQ = `
		SELECT id, article_id, tag_type, tag_value, player_image_url, player_fide_elo, opening_board_image_url, opening_variation_count, created_at
		FROM article_tags
		WHERE article_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, tagQ, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]ArticleTag, 0)
	for rows.Next() {
		var t ArticleTag
		if err := rows.Scan(
			&t.ID,
			&t.ArticleID,
			&t.TagType,
			&t.TagValue,
			&t.PlayerImageURL,
			&t.PlayerFIDEElo,
			&t.OpeningBoardImageURL,
			&t.OpeningVariationCount,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
