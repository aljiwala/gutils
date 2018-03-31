package errorsx

import (
	"fmt"
	"net/http"
	"strings"
)

// DomainError represents an error thrown by the domain.
type DomainError struct {
	Status  int    `json:"status_code"`
	Message string `json:"error"`
}

// NewDomainError holds the value of status code, message and errors.
func NewDomainError(status int, message string) *DomainError {
	return &DomainError{status, message}
}

// Error should return the string value of the error.
func (err *DomainError) Error() string {
	return err.Message
}

// 1xx -------------------------------------------------------------------------

// Continue should return `http.StatusContinue` with custom message.
func Continue(message string) *DomainError { // 100
	return &DomainError{http.StatusContinue, message}
}

// SwitchingProtocols should return `http.StatusSwitchingProtocols` with custom
// message.
func SwitchingProtocols(message string) *DomainError { // 101
	return &DomainError{http.StatusSwitchingProtocols, message}
}

// Processing should return `http.StatusProcessing` with custom message.
func Processing(message string) *DomainError { // 102
	return &DomainError{http.StatusProcessing, message}
}

// 2xx -------------------------------------------------------------------------

// OK should return `http.StatusOK` with custom message.
func OK(message string) *DomainError { // 200
	return &DomainError{http.StatusOK, message}
}

// Created should return `http.StatusOK` with custom message.
func Created(message string) *DomainError { // 201
	return &DomainError{http.StatusCreated, message}
}

// Accepted should return `http.StatusAccepted` with custom message.
func Accepted(message string) *DomainError { // 202
	return &DomainError{http.StatusAccepted, message}
}

// NonAuthoritativeInfo should return `http.StatusNonAuthoritativeInfo`
// with custom message.
func NonAuthoritativeInfo(message string) *DomainError { // 203
	return &DomainError{http.StatusNonAuthoritativeInfo, message}
}

// NoContent should return `http.StatusNoContent` with custom message.
func NoContent(message string) *DomainError { // 204
	return &DomainError{http.StatusNoContent, message}
}

// ResetContent should return `http.StatusResetContent` with custom message.
func ResetContent(message string) *DomainError { // 205
	return &DomainError{http.StatusResetContent, message}
}

// PartialContent should return `http.StatusPartialContent` with custom message.
func PartialContent(message string) *DomainError { // 206
	return &DomainError{http.StatusPartialContent, message}
}

// MultiStatus should return `http.StatusMultiStatus` with custom message.
func MultiStatus(message string) *DomainError { // 207
	return &DomainError{http.StatusMultiStatus, message}
}

// AlreadyReported should return `http.StatusAlreadyReported` with custom message.
func AlreadyReported(message string) *DomainError { // 208
	return &DomainError{http.StatusAlreadyReported, message}
}

// IMUsed should return `http.StatusIMUsed` with custom message.
func IMUsed(message string) *DomainError { // 209
	return &DomainError{http.StatusIMUsed, message}
}

// 3xx -------------------------------------------------------------------------

// MultipleChoices should return `http.StatusMultipleChoices` with custom message.
func MultipleChoices(message string) *DomainError { // 300
	return &DomainError{http.StatusMultipleChoices, message}
}

// MovedPermanently should return `http.StatusMovedPermanently` with custom message.
func MovedPermanently(message string) *DomainError { // 301
	return &DomainError{http.StatusMovedPermanently, message}
}

// Found should return `http.StatusFound` with custom message.
func Found(message string) *DomainError { // 302
	return &DomainError{http.StatusFound, message}
}

// SeeOther should return `http.StatusSeeOther` with custom message.
func SeeOther(message string) *DomainError { // 303
	return &DomainError{http.StatusSeeOther, message}
}

// NotModified should return `http.StatusNotModified` with custom message.
func NotModified(message string) *DomainError { // 304
	return &DomainError{http.StatusNotModified, message}
}

// UseProxy should return `http.StatusUseProxy` with custom message.
func UseProxy(message string) *DomainError { // 305
	return &DomainError{http.StatusUseProxy, message}
}

// TemporaryRedirect should return `http.Status` with custom message.
func TemporaryRedirect(message string) *DomainError { // 307
	return &DomainError{http.StatusTemporaryRedirect, message}
}

// PermanentRedirect should return `http.Status` with custom message.
func PermanentRedirect(message string) *DomainError { // 308
	return &DomainError{http.StatusPermanentRedirect, message}
}

// 4xx -------------------------------------------------------------------------

// BadRequest should return `http.StatusBadRequest` with custom message.
func BadRequest(message string) *DomainError { // 400
	return &DomainError{http.StatusBadRequest, message}
}

// Unauthorized should return `http.StatusUnauthorized` with custom message.
func Unauthorized(message string) *DomainError { // 401
	return &DomainError{http.StatusUnauthorized, message}
}

// PaymentRequired should return `http.PaymentRequired` with custom message.
func PaymentRequired(message string) *DomainError { // 402
	return &DomainError{http.StatusPaymentRequired, message}
}

// Forbidden should return `http.StatusForbidden` with custom message.
func Forbidden(message string) *DomainError { // 403
	return &DomainError{http.StatusForbidden, message}
}

// NotFound should return `http.StatusNotFound` with custom message.
func NotFound(message string) *DomainError { // 404
	return &DomainError{http.StatusNotFound, message}
}

// MethodNotAllowed should return `http.Status` with custom message.
func MethodNotAllowed(message string) *DomainError { // 405
	return &DomainError{http.StatusMethodNotAllowed, message}
}

// NotAcceptable should return `http.StatusNotAcceptable` with custom message.
func NotAcceptable(message string) *DomainError { // 406
	return &DomainError{http.StatusNotAcceptable, message}
}

// ProxyAuthRequired should return `http.StatusProxyAuthRequired` with custom message.
func ProxyAuthRequired(message string) *DomainError { // 407
	return &DomainError{http.StatusProxyAuthRequired, message}
}

// RequestTimeout should return `http.Status` with custom message.
func RequestTimeout(message string) *DomainError { // 408
	return &DomainError{http.StatusRequestTimeout, message}
}

// Conflict should return `http.StatusConflict` with custom message.
func Conflict(message string) *DomainError { // 409
	return &DomainError{http.StatusConflict, message}
}

// Gone should return `http.Status` with custom message.
func Gone(message string) *DomainError { // 410
	return &DomainError{http.StatusGone, message}
}

// LengthRequired should return `http.Status` with custom message.
func LengthRequired(message string) *DomainError { // 411
	return &DomainError{http.StatusLengthRequired, message}
}

// PreconditionFailed should return `http.Status` with custom message.
func PreconditionFailed(message string) *DomainError { // 412
	return &DomainError{http.StatusPreconditionFailed, message}
}

// RequestEntityTooLarge should return `http.Status` with custom message.
func RequestEntityTooLarge(message string) *DomainError { // 413
	return &DomainError{http.StatusRequestEntityTooLarge, message}
}

// RequestURITooLong should return `http.Status` with custom message.
func RequestURITooLong(message string) *DomainError { // 414
	return &DomainError{http.StatusRequestURITooLong, message}
}

// UnsupportedMediaType should return `http.StatusUnsupportedMediaType` with
// custom message.
func UnsupportedMediaType(message string) *DomainError { // 415
	return &DomainError{http.StatusUnsupportedMediaType, message}
}

// RequestedRangeNotSatisfiable should return
// `http.StatusRequestedRangeNotSatisfiable` with custom message.
func RequestedRangeNotSatisfiable(message string) *DomainError { // 416
	return &DomainError{http.StatusRequestedRangeNotSatisfiable, message}
}

// ExpectationFailed should return `http.StatusExpectationFailed`
// with custom message.
func ExpectationFailed(message string) *DomainError { // 418
	return &DomainError{http.StatusExpectationFailed, message}
}

// Teapot should return `http.StatusTeapot` with custom message.
func Teapot(message string) *DomainError { // 418
	return &DomainError{http.StatusTeapot, message}
}

// UnprocessableEntity should return `http.StatusUnprocessableEntity` with
// custom message.
func UnprocessableEntity(message string) *DomainError { // 422
	return &DomainError{http.StatusUnprocessableEntity, message}
}

// Locked should return `http.StatusLocked` with custom message.
func Locked(message string) *DomainError { // 423
	return &DomainError{http.StatusLocked, message}
}

// FailedDependency should return `http.StatusFailedDependency`
// with custom message.
func FailedDependency(message string) *DomainError { // 424
	return &DomainError{http.StatusFailedDependency, message}
}

// UpgradeRequired should return `http.StatusUpgradeRequired`
// with custom message.
func UpgradeRequired(message string) *DomainError { // 426
	return &DomainError{http.StatusUpgradeRequired, message}
}

// PreconditionRequired should return `http.StatusPreconditionRequired`
// with custom message.
func PreconditionRequired(message string) *DomainError { // 428
	return &DomainError{http.StatusPreconditionRequired, message}
}

// TooManyRequests should return `http.Status` with custom message.
func TooManyRequests(message string) *DomainError { // 429
	return &DomainError{http.StatusTooManyRequests, message}
}

// RequestHeaderFieldsTooLarge should return
// `http.StatusRequestHeaderFieldsTooLarge` with custom message.
func RequestHeaderFieldsTooLarge(message string) *DomainError { // 431
	return &DomainError{http.StatusRequestHeaderFieldsTooLarge, message}
}

// UnavailableForLegalReasons should return
// `http.StatusUnavailableForLegalReasons` with custom message.
func UnavailableForLegalReasons(message string) *DomainError { // 451
	return &DomainError{http.StatusUnavailableForLegalReasons, message}
}

// 5xx -------------------------------------------------------------------------

// InternalServer should return `http.StatusInternalServerError`
// with custom message.
func InternalServer(message string) *DomainError { // 500
	return &DomainError{http.StatusInternalServerError, message}
}

// NotImplemented should return `http.StatusNotImplemented` with custom message.
func NotImplemented(message string) *DomainError { // 501
	return &DomainError{http.StatusNotImplemented, message}
}

// BadGateway should return `http.Status` with custom message.
func BadGateway(message string) *DomainError { // 502
	return &DomainError{http.StatusBadGateway, message}
}

// ServiceUnavailable should return `http.StatusServiceUnavailable`
// with custom message.
func ServiceUnavailable(message string) *DomainError { // 503
	return &DomainError{http.StatusServiceUnavailable, message}
}

// GatewayTimeout should return `http.StatusGatewayTimeout`
// with custom message.
func GatewayTimeout(message string) *DomainError { // 504
	return &DomainError{http.StatusGatewayTimeout, message}
}

// HTTPVersionNotSupported should return `http.StatusHTTPVersionNotSupported`
// with custom message.
func HTTPVersionNotSupported(message string) *DomainError { // 505
	return &DomainError{http.StatusHTTPVersionNotSupported, message}
}

// VariantAlsoNegotiates should return `http.StatusVariantAlsoNegotiates`
// with custom message.
func VariantAlsoNegotiates(message string) *DomainError { // 506
	return &DomainError{http.StatusVariantAlsoNegotiates, message}
}

// InsufficientStorage should return `http.StatusInsufficientStorage`
// with custom message.
func InsufficientStorage(message string) *DomainError { // 5
	return &DomainError{http.StatusInsufficientStorage, message}
}

// LoopDetected should return `http.StatusLoopDetected` with custom message.
func LoopDetected(message string) *DomainError { // 508
	return &DomainError{http.StatusLoopDetected, message}
}

// NotExtended should return `http.StatusNotExtended` with custom message.
func NotExtended(message string) *DomainError { // 510
	return &DomainError{http.StatusNotExtended, message}
}

// NetworkAuthenticationRequired should return
// `http.StatusNetworkAuthenticationRequired` with custom message.
func NetworkAuthenticationRequired(message string) *DomainError { // 511
	return &DomainError{http.StatusNetworkAuthenticationRequired, message}
}

// -----------------------------------------------------------------------------

// IsStatusNotModified should return true if HTTP status of an error is 204.
func (err *DomainError) IsStatusNotModified() bool {
	return err.Status == http.StatusNotModified
}

// IsStatusBadRequest should return true if HTTP status of an error is 400.
func (err *DomainError) IsStatusBadRequest() bool {
	return err.Status == http.StatusBadRequest
}

// IsStatusUnauthorized should return true if HTTP status of an error is 401.
func (err *DomainError) IsStatusUnauthorized() bool {
	return err.Status == http.StatusUnauthorized
}

// IsStatusNotFound should return true if HTTP status of an error is 404.
func (err *DomainError) IsStatusNotFound() bool {
	return err.Status == http.StatusNotFound
}

// IsStatusConflict should return true if HTTP status of an error is 409.
func (err *DomainError) IsStatusConflict() bool {
	return err.Status == http.StatusConflict
}

// -----------------------------------------------------------------------------

// NotUniqueTogether should return error for unique together fields.
func NotUniqueTogether(ss ...string) error {
	if len(ss) == 0 {
		return nil
	}

	return fmt.Errorf("'%s' must be unique", strings.Join(ss, ", "))
}

// NotUnique should return an error for unique together fields.
func NotUnique(s string) error {
	return NotUniqueTogether(s)
}

// -----------------------------------------------------------------------------
