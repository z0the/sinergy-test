package tcp

type Controller interface {
	Run(port string) error
}
