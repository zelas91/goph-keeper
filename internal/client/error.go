package client

type ErrAuth struct {
	Err error
}

func (e ErrAuth) Error() string {
	return e.Err.Error()
}

//type ErrWorkingData struct {
//	Err error
//}
//
//func (e ErrWorkingData) Error() string {
//	return e.Err.Error()
//}
