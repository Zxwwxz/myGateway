package balance

type LoadBalance interface {
	Add(...string) error
	Next(string) (string, error)
}
