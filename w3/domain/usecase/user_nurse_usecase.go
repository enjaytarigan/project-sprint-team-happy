package usecase

import (
	"halosuster/domain/entity"
	"halosuster/domain/repository"
	"halosuster/internal/helper"
	"strconv"
)

type IUserNurseUsecase interface {
	Create(in CreateNurseRequest, createdBy string) (entity.User, error)
	Update() error
	Delete(nurseId string) error
	SetAccess(nurseId string) error
}

type userNurseUseCase struct {
	userRepository repository.IUserRepository
}

func NewUserNurseUseCase(userRepo repository.IUserRepository) *userNurseUseCase {
	return &userNurseUseCase{
		userRepository: userRepo,
	}
}

type CreateNurseRequest struct {
	NIP                 int    `json:"nip" validate:"required,numeric,min=3030000000000,max=3039999999999"`
	Name                string `json:"name" validate:"required,min=5,max=50"`
	IndetityCardScanImg string `json:"identityCardScanImg" validate:"required,url"`
}

func (u *userNurseUseCase) Create(in CreateNurseRequest, createdBy string) (entity.User, error) {
	userNip := strconv.FormatInt(int64(in.NIP), 10)
	newNurse := entity.User{
		NIP:                 userNip,
		Name:                in.Name,
		IdentityCardScanImg: in.IndetityCardScanImg,
		Role:                string(entity.NURSE),
	}

	if !entity.ValidateUserNIP(userNip, entity.NURSE) {
		return entity.User{}, helper.CustomError{
			Message: "nip is invalid",
			Code:    400,
		}
	}

	if !entity.ValidateIdentityCardScanImageURL(newNurse.IdentityCardScanImg) {
		return entity.User{}, helper.CustomError{
			Message: "identityCardScanImgUrl is invalid",
			Code:    400,
		}
	}

	isNIPExist, err := u.userRepository.CheckNIPExist(newNurse.NIP)

	if err != nil {
		return entity.User{}, err
	}

	if isNIPExist {
		return entity.User{}, helper.CustomError{
			Message: "nip is already used",
			Code:    409,
		}
	}

	if createdNurse, err := u.userRepository.InsertUser(newNurse); err != nil {
		return entity.User{}, err
	} else {
		newNurse.ID = createdNurse.ID
	}

	return newNurse, nil
}

func (u *userNurseUseCase) Update() error {
	return nil
}

func (u *userNurseUseCase) Delete(nurseId string) error {
	return nil
}

func (u *userNurseUseCase) SetAccess(nurseId string) error {
	return nil
}
