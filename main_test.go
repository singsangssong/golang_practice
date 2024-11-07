package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// 초기 서버의 문제점
// 1. 출력 검증 어려움
// 2. 테스트 완료 후 종료방법 x
// 3. 이상 처리시 서버가 꺼짐. (os.Exit)
// 4. 포트번호 고정으로 테스트에서 서버실행 안될수도 있음.
// func TestMainFunc(t *testing.T) {
// 	go main()
// }

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
