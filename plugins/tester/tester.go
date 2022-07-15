package tester

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tester struct {
	mux  *runtime.ServeMux
	opts []grpc.DialOption
}

var instance *tester

func New() *tester {
	instance = &tester{
		mux:  runtime.NewServeMux(),
		opts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}
	return instance
}

// Instance can be called after run the constructor function New().
// It will return the tester instance.
// func Instance() *tester {
// 	return instance
// }

func (t *tester) Load() error {
	return nil
}

func (t *tester) Name() string {
	return "hrpc-tester"
}

func (t *tester) DependsOn() []string {
	return nil
}

func HandlerFromEndpoint() (*runtime.ServeMux, string, []grpc.DialOption) {
	return instance.mux, "localhost:8888", instance.opts
}

// Run should be called after pb.RegisterXXXXXHandlerFromEndpoint(ctx, )
func Run() error {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		if err := http.ListenAndServe(":8080", instance.mux); err != nil {
			fmt.Println(err.Error())
		}
	}()
	return nil
}
