package repository

type Store interface {
	SaveLogInCash() error
	GetCash() Cash
	AddLog(log log) error
}