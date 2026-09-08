package http

import (
	"github.com/Nok1o/lab1_rsoi/backend/internal/domain"
	"github.com/Nok1o/lab1_rsoi/backend/internal/usecase/person"
)

type createPersonRequest struct {
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type updatePersonRequest struct {
	Name    *string `json:"name"`
	Age     *int32  `json:"age"`
	Address *string `json:"address"`
	Work    *string `json:"work"`
}

type personResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type validationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (request createPersonRequest) toInput() person.CreateInput {
	return person.CreateInput{
		Name:    request.Name,
		Age:     request.Age,
		Address: request.Address,
		Work:    request.Work,
	}
}

func (request updatePersonRequest) toInput() person.UpdateInput {
	return person.UpdateInput{
		Name:    request.Name,
		Age:     request.Age,
		Address: request.Address,
		Work:    request.Work,
	}
}

func toPersonResponse(person domain.Person) personResponse {
	return personResponse{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	}
}

func toPersonResponses(persons []domain.Person) []personResponse {
	responses := make([]personResponse, 0, len(persons))
	for _, person := range persons {
		responses = append(responses, toPersonResponse(person))
	}
	return responses
}
