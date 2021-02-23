package http

import (
	"fmt"
	"net/http"
	"time"
)

var ErrQuotaReached = fmt.Errorf("request quota reached")
var ErrTimeout = fmt.Errorf("timeout")

func NewRLClient(requestQuota int, failOnReachingLimit bool) *RLClient {
	requestQuotaChan := make(chan struct{}, requestQuota)

	for i := 0; i < requestQuota; i++ {
		requestQuotaChan <- struct{}{}
	}

	return &RLClient{
		requestQuota:        requestQuotaChan,
		failOnReachingLimit: failOnReachingLimit,
	}
}

type RLClient struct {
	requestQuota        chan struct{}
	failOnReachingLimit bool
}

func (c *RLClient) Get(url string) (resp *http.Response, err error) {
	if c.failOnReachingLimit {
		return c.getWithFailOnReachingLimit(url)
	}
	return c.getWithLockOnReachingLimit(url)
}

func (c *RLClient) getWithLockOnReachingLimit(url string) (resp *http.Response, err error) {
	select {
	case <-c.requestQuota:
		defer func() { c.requestQuota <- struct{}{} }()
		return http.Get(url)
	case <-time.After(10 * time.Second):
		return nil, ErrTimeout
	}
}

func (c *RLClient) getWithFailOnReachingLimit(url string) (resp *http.Response, err error) {
	select {
	case <-c.requestQuota:
		defer func() { c.requestQuota <- struct{}{} }()
		return http.Get(url)
	default:
		return nil, ErrQuotaReached
	}
}
