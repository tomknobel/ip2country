package rate

type ILimiter interface {
	Allow() bool
}
