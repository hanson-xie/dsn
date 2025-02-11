package build

// BuildVersion is the local build version
const BuildVersion = "1.0.0"
const BuildUnit = "-bedrock"

var CurrentCommit string

func UserVersion() string {
	return BuildVersion + BuildUnit + CurrentCommit
}
