package contextx

import "context"

func NewContextWithValue(key, value string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, key, value)
}

func GetValueFromContext(ctx context.Context, key string) string {
	value := ctx.Value(key)
	if value == nil {
		return ""
	}
	return value.(string)
}
