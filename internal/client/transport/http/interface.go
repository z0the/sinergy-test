package http

type Controller interface {
	Run(port string) error
}
