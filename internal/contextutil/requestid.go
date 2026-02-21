package contextutil

import "context"

type ctxKey string

const ctxKeyRequestID ctxKey = "request_id"

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, id)
}

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyRequestID).(string)
	return id
}
