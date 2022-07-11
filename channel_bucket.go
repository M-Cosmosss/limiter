package limiter

import (
	"time"

	"github.com/flamego/flamego"
	"github.com/pkg/errors"
)

var buckets = map[string]*ChannelBucket{}

type ChannelBucket struct {
	ch       chan struct{}
	interval time.Duration
}

func NewChannelBucket(rate int, capacity int) *ChannelBucket {
	ch := make(chan struct{}, capacity)
	go func() {
		interval := time.Second / time.Duration(rate)
		for {
			ch <- struct{}{}
			time.Sleep(interval)
		}
	}()
	return &ChannelBucket{
		ch: ch,
	}
}

func (b *ChannelBucket) Take() bool {
	select {
	case <-b.ch:
		return true
	default:
		return false
	}
}

type ChannelBucketOption struct {
	Rate     int
	Capacity int
	Key      func(c flamego.Context) string
}

func Limiter(opt ChannelBucketOption) flamego.Handler {
	return func(c flamego.Context) error {
		var b *ChannelBucket
		var ok bool
		key := opt.Key(c)
		if b, ok = buckets[key]; !ok {
			b = NewChannelBucket(opt.Rate, opt.Capacity)
			buckets[key] = b
		}
		if !b.Take() {
			return errors.New("too many requests")
		}
		return nil
	}
}
