package types

type ConnectionString struct {
	Path string
}

type App struct {
	ServerAddress   string
	BackEndAddress  string
	FrontEndAddress string
	Debug           bool
}
