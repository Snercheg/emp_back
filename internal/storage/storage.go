package storage

import "errors"

var (
	ErrUserExist              = errors.New("user exist")
	ErrUserNotFound           = errors.New("user not found")
	ErrAppNotFound            = errors.New("app not found")
	ErrModuleNotFound         = errors.New("module not found")
	ErrRecommendationNotFound = errors.New("recommendation not found")
	ErrRecommendationExist    = errors.New("recommendation exist")
	ErrPlantFamilyNotFound    = errors.New("plant family not found")
	ErrPlantFamilyExist       = errors.New("plant family exist")
	ErrSettingNotFound        = errors.New("setting not found")
	ErrDataNotFound           = errors.New("data not found")
)
