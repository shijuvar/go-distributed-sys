package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	ordersvc "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc"
	svctransport "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc/transport"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// NewService wires Go kit endpoints to the HTTP transport
func NewService(svcEndpoints svctransport.Endpoints, options []kithttp.ServerOption, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	/*
		   options := []kithttp.ServerOption{
				kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
				kithttp.ServerErrorEncoder(encodeError),
			}
	*/
	errorLogger := kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger))
	errorEncoder := kithttp.ServerErrorEncoder(encodeError)
	options = append(options, errorLogger, errorEncoder)

	// Configure routes with Gorilla Mux package
	r.Methods("POST").Path("/api/orders").Handler(kithttp.NewServer(
		svcEndpoints.CreateOrder,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/api/orders/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetOrderByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req svctransport.CreateOrderRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Order); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return svctransport.GetOrderByIDRequest{ID: id}, nil
}

// Errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type Errorer interface {
	Error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Errorer); ok && e.Error() != nil { // Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
func codeFrom(err error) int {
	switch err {
	case ordersvc.ErrNotFound:
		return http.StatusNotFound
	//case ErrAlreadyExists, ErrInconsistentIDs:
	//	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
