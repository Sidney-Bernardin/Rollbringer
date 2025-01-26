package service

type GamesDatabaseRepository interface {
	Close() error
}
