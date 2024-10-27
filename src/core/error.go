package core

// no packages are imported here
// core is a package that contains the application's core logic
// core is imported in other packages

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Internal Success codes
const (
	SuccessCode = http.StatusOK      //200
	CreatedCode = http.StatusCreated //201
)

// Success messages
const (
	OrderEventSuccess = "Order Event Received Successfully"
)

// Internal error codes
const (
	BadRequestCode   = http.StatusBadRequest   //400
	UnauthorizedCode = http.StatusUnauthorized //401

	InternalServerErrorCode = http.StatusInternalServerError //500
)

// Error codes for the application 800 series
const (
	TransactionIDExistsErrorCode = 801
	InvalidJSONErrorCode         = 802
	DataBaseConnectionErrorCode  = 803
	DataBaseErrorCode            = 804
)

const (
	ErrUnauthorized          = "Unauthorized"
	ErrBadRequest            = "Bad request"
	ErrInternalServer        = "Internal server error"
	ErrInvalidJSON           = "Invalid JSON"
	ErrInvalidOrderRequest   = "Invalid order request"
	ErrOrderNotFound         = "Order not found"
	ErrOrderUpdateFailed     = "Order update failed"
	ErrTransactionIDExists   = "Transaction ID already exists"
	ErrTransactionIDNotFound = "Transaction ID not found"
	ErrEventCreationFailed   = "Event creation failed"
	ErrEventUpdateFailed     = "Event update failed"
	ErrEventProducerFailed   = "Event producer failed"
	ErrEventConsumerFailed   = "Event consumer failed"
	ErrEventNotFound         = "Event not found"
	ErrEventStatusUpdate     = "Event status update failed"
	ErrEventStageUpdate      = "Event stage update failed"
	ErrEventStatusInvalid    = "Invalid event status"
	ErrEventStageInvalid     = "Invalid event stage"
	ErrEventStatusExists     = "Event status already exists"
	ErrDatabaseConnection    = "Database connection error"
	ErrDatabase              = "Database error"
)

var (
	ErrUnauthorizedError          = errors.New("unauthorized")
	ErrBadRequestError            = errors.New("bad request")
	ErrInternalError              = errors.New("internal server error")
	ErrInvalidJSONError           = errors.New("invalid json")
	ErrOrderNotFoundError         = errors.New("order not found")
	ErrOrderUpdateError           = errors.New("order update failed")
	ErrTransactionIDError         = errors.New("transaction id already exists")
	ErrTransactionIDNotFoundError = errors.New("transaction id not found")
	ErrEventCreationError         = errors.New("event creation failed")
	ErrEventUpdateError           = errors.New("event update failed")
	ErrEventProducerError         = errors.New("event producer failed")
	ErrEventConsumerError         = errors.New("event consumer failed")
	ErrEventNotFoundError         = errors.New("event not found")
	ErrEventStatusUpdateError     = errors.New("event status update failed")
	ErrEventStageUpdateError      = errors.New("event stage update failed")
	ErrEventStatusInvalidError    = errors.New("invalid event status")
	ErrEventStageInvalidError     = errors.New("invalid event stage")
	ErrEventStatusExistsError     = errors.New("event status already exists")
	ErrDatabaseConnectionError    = errors.New("database connection error")
	ErrDatabaseError              = errors.New("database error")
)

// RestErr Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
	ErrBody() RestError
}

// RestError Rest error struct
type RestError struct {
	ErrStatus  int         `json:"status,omitempty"`
	ErrError   string      `json:"error,omitempty"`
	ErrMessage interface{} `json:"message,omitempty"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
}

// Status Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrMessage)
}

// Causes RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrMessage
}

// ErrBody Error body
func (e RestError) ErrBody() RestError {
	return e
}

func SuccessResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"status": true, "code": statusCode, "message": message})
}

// NewRestError constructs a new RestError instance
func NewRestError(status int, err error, causes interface{}, debug bool) RestError {
	restError := RestError{
		ErrStatus: status,
		ErrError:  err.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return restError
}

// BadRequestError constructs a c.JSON response for bad request errors
func BadRequestError(c *gin.Context, causes interface{}, debug bool) {
	restError := NewRestError(BadRequestCode, ErrBadRequestError, causes, debug)
	c.JSON(BadRequestCode, restError)
}

// UnauthorizedError constructs a c.JSON response for unauthorized errors
func UnauthorizedError(c *gin.Context, causes interface{}, debug bool) {
	restError := NewRestError(UnauthorizedCode, ErrUnauthorizedError, causes, debug)
	c.JSON(UnauthorizedCode, restError)
}

// InternalServerError constructs a c.JSON response for internal server errors
func InternalServerError(c *gin.Context, causes interface{}, debug bool) {
	restError := NewRestError(InternalServerErrorCode, ErrInternalError, causes, debug)
	c.JSON(InternalServerErrorCode, restError)
}

// Invalid JSON Error constructs a c.JSON response for invalid JSON errors
func InvalidJSONError(c *gin.Context, causes interface{}, debug bool) {
	restError := NewRestError(InvalidJSONErrorCode, ErrInvalidJSONError, causes, debug)
	c.JSON(InvalidJSONErrorCode, restError)
}

// TransactionIDExistsError constructs a c.JSON response for transaction ID exists errors
func TransactionIDExistsError(c *gin.Context, causes interface{}, debug bool) {
	restError := NewRestError(TransactionIDExistsErrorCode, ErrTransactionIDError, causes, debug)
	c.JSON(TransactionIDExistsErrorCode, restError)
}
