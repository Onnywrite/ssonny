package repo

type Transactor interface {
	Rollback() error
	Commit() error
}
