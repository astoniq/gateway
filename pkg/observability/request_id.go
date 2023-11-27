package observability

import "context"

type reqIdKeyType int

const reqIdKey reqIdKeyType = iota

func RequestIdToContext(ctx context.Context, requestId string) context.Context {
	if ctx == nil {
		panic("Can not put request Id to empty context")
	}
	return context.WithValue(ctx, reqIdKey, requestId)
}

func RequestIdFromContext(ctx context.Context) string {
	if ctx == nil {
		panic("Can not get request Id from empty context")
	}

	if requestId, ok := ctx.Value(reqIdKey).(string); ok {
		return requestId
	}
	return ""
}
