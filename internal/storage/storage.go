package storage

import "errors"

var (
	ErrUserExist                     = errors.New("user exist")
	ErrUserNotFound                  = errors.New("user not found")
	ErrModuleNotFound                = errors.New("module not found")
	ErrModuleCannotBeChanged         = errors.New("cannot change module")
	ErrModuleExist                   = errors.New("module exist")
	ErrRecommendationNotFound        = errors.New("recommendation not found")
	ErrRecommendationExist           = errors.New("recommendation exist")
	ErrRecommendationCannotBeChanged = errors.New("cannot change recommendation")
	ErrPlantFamilyNotFound           = errors.New("plant family not found")
	ErrPlantFamilyExist              = errors.New("plant family exist")
	ErrPlantFamilyCannotBeChanged    = errors.New("cannot change plant family")
	ErrSettingNotFound               = errors.New("setting not found")
	ErrSettingExist                  = errors.New("setting exist")
	ErrSettingCannotBeChanged        = errors.New("cannot change setting")
	ErrDataNotFound                  = errors.New("data not found")
)
