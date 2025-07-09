package applications

type Exception interface {
	GetCode() int
	GetError() error
	GetMessage() string
}

func NewApplicationException(code int, err error) *ApplicationException {
	return &ApplicationException{Code: code, Err: err}
}

type ApplicationException struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

func (e ApplicationException) GetCode() int {
	return e.Code
}

func (e ApplicationException) GetError() error {
	return e.Err
}

func (e ApplicationException) GetMessage() string {
	return e.Err.Error()
}
