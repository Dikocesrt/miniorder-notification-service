package common

import "github.com/gofiber/fiber/v2"

type Response[T any] struct {
	Body   T
	Status int
	Header map[string]string
	Cookie []fiber.Cookie
}

type ResponseBodySuccess[T any] struct {
	Message string `json:"message,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

type ResponseBodySuccessList[T any] struct {
	Message string `json:"message,omitempty"`
	Data    []T    `json:"data"`
}

type ResponseBodyPaginated[T any] struct {
	Message  string                    `json:"message,omitempty"`
	Metadata ResponseMetadataPaginated `json:"metadata,omitempty"`
	Data     []T                       `json:"data,omitempty"`
}

type ResponseMetadataPaginated struct {
	Count      int `json:"count,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
}

type ResponseBodyError struct {
	Message string `json:"message,omitempty"`
}
