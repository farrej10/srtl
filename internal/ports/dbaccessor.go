package ports

type IDatabaseAccessor interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	Close() error
}
