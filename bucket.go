package limiter

type Bucket interface {
	Take() bool
}
