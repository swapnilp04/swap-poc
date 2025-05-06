package swapErr

import "errors"

var ErrInvalidUser = errors.New("invalid username or password")
var ErrInternalServer = errors.New("internal server error")
var ErrPasswordMisMatch = errors.New("password and confirm password does match")
var ErrBadData = errors.New("bad data")
var ErrForbidden = errors.New("forbidden")
var ErrAlreadyChecked = errors.New("Already Checked")
var ErrAlreadyHasClass = errors.New("Already Assigned to Class")
var ErrRemoveHosteStudents = errors.New("Remove Hostel Student")
var ErrEmptyRole = errors.New("Empty Role")
