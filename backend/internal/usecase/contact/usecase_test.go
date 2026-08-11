package contact

import (
	"context"
	"testing"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/usecase"
)

type memoryContactRepo struct{ messages []entity.ContactMessage }

func (r *memoryContactRepo) Create(_ context.Context, message *entity.ContactMessage) error {
	message.ID = uint64(len(r.messages) + 1)
	r.messages = append(r.messages, *message)
	return nil
}
func (r *memoryContactRepo) List(_ context.Context, _, _ int, _ string) ([]entity.ContactMessage, int64, error) {
	return r.messages, int64(len(r.messages)), nil
}
func (r *memoryContactRepo) GetByID(_ context.Context, id uint64) (*entity.ContactMessage, error) {
	return &r.messages[id-1], nil
}
func (r *memoryContactRepo) UpdateStatus(_ context.Context, id uint64, status entity.ContactStatus) error {
	r.messages[id-1].Status = status
	return nil
}
func (r *memoryContactRepo) Delete(_ context.Context, id uint64) error {
	r.messages = append(r.messages[:id-1], r.messages[id:]...)
	return nil
}

func TestCreateContactMessage(t *testing.T) {
	repo := &memoryContactRepo{}
	uc := NewUseCase(repo)
	message, err := uc.Create(context.Background(), usecase.CreateContactInput{
		Name: "  Reader  ", Email: "READER@example.com", Message: "This is a useful message.", IPHash: "hash",
	})
	if err != nil {
		t.Fatal(err)
	}
	if message.Name != "Reader" || message.Email != "reader@example.com" || message.Status != entity.ContactStatusUnread {
		t.Fatalf("unexpected message: %#v", message)
	}
}

func TestCreateContactMessageRejectsInvalidInput(t *testing.T) {
	uc := NewUseCase(&memoryContactRepo{})
	if _, err := uc.Create(context.Background(), usecase.CreateContactInput{Name: "Reader", Email: "invalid", Message: "too short"}); err == nil {
		t.Fatal("invalid message was accepted")
	}
}
