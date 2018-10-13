package main

type DocumentConflictError struct {
	message string
}

func NewDocumentConflictError(message string) *DocumentConflictError {
	return &DocumentConflictError{message}
}

func (e *DocumentConflictError) Error() string {
	return e.message
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{message}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type OperationTimeoutError struct {
	message string
}

func NewOperationTimeoutError(message string) *OperationTimeoutError {
	return &OperationTimeoutError{message}
}

func (e *OperationTimeoutError) Error() string {
	return e.message
}
