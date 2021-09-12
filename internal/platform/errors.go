package platform

const (
	ErrProductNotFound = constError("product not found")
)

type constError string

func (err constError) Error() string {
	return string(err)
}
