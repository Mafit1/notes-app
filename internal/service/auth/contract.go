package auth

type Service interface {
	Register(in RegisterIn) (RegisterOut, error)
	Login(in LoginIn) (LoginOut, error)
}
