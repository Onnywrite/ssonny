package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

func RunSuitsParallel(wg *sync.WaitGroup, t *testing.T, ss ...suite.TestingSuite) {
	for i := range ss {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			suite.Run(t, ss[ii])
		}(i)
	}
}

func RunSuitsSync(ctx context.Context, t *testing.T, ss ...suite.TestingSuite) error {
	for i := range ss {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		default:
			suite.Run(t, ss[i])
		}
	}
	return nil
}
