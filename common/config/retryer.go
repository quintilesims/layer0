package config

import (
	"math/rand"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
)

type Retryer struct {
	c            client.DefaultRetryer
	minRetryTime int
	maxRetryTime int
}

func NewRetryer(minRetryTime, maxRetryTime int, a *aws.Config) *Retryer {
	// minRetryTime default is 500ms
	// maxRetryTime default is 5 minutes
	retryer := &Retryer{
		minRetryTime: 500,
		maxRetryTime: 300000,
	}

	if minRetryTime > 0 {
		retryer.minRetryTime = minRetryTime
	}

	if maxRetryTime > 0 {
		retryer.maxRetryTime = maxRetryTime
	}

	return retryer
}

var seededRand = rand.New(&lockedSource{src: rand.NewSource(time.Now().UnixNano())})

// RetryRules returns the delay duration before retrying this request again
func (d Retryer) RetryRules(r *request.Request) time.Duration {
	minTime := 30
	throttle := d.shouldThrottle(r)
	if throttle {
		minTime = 500
	}

	retryCount := r.RetryCount
	if retryCount > 13 {
		retryCount = 13
	} else if throttle && retryCount > 8 {
		retryCount = 8
	}

	delay := (1 << uint(retryCount)) * (seededRand.Intn(minTime) + minTime)
	if delay < d.minRetryTime {
		delay = d.minRetryTime
	}

	if delay > d.maxRetryTime {
		delay = d.maxRetryTime
	}

	return time.Duration(delay) * time.Millisecond
}

// ShouldRetry returns true if the request should be retried.
func (d Retryer) ShouldRetry(r *request.Request) bool {
	// If one of the other handlers already set the retry state
	// we don't want to override it based on the service's state
	if r.Retryable != nil {
		return *r.Retryable
	}

	if r.HTTPResponse.StatusCode >= 500 {
		return true
	}
	return r.IsErrorRetryable() || d.shouldThrottle(r)
}

// ShouldThrottle returns true if the request should be throttled.
func (d Retryer) shouldThrottle(r *request.Request) bool {
	if r.HTTPResponse.StatusCode == 502 ||
		r.HTTPResponse.StatusCode == 503 ||
		r.HTTPResponse.StatusCode == 504 {
		return true
	}
	return r.IsErrorThrottle()
}

// lockedSource is a thread-safe implementation of rand.Source
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
