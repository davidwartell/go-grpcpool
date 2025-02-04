package grpcpool

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func DecodeGrpcError(err error) (connectionFailure bool, returnErr error) {
	st, ok := status.FromError(err)
	if ok {
		// application specific error - Use st.Message() and st.Code()
		if st.Code() == codes.Canceled || st.Code() == codes.Unavailable {
			// this happens when server is gracefully shutting down with msg on client side:
			// https://github.com/grpc/grpc-go/blob/750abe8f95cd270ab68f7298a2854148d0c33030/clientconn.go grpc: the client connection is closing
			return true, err
		} else {
			rpcErr := errors.Errorf("error on grpc request %s (%d)", st.Message(), st.Code())
			return false, rpcErr
		}
	} else if errors.Is(err, io.EOF) {
		// stream returns EOF when gracefully shutting down
		return true, err
	} else if strings.Contains(err.Error(), "use of closed network connection") {
		return true, err
	} else if strings.Contains(err.Error(), "Error while dialing") {
		return true, err
	} else if strings.Contains(err.Error(), "context canceled") {
		return true, err
	} else {
		// general GRPC error
		return true, err
	}
}

//goland:noinspection GoUnusedExportedFunction
func GrpcErrIsConnFail(err error) (connectionFailure bool) {
	connectionFailure, _ = DecodeGrpcError(err)
	return
}

//goland:noinspection GoUnusedExportedFunction
func DecodeGrpcErrorCode(err error) (hasCode bool, code codes.Code) {
	st, ok := status.FromError(err)
	if ok {
		hasCode = true
		code = st.Code()
	}
	return
}
