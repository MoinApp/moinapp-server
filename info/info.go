package info

const (
	AppName = "MoinApp-Server"
)

var AppVersion string

func CheckCorrectCompilation() {
	if AppVersion == "" {
		AppVersion = "unknown"
	}
}
