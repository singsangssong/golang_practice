package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	// - 초기 서버: 18080포트 번호 고정
	// err := http.ListenAndServe(
	// 	":18080",
	// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	// 	}),
	// )

	if err := run(context.Background()); err != nil { // 3. err에 반환값이 저장되면 이를 반환값 오류를 보냄.
		log.Printf("Failed to terminate server: %v", err)
	}

	s := http.Server{ // 2. context에 서버 중단 명령이 있으면, shutdown메서드로 http 서버 기능 종료
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}
	s.ListenAndServe() // 1. http 요청받기
}

func run(ctx context.Context) error {
	s := http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}
