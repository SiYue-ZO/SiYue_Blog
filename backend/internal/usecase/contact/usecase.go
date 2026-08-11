package contact

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct{ repo repo.ContactRepo }

var ErrInvalidInput = errors.New("invalid contact input")

func NewUseCase(contactRepo repo.ContactRepo) usecase.ContactUseCase {
	return &useCase{repo: contactRepo}
}

func (uc *useCase) Create(ctx context.Context, input usecase.CreateContactInput) (*entity.ContactMessage, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Message = strings.TrimSpace(input.Message)
	if utf8.RuneCountInString(input.Name) < 1 || utf8.RuneCountInString(input.Name) > 100 {
		return nil, fmt.Errorf("%w: name must contain 1 to 100 characters", ErrInvalidInput)
	}
	address, err := mail.ParseAddress(input.Email)
	if err != nil || !strings.EqualFold(address.Address, input.Email) {
		return nil, fmt.Errorf("%w: a valid email address is required", ErrInvalidInput)
	}
	if utf8.RuneCountInString(input.Message) < 10 || utf8.RuneCountInString(input.Message) > 5000 {
		return nil, fmt.Errorf("%w: message must contain 10 to 5000 characters", ErrInvalidInput)
	}
	message := &entity.ContactMessage{
		Name: input.Name, Email: input.Email, Message: input.Message,
		Status: entity.ContactStatusUnread, IPHash: input.IPHash,
	}
	if err := uc.repo.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("create contact message: %w", err)
	}
	return message, nil
}

func (uc *useCase) List(ctx context.Context, page, pageSize int, status string) ([]entity.ContactMessage, int64, error) {
	if status != "" && status != string(entity.ContactStatusUnread) && status != string(entity.ContactStatusRead) && status != string(entity.ContactStatusArchived) {
		return nil, 0, fmt.Errorf("invalid contact status")
	}
	return uc.repo.List(ctx, (page-1)*pageSize, pageSize, status)
}

func (uc *useCase) UpdateStatus(ctx context.Context, id uint64, status string) error {
	if status != string(entity.ContactStatusUnread) && status != string(entity.ContactStatusRead) && status != string(entity.ContactStatusArchived) {
		return fmt.Errorf("invalid contact status")
	}
	return uc.repo.UpdateStatus(ctx, id, entity.ContactStatus(status))
}

func (uc *useCase) Delete(ctx context.Context, id uint64) error { return uc.repo.Delete(ctx, id) }
