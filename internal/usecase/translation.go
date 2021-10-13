package usecase

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
)

// TranslationUseCase -.
// Business logic.
//
// Methods are grouped by area of application (on a common basis)
// Each group has its own structure
// One file - one structure
// Repositories, webapi, rpc, and other business logic structures are injected into business logic structures (see Dependency Injection).
type TranslationUseCase struct {
	repo   TranslationRepo
	webAPI TranslationWebAPI
}

// New -.
func New(r TranslationRepo, w TranslationWebAPI) *TranslationUseCase {
	return &TranslationUseCase{
		repo:   r,
		webAPI: w,
	}
}

// History - getting translate history from store.
func (uc *TranslationUseCase) History(ctx context.Context) ([]entity.Translation, error) {
	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}

// Translate -.
func (uc *TranslationUseCase) Translate(ctx context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.webAPI.Translate: %w", err)
	}

	err = uc.repo.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
