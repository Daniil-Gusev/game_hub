package core

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"os"
)

func EncodeData(data any) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, NewAppError(Err, "encode_error", map[string]any{
			"error": err,
		})
	}
	return bytes, nil
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, NewAppError(Err, "file_open_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, NewAppError(Err, "file_read_error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	return data, nil
}

func DecodeData(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return NewAppError(Err, "decode_error", map[string]any{
			"error": err,
		})
	}
	if err := ValidateData(target); err != nil {
		return NewAppError(Err, "decode_error", map[string]any{"error": err})
	}
	return nil
}

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateData(data any) error {
	err := validate.Struct(data)
	if err != nil {
		return err
	}
	return nil
}

func LoadData(filePath string, target any) error {
	data, err := ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := DecodeData(data, target); err != nil {
		return err
	}
	if err := ValidateData(target); err != nil {
		return err
	}
	return nil
}
