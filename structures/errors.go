package Structures

import "errors"

var ErrCredentialsMissing = errors.New("Authorization credentials missing.")
var ErrDeveloperOnly = errors.New("This endpoint is only for developers.")
