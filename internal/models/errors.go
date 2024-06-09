package models

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateUser      = errors.New("user is already exist")
	ErrNilUUID            = errors.New("uuid id nil")
	ErrNotAllowed         = errors.New("not allowed")
	ErrDuplicateCompany   = errors.New("duplicate company")
	ErrCompanyNotFound    = errors.New("company not found")
	ErrCompanyNameIsEmpty = errors.New("name is empty")
	ErrUserNameIsEmpty    = errors.New("username is empty")
	ErrEmailIsEmpty       = errors.New("email is empty")
	ErrPhoneIsEmpty       = errors.New("phone is empty")
)
