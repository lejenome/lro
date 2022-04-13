package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	ValidationErrorCode = iota + 100
)

type JsonApiErrorSource struct {
	Pointer  string `json:"pointer,omitempty"`
	Paramter string `json:"paramter,omitempty"`
}
type JsonApiError struct {
	ID     uint                `json:"id,omitempty"`
	Code   uint                `json:"code,omitempty"`
	Status uint                `json:"status"`
	Title  string              `json:"title"`
	Detail string              `json:"detail,omitempty"`
	Source *JsonApiErrorSource `json:"source,omitempty"`
	Meta   interface{}         `json:"meta,omitempty"`
}

type Links struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}
type Meta map[string]interface{}
type JsonApiData struct {
	Data  interface{} `json:"data"`
	Links *Links      `json:"links,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

func AbortValidatorError(c *gin.Context, err error) {
	type ValidatorError struct {
		Tag       string
		Namespace string
		Kind      string
		Field     string
		Type      string
		Value     string
		Param     string
	}
	errs := err.(validator.ValidationErrors)
	res := make([]JsonApiError, len(errs))
	for i, err := range errs {
		res[i] = JsonApiError{
			Status: 422,
			Title:  "Field Validation Error",
			Code:   ValidationErrorCode,
			Detail: "Given request body was invalid.",
			Source: &JsonApiErrorSource{
				Pointer: err.Field(),
			},
			Meta: ValidatorError{
				Tag:       err.ActualTag(),
				Namespace: err.Namespace(),
				Field:     err.Field(),
				Kind:      err.Kind().String(),
				Type:      err.Type().String(),
				Value:     fmt.Sprintf("%v", err.Value()),
				Param:     err.Param(),
			},
		}
	}
	ErrorResponse(c, 422, res)
	c.Abort()
}

func ErrorResponse(c *gin.Context, status int, msg interface{}) {
	var errs []JsonApiError
	switch msg.(type) {
	case []JsonApiError:
		errs = msg.([]JsonApiError)
	case JsonApiError:
		errs = append(errs, msg.(JsonApiError))
	case *JsonApiError:
		errs = append(errs, *msg.(*JsonApiError))
	case error:
		errs = append(errs, JsonApiError{
			Title:  msg.(error).Error(),
			Status: uint(status),
		})
	case *error:
		errs = append(errs, JsonApiError{
			Title:  (*msg.(*error)).Error(),
			Status: uint(status),
		})
	case string:
		errs = append(errs, JsonApiError{
			Title:  msg.(string),
			Status: uint(status),
		})
	default:
		errs = append(errs, JsonApiError{
			Title:  fmt.Sprintf("%+v", msg),
			Status: uint(status),
		})
	}
	res := struct {
		Errors []JsonApiError `json:"errors"`
	}{Errors: errs}
	c.JSON(status, &res)
}

func DataResponse(c *gin.Context, status int, msg interface{}) {
	switch {
	case status == http.StatusNoContent:
		c.JSON(http.StatusNoContent, nil)
	/*
		case status == http.StatusCreated:
			c.JSON(http.StatusNoContent, nil)
	*/
	default:
		c.JSON(status, &JsonApiData{
			Data: msg,
		})
	}

}
