package core

import (
	"fmt"
	"game_hub/utils"
	"strings"
)

type ErrorCode string

const (
	Err             ErrorCode = "ERROR"
	ErrUnknown      ErrorCode = "UNKNOWN_ERROR"
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	ErrOutOfRange   ErrorCode = "OUT_OF_RANGE"
	ErrInvalidRange ErrorCode = "INVALID_RANGE"
	ErrEOF          ErrorCode = "END_OF_INPUT"
	ErrStateStack             = "STATE_STACK_ERROR"
	ErrLocalization           = "LOCALIZATION_ERROR"
	ErrCommand                = "COMMAND_ERROR"
	ErrInit         ErrorCode = "initialization_error"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
}

func (e *AppError) Error() string {
	msg := e.Message
	if e.Details == nil || len(e.Details) == 0 {
		return msg
	}
	msg += "\r\nDetails:"
	for key, value := range e.Details {
		msg += fmt.Sprintf("\r\n%s: %v", key, value)
	}
	return msg
}

func NewAppError(code ErrorCode, message string, details map[string]any) *AppError {
	return &AppError{Code: code, Message: message, Details: details}
}

type AppErrors struct {
	Errors []error
}

func NewAppErrors(errs []error) *AppErrors {
	appErrors := &AppErrors{
		Errors: make([]error, 0, 10),
	}
	if errs != nil {
		for _, err := range errs {
			appErrors.Add(err)
		}
	}
	return appErrors
}
func (e *AppErrors) Add(err error) {
	e.Errors = append(e.Errors, err)
}
func (e *AppErrors) Error() string {
	var text string
	for _, err := range e.Errors {
		text += (err.Error() + "\r\n")
	}
	return text
}

type ErrorHandler interface {
	Handle(err error) string
}

type LocalizedErrorHandler struct {
	localizer *MessageLocalizer
}

func NewLocalizedErrorHandler(localizer *MessageLocalizer) *LocalizedErrorHandler {
	return &LocalizedErrorHandler{
		localizer: localizer,
	}
}
func (h *LocalizedErrorHandler) Handle(err error) string {
	if err == nil {
		return ""
	}
	if appErrs, ok := err.(*AppErrors); ok {
		if len(appErrs.Errors) == 0 {
			return ""
		}
		var text string
		for i, e := range appErrs.Errors {
			text += (h.Handle(e))
			if i != (len(appErrs.Errors) - 1) {
				text += "\r\n"
			}
		}
		return text
	}
	if appErr, ok := err.(*AppError); ok {
		errMsg := appErr.Message
		var information string
		var informationKey string
		switch appErr.Code {
		case ErrInternal:
			informationKey = "internal_error"
		case ErrLocalization:
			informationKey = "localization_error"
		}
		if informationKey != "" {
			informationBuf, locErr := h.localizer.Get(informationKey)
			if locErr != nil {
				information = string(appErr.Code)
			}
			information = informationBuf
		}
		if appErr.Details == nil || appErr.Details["IsLocalized"] == nil {
			localizedMsg, locErr := h.localizer.Get(appErr.Message)
			if locErr != nil {
				failMsg, failMsgLocErr := h.localizer.Get("message_localization_failure")
				if failMsgLocErr != nil {
					return fmt.Sprintf("%v: failed to localize for message %s.\r\nError: %s\r\n", information, appErr.Message, err.Error())
				}
				return h.Handle(NewAppError(Err, failMsg, map[string]any{
					"IsLocalized": true,
					"message":     appErr.Message,
					"error":       err.Error(),
				}))
			}
			errMsg = localizedMsg
		}
		for key, value := range appErr.Details {
			if appErrValue, ok := value.(*AppError); ok {
				errMsg = strings.ReplaceAll(errMsg, "$"+key, h.Handle(appErrValue))
			}
			if appErrsValue, ok := value.(*AppErrors); ok {
				errMsg = strings.ReplaceAll(errMsg, "$"+key, h.Handle(appErrsValue))
			}
		}
		errMsg = utils.SubstituteParams(errMsg, appErr.Details)
		if information != "" {
			return fmt.Sprintf("%s: %s", information, errMsg)
		}
		return errMsg
	}
	appErr := NewAppError(ErrUnknown, fmt.Sprintf("Error: %v", err), map[string]any{"IsLocalized": true})
	return h.Handle(appErr)
}
