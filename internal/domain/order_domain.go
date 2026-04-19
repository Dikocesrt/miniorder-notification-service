package domain

type OrderKafkaMessage struct {
	OrderID       string `json:"order_id" validate:"required,uuid"`
	CustomerEmail string `json:"customer_email" validate:"required,email"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
}
