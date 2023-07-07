package agent

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"strings"
)

// service is the interface that wraps the basic methods.
type service interface {
	Invoke(ctx context.Context, method string, content []byte) ([]byte, error)
}

const (
	GRPC uint8 = iota
)

func Invoke(typ uint8, ctx context.Context, method string, content []byte) ([]byte, error) {
	switch typ {
	case GRPC:
		return grpcServiceInstance.Invoke(ctx, method, content)
	default:
		return nil, fmt.Errorf("unsupported type: %d", typ)
	}
}

type grpcService struct {
}

var (
	serviceClientMap = make(map[string]*grpc.ClientConn)
	hackAddressMap   = map[string]string{
		"account":     "127.0.0.1:6200",
		"matchmaking": "127.0.0.1:6202",
	}
	grpcServiceInstance = &grpcService{}
)

// Invoke is a method that invokes the grpc service.
//
//	method: the grpc method name, format: /{package}.{service}/{rpc_method}
func (gs *grpcService) Invoke(ctx context.Context, method string, content []byte) (reply []byte, err error) {
	packageName := getPackageName(method)
	if conn, ok := serviceClientMap[packageName]; ok {
		err = conn.Invoke(ctx, method, content, reply)
	} else {
		address, ok := hackAddressMap[packageName]
		if !ok {
			return nil, fmt.Errorf("unsupported service: %s", packageName)
		}

		conn, err = grpc.DialContext(ctx, address, grpc.WithBlock())
		if err != nil {
			return
		}
		serviceClientMap[packageName] = conn
		err = conn.Invoke(ctx, method, content, reply)
	}
	return
}

func getPackageName(method string) string {
	index := strings.Index(method, "/")
	if index == -1 {
		return ""
	}
	dotIndex := strings.Index(method, ".")
	if dotIndex == -1 || dotIndex < index {
		return ""
	}
	return method[index+1 : dotIndex]
}
