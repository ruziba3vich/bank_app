package entrepreneur

import (
	"context"
	"errors"

	"github.com/google/uuid"

	domain "github.com/prodonik/bank_app/internal/domain/entrepreneur"
	domainn "github.com/prodonik/bank_app/internal/domain/inn"
)

type Service struct {
	repo    domain.Repository
	innRepo domainn.Repository
}

func NewService(repo domain.Repository, innRepo domainn.Repository) *Service {
	return &Service{repo: repo, innRepo: innRepo}
}

type CreateInput struct {
	InnName               string
	LegalName             string
	RegistrationAuthority string
	RegistrationDate      string
	RegistrationNumber    string
	LegalForm             string
	IfutCode              int32
	DbibtCode             int32
	ActivityStatus        bool
	CharterFund           int32
	Founders              string
	Email                 string
	Phone                 string
	MhobtCode             string
	Address               string
	DirectorName          string
}

type UpdateInput struct {
	LegalName             *string
	RegistrationAuthority *string
	RegistrationDate      *string
	RegistrationNumber    *string
	LegalForm             *string
	IfutCode              *int32
	DbibtCode             *int32
	ActivityStatus        *bool
	CharterFund           *int32
	Founders              *string
	Email                 *string
	Phone                 *string
	MhobtCode             *string
	Address               *string
	DirectorName          *string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*domain.Entrepreneur, error) {
	// Try to find existing INN by name, create if not found
	innRecord, err := s.innRepo.GetByName(ctx, input.InnName)
	if err != nil {
		if !errors.Is(err, domainn.ErrInnNotFound) {
			return nil, err
		}
		innRecord, err = s.innRepo.Create(ctx, &domainn.Inn{Name: input.InnName})
		if err != nil {
			return nil, err
		}
	}

	e := &domain.Entrepreneur{
		InnID:                 innRecord.ID,
		InnName:               innRecord.Name,
		LegalName:             input.LegalName,
		RegistrationAuthority: input.RegistrationAuthority,
		RegistrationDate:      input.RegistrationDate,
		RegistrationNumber:    input.RegistrationNumber,
		LegalForm:             input.LegalForm,
		IfutCode:              input.IfutCode,
		DbibtCode:             input.DbibtCode,
		ActivityStatus:        input.ActivityStatus,
		CharterFund:           input.CharterFund,
		Founders:              input.Founders,
		Email:                 input.Email,
		Phone:                 input.Phone,
		MhobtCode:             input.MhobtCode,
		Address:               input.Address,
		DirectorName:          input.DirectorName,
	}

	return s.repo.Create(ctx, e)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Entrepreneur, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAll(ctx context.Context, filter domain.EntrepreneurFilter) ([]*domain.Entrepreneur, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*domain.Entrepreneur, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.LegalName != nil {
		existing.LegalName = *input.LegalName
	}
	if input.RegistrationAuthority != nil {
		existing.RegistrationAuthority = *input.RegistrationAuthority
	}
	if input.RegistrationDate != nil {
		existing.RegistrationDate = *input.RegistrationDate
	}
	if input.RegistrationNumber != nil {
		existing.RegistrationNumber = *input.RegistrationNumber
	}
	if input.LegalForm != nil {
		existing.LegalForm = *input.LegalForm
	}
	if input.IfutCode != nil {
		existing.IfutCode = *input.IfutCode
	}
	if input.DbibtCode != nil {
		existing.DbibtCode = *input.DbibtCode
	}
	if input.ActivityStatus != nil {
		existing.ActivityStatus = *input.ActivityStatus
	}
	if input.CharterFund != nil {
		existing.CharterFund = *input.CharterFund
	}
	if input.Founders != nil {
		existing.Founders = *input.Founders
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}
	if input.Phone != nil {
		existing.Phone = *input.Phone
	}
	if input.MhobtCode != nil {
		existing.MhobtCode = *input.MhobtCode
	}
	if input.Address != nil {
		existing.Address = *input.Address
	}
	if input.DirectorName != nil {
		existing.DirectorName = *input.DirectorName
	}

	return s.repo.Update(ctx, existing)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
