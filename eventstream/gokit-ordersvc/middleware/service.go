package middleware

import (
	ordersvc "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc"
)

// Middleware describes a service middleware.
type Middleware func(service ordersvc.Service) ordersvc.Service
