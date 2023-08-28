package river

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
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

func (s *service) UpdateRiver(ctx context.Context) error {
	riverChans := make(map[string]chan *UpdateRiver)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				riverIDs, err := s.Repository.GetRiverId(ctx)
				if err != nil {
					log.Printf("failed to get river IDs: %v", err)
					// Continue with the loop even if there's an error
					continue
				}

				// Create a new map to store the updated river channels
				newRiverChans := make(map[string]chan *UpdateRiver)

				for _, id := range riverIDs {
					riverMutex.Lock()
					if _, exists := riverChans[id]; !exists {
						riverChans[id] = make(chan *UpdateRiver)
						go s.Repository.UpdateRiver(ctx, riverChans[id], id)
					}
					riverMutex.Unlock()

					// Copy existing river channels to the new map
					newRiverChans[id] = riverChans[id]
				}

				// Acquire the mutex before updating riverChans
				riverMutex.Lock()

				// Update riverChans with the new map
				riverChans = newRiverChans

				// Release the mutex after updating riverChans
				riverMutex.Unlock()

				time.Sleep(5 * time.Second)
			}
		}
	}()

	for {
		// Acquire the mutex before accessing riverChans
		riverMutex.Lock()

		for _, riverChan := range riverChans {
			select {
			case r := <-riverChan:
				msg := UpdateRiver{
					Id:     r.Id,
					Height: r.Height,
					Status: r.Status,
				}

				// Acquire the mutex before updating the clients
				clientsMutex.Lock()
				for client := range clients {
					err := client.WriteJSON(msg)
					if err != nil {
						delete(clients, client)
						log.Printf("failed to write JSON to client: %v", err)
					}
				}
				// Release the mutex after updating the clients
				clientsMutex.Unlock()
			case <-ctx.Done():
				// Release the mutex before returning
				riverMutex.Unlock()
				return nil
			}
		}

		// Release the mutex before waiting for the next iteration
		riverMutex.Unlock()
	}
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
