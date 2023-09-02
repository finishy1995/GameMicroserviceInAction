package contextx

import (
	"context"
	"google.golang.org/grpc/metadata"
)

var (
	defaultMetadata = metadata.MD{}
)

func NewContextWithValue(key, value string) context.Context {
	ctx := metadata.NewOutgoingContext(context.Background(), defaultMetadata)

	return SetKeyValue(ctx, key, value)
}

func getMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if md, ok = metadata.FromOutgoingContext(ctx); !ok {
			return metadata.MD{}
		}
	}

	return md
}

func SetKeyValue(ctx context.Context, key string, value string) context.Context {
	md := getMetadata(ctx)
	md.Set(key, value)
	return metadata.NewOutgoingContext(ctx, md)
}

func GetValueFromContext(ctx context.Context, key string) string {
	values := getValues(ctx, key)

	if size := len(values); size != 0 {
		return values[size-1]
	}

	return ""
}

func getValues(ctx context.Context, key string) []string {
	md := getMetadata(ctx)
	values := md.Get(key)
	return values
}
