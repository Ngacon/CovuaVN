package article

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidTagPayload = errors.New("invalid tag payload")

type Service interface {
	Create(ctx context.Context, input CreateArticleInput) (Article, error)
	List(ctx context.Context) ([]Article, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input CreateArticleInput) (Article, error) {
	if err := validateTagPayload(input.Tags); err != nil {
		return Article{}, fmt.Errorf("%w: %v", ErrInvalidTagPayload, err)
	}
	return s.repo.Create(ctx, input)
}

func (s *service) List(ctx context.Context) ([]Article, error) {
	return s.repo.List(ctx)
}

func validateTagPayload(tags []CreateArticleTagInput) error {
	if len(tags) == 0 {
		return errors.New("tags is required")
	}

	for i, tag := range tags {
		switch tag.TagType {
		case TagTypeOpening:
			if isEmptyStringPtr(tag.OpeningBoardImageURL) {
				return fmt.Errorf("tags[%d]: opening_board_image_url is required for opening tag", i)
			}
			if tag.OpeningVariationCount == nil || *tag.OpeningVariationCount <= 0 {
				return fmt.Errorf("tags[%d]: opening_variation_count must be > 0 for opening tag", i)
			}
		case TagTypePlayer:
			if isEmptyStringPtr(tag.PlayerImageURL) {
				return fmt.Errorf("tags[%d]: player_image_url is required for player tag", i)
			}
			if tag.PlayerFIDEElo == nil || *tag.PlayerFIDEElo <= 0 {
				return fmt.Errorf("tags[%d]: player_fide_elo must be > 0 for player tag", i)
			}
		case TagTypeVariation:
			continue
		default:
			return fmt.Errorf("tags[%d]: tag_type must be one of opening, player, variation", i)
		}
	}

	return nil
}

func isEmptyStringPtr(v *string) bool {
	return v == nil || *v == ""
}
