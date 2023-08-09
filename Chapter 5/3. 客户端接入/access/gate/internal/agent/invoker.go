package agent

import (
	"ProjectX/service/matchmaking/pb/matchmaking"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	var msg, repMsg proto.Message

	switch method {
	case matchmaking.Matchmaking_Start_FullMethodName:
		msg = &matchmaking.StartRequest{}
		repMsg = &matchmaking.StartResponse{}
		break
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	err = proto.Unmarshal(content, msg)
	if err != nil {
		return
	}

	packageName := getPackageName(method)
	if conn, ok := serviceClientMap[packageName]; ok {
		err = conn.Invoke(ctx, method, msg, repMsg)
	} else {
		address, ok := hackAddressMap[packageName]
		if !ok {
			return nil, fmt.Errorf("unsupported service: %s", packageName)
		}

		conn, err = grpc.DialContext(ctx, address, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return
		}
		serviceClientMap[packageName] = conn
		err = conn.Invoke(ctx, method, msg, repMsg)
	}

	if err == nil {
		reply, err = proto.Marshal(repMsg)
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
