package river

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
)

type service struct {
	Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

var (
	riverMutex   sync.Mutex
	clientsMutex sync.Mutex
)

func (s *service) UpdateRiver(ctx context.Context, id string, dataChan chan *UpdateRiver) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-dataChan:

				if update.Height < 5 {
					update.Status = "Aman"
				} else if update.Height < 10 {
					update.Status = "Siaga"
				} else {
					update.Status = "Bahaya"
				}

				riverMutex.Lock()
				err := s.Repository.UpdateRiver(ctx, update)
				if err != nil {
					log.Printf("Failed to update river: %v", err)
				}
				riverMutex.Unlock()

				msg := UpdateRiver{
					Id:     update.Id,
					Height: update.Height,
					Status: update.Status,
				}
				fmt.Println("From service:", msg)

				clientsMutex.Lock()
				for client := range clients {
					err := client.WriteJSON(msg)
					if err != nil {
						delete(clients, client)
						log.Printf("Failed to write JSON to client: %v", err)
					}
				}
				clientsMutex.Unlock()
			}
		}
	}()

	return nil
}

func (s *service) AddRiver(ctx context.Context, req *CreateRiverRequest) error {
	river := River{
		Id:        req.Id,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
	}
	err := s.Repository.CreateRiver(ctx, &river)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ViewRiverById(ctx context.Context, id string) (*RiverResponse, error) {
	res, err := s.Repository.GetRiverById(ctx, id)
	if err != nil {
		return &RiverResponse{}, err
	}

	news := RiverResponse{
		Id:        res.Id,
		Latitude:  res.Latitude,
		Longitude: res.Longitude,
		Location:  res.Location,
		Height:    res.Height,
		Status:    res.Status,
	}

	return &news, nil
}

func (s *service) RemoveRiver(ctx context.Context, id string) error {
	err := s.Repository.DeleteRiver(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("river with ID %s not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ChangeRiverDetail(ctx context.Context, req *UpdateRiverRequest, id string) error {
	river := River{
		Id:        req.Id,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
		Height:    req.Height,
		Status:    req.Status,
	}

	err := s.Repository.UpdateRiverDetail(ctx, &river, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("river with ID %s not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewRiver(ctx context.Context) (*[]RiverResponse, error) {
	r, err := s.Repository.GetRiver(ctx)
	if err != nil {
		return &[]RiverResponse{}, err
	}

	var rivers []RiverResponse

	for _, res := range *r {
		riverRes := NewRiverResponse(res)

		rivers = append(rivers, *riverRes)
	}

	return &rivers, nil
}

func (s *service) ViewRiverByStatus(ctx context.Context, status string) (*[]RiverResponse, error) {
	r, err := s.Repository.GetRiverByStatus(ctx, status)
	if err != nil {
		return &[]RiverResponse{}, err
	}

	var rivers []RiverResponse

	for _, res := range *r {
		riverRes := NewRiverResponse(res)

		rivers = append(rivers, *riverRes)
	}

	return &rivers, nil
}

func NewRiverResponse(res River) *RiverResponse {
	return &RiverResponse{
		Id:        res.Id,
		Latitude:  res.Latitude,
		Longitude: res.Longitude,
		Location:  res.Location,
		Height:    res.Height,
		Status:    res.Status,
	}
}

func (s *service) SearchRiver(ctx context.Context, searchTerm string) (*RiverResponse, error) {
	river, err := s.Repository.FindRiver(ctx, searchTerm)
	if err != nil {
		return nil, err
	}

	riverRes := NewRiverResponse(*river)

	return riverRes, nil
}

func (s *service) SortRiver(ctx context.Context, sortBy string) (*[]RiverResponse, error) {
	rivers, err := s.Repository.FilterRiver(ctx, sortBy)
	if err != nil {
		return nil, err
	}

	var riverResponses []RiverResponse
	for _, river := range *rivers {
		riverRes := NewRiverResponse(river)

		riverResponses = append(riverResponses, *riverRes)
	}

	return &riverResponses, nil
}

func (s *service) ViewAllRiverCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllRiverCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}
