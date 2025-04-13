package core

import (
	"fmt"
	"game_hub/utils"
)

type ErrorCode string

const (
	ErrUnknown      ErrorCode = "UNKNOWN_ERROR"
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	ErrOutOfRange   ErrorCode = "OUT_OF_RANGE"
	ErrInvalidRange ErrorCode = "INVALID_RANGE"
	ErrEOF          ErrorCode = "END_OF_INPUT"
	ErrStateStack             = "STATE_STACK_ERROR"
	ErrLocalization           = "LOCALIZATION_ERROR"
	ErrCommand                = "COMMAND_ERROR"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]any
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code ErrorCode, message string, details map[string]any) *AppError {
	return &AppError{Code: code, Message: message, Details: details}
}

// интерфейс обработчика ошибок
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
	if appErr, ok := err.(*AppError); ok {
		errMsg := appErr.Message
		if appErr.Details == nil || appErr.Details["IsLocalized"] == nil {
			msg, locErr := h.localizer.Get(appErr.Message)
			if locErr != nil {
				if appErr.Code == ErrLocalization {
					return fmt.Sprintf("%v: failed to localize for message %s.\r\nDetails: %v.\r\n", appErr.Code, appErr.Message, appErr.Details)
				}
				return h.Handle(locErr)
			}
			errMsg = msg
		}
		errMsg = utils.SubstituteParams(errMsg, appErr.Details)
		var informationKey string
		switch appErr.Code {
		case ErrInternal:
			informationKey = "internal_error"
		case ErrLocalization:
			informationKey = "localization_error"
		}
		if informationKey != "" {
			information, locErr := h.localizer.Get(informationKey)
			if locErr != nil {
				information = string(appErr.Code)
			}
			return fmt.Sprintf("%s: %s", information, errMsg)
		}
		return errMsg
	}
	appErr := NewAppError(ErrUnknown, "error", map[string]any{"error": fmt.Sprintf("%v", err)})
	return h.Handle(appErr)
}
