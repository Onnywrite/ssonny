package tests_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite1 struct {
	suite.Suite
}

func (s *TestSuite1) Test1() {
	time.Sleep(time.Second)
}

func (s *TestSuite1) Test2() {
	time.Sleep(time.Second)
}

type TestSuite2 struct {
	suite.Suite
}

func (s *TestSuite2) Test1() {
	time.Sleep(time.Second)
}

func (s *TestSuite2) Test2() {
	time.Sleep(time.Second)
}

type TestSuite3 struct {
	suite.Suite
}

func (s *TestSuite3) Test1() {
	time.Sleep(time.Second)
}

func (s *TestSuite3) Test2() {
	time.Sleep(time.Second)
}

type TestSuite4 struct {
	suite.Suite
}

func (s *TestSuite4) Test1() {
	time.Sleep(time.Second)
}

func (s *TestSuite4) Test2() {
	time.Sleep(time.Second)
}

type TestSuite5 struct {
	suite.Suite
}

func (s *TestSuite5) Test1() {
	time.Sleep(time.Second)
}

func (s *TestSuite5) Test2() {
	time.Sleep(time.Second)
}

type TestSuite6 struct {
	suite.Suite
}

func (s *TestSuite6) Test1() {
	time.Sleep(time.Second)
}

func (s *TestSuite6) Test2() {
	time.Sleep(time.Second)
}

func (s *TestSuite6) Test3() {
	time.Sleep(time.Second)
}

func (s *TestSuite6) Test4() {
	time.Sleep(time.Second)
}

func TestParallelism(t *testing.T) {
	wg := sync.WaitGroup{}
	start := time.Now()
	tests.RunSuitsParallel(t, &wg, new(TestSuite1), new(TestSuite2), new(TestSuite3), new(TestSuite4), new(TestSuite5), new(TestSuite6))
	wg.Wait()
	end := time.Now()
	require.LessOrEqual(t, end.Sub(start), time.Second*5)
}

func TestNotParallelism(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*5)
	defer c()
	err := tests.RunSuitsSync(ctx, t, new(TestSuite1), new(TestSuite2), new(TestSuite3), new(TestSuite4), new(TestSuite5), new(TestSuite6))
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
