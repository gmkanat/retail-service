package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	log.Printf("Request - Method: %s; Request: %v", info.FullMethod, req)

	resp, err := handler(ctx, req)

	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Response - Method: %s; Error: %v", info.FullMethod, st.Message())
	} else {
		log.Printf("Response - Method: %s; Response: %v", info.FullMethod, resp)
	}

	duration := time.Since(start)
	log.Printf("Duration: %s", duration)

	return resp, err
}
