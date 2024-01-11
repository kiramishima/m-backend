package errors

import "errors"

// Repository Errors
var (
	ErrExecuteQuery      = errors.New("failed to execute query")
	ErrScanData          = errors.New("failed to scan data")
	ErrPrepareStatement  = errors.New("failed to prepare SQL statement")
	ErrBeginTransaction  = errors.New("failed to begin transaction")
	ErrRollback          = errors.New("failed to rollback transaction")
	ErrCommit            = errors.New("failed to commit transaction")
	ErrRetrieveRows      = errors.New("failed to retrieve rows affected")
	ErrAlreadyExists     = errors.New("email already exists")
	ErrNoRecords         = errors.New("not records")
	ErrItemNotFound      = errors.New("item don't exist")
	ErrUpdatingRecord    = errors.New("failed to update record")
	ErrExecuteStatement  = errors.New("failed to execute statement")
	ErrBondAlreadyExists = errors.New("bond already exists")
	ErrBondNotExist      = errors.New("bond doesn't exist")
	ErrDeleteBond        = errors.New("failed deleting the bond")
	ErrNoAvailableBonds  = errors.New("requested num of bonds no available")
)
