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

type SheetCredential struct {
	SheetID      string
	Scope        string
	Type         string
	ProjectID    string
	PrivateKeyID string
	PrivateKey   string
	ClientEmail  string
	ClientID     string
	TokenURL     string
}
