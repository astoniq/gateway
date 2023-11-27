package middleware

import (
	"github.com/astoniq/janus/pkg/observability"
	"github.com/gofrs/uuid"
	"net/http"
)

type regIdKeyType int

const (
	reqIdKey        regIdKeyType = iota
	requestIdHeader              = "X-Request-Id"
)

func RequestId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestId := request.Header.Get(requestIdHeader)
		if requestId == "" {
			requestId = uuid.Must(uuid.NewV4()).String()
		}
		request.Header.Set(requestIdHeader, requestId)
		writer.Header().Set(requestIdHeader, requestId)

		handler.ServeHTTP(writer, request.WithContext(observability.RequestIdToContext(request.Context(), requestId)))
	})
}
