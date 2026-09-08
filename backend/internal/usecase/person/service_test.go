package person

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/Nok1o/lab1_rsoi/backend/internal/domain"
)

func TestServiceCreate(t *testing.T) {
	repository := newMemoryRepository()
	service := NewService(repository)

	id, err := service.Create(context.Background(), CreateInput{
		Name: "  Ada Lovelace  ", Age: 36, Address: "London", Work: "Mathematician",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	created := repository.persons[id]
	if created.Name != "Ada Lovelace" || created.Age != 36 || created.Address != "London" || created.Work != "Mathematician" {
		t.Fatalf("created person = %+v", created)
	}
}

func TestServiceGetByID(t *testing.T) {
	repository := newMemoryRepository()
	repository.persons[7] = domain.Person{ID: 7, Name: "Grace Hopper"}
	service := NewService(repository)

	got, err := service.GetByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ID != 7 || got.Name != "Grace Hopper" {
		t.Fatalf("GetByID() = %+v", got)
	}
}

func TestServiceList(t *testing.T) {
	repository := newMemoryRepository()
	repository.persons[2] = domain.Person{ID: 2, Name: "Second"}
	repository.persons[1] = domain.Person{ID: 1, Name: "First"}
	service := NewService(repository)

	got, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("List() = %+v", got)
	}
}

func TestServiceUpdateKeepsMissingFields(t *testing.T) {
	repository := newMemoryRepository()
	repository.persons[1] = domain.Person{
		ID: 1, Name: "Old name", Age: 31, Address: "Old address", Work: "Engineer",
	}
	service := NewService(repository)
	name := "New name"
	address := "New address"

	got, err := service.Update(context.Background(), 1, UpdateInput{Name: &name, Address: &address})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if got.Name != name || got.Address != address || got.Age != 31 || got.Work != "Engineer" {
		t.Fatalf("Update() = %+v", got)
	}
}

func TestServiceDelete(t *testing.T) {
	repository := newMemoryRepository()
	repository.persons[1] = domain.Person{ID: 1, Name: "Deleted"}
	service := NewService(repository)

	if err := service.Delete(context.Background(), 1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, ok := repository.persons[1]; ok {
		t.Fatal("Delete() did not remove person")
	}
	if err := service.Delete(context.Background(), 1); !errors.Is(err, domain.ErrPersonNotFound) {
		t.Fatalf("second Delete() error = %v", err)
	}
}

type memoryRepository struct {
	persons map[int64]domain.Person
	nextID  int64
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{persons: make(map[int64]domain.Person), nextID: 1}
}

func (r *memoryRepository) Create(_ context.Context, person domain.Person) (int64, error) {
	person.ID = r.nextID
	r.nextID++
	r.persons[person.ID] = person
	return person.ID, nil
}

func (r *memoryRepository) GetByID(_ context.Context, id int64) (domain.Person, error) {
	person, ok := r.persons[id]
	if !ok {
		return domain.Person{}, domain.ErrPersonNotFound
	}
	return person, nil
}

func (r *memoryRepository) List(_ context.Context) ([]domain.Person, error) {
	persons := make([]domain.Person, 0, len(r.persons))
	for _, person := range r.persons {
		persons = append(persons, person)
	}
	sort.Slice(persons, func(i, j int) bool { return persons[i].ID < persons[j].ID })
	return persons, nil
}

func (r *memoryRepository) Update(_ context.Context, person domain.Person) error {
	if _, ok := r.persons[person.ID]; !ok {
		return domain.ErrPersonNotFound
	}
	r.persons[person.ID] = person
	return nil
}

func (r *memoryRepository) Delete(_ context.Context, id int64) error {
	if _, ok := r.persons[id]; !ok {
		return domain.ErrPersonNotFound
	}
	delete(r.persons, id)
	return nil
}
