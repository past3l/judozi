package techniques

// Technique represents a single recon technique
type Technique interface {
	Name() string
	Description() string
	StealthLevel() string // "passive", "low", "medium"
	Run() string
}

// GetAll returns all available recon techniques
func GetAll() []Technique {
	return []Technique{
		NewSystemInfo(),
		NewNetworkInfo(),
		NewUserEnum(),
		NewSUIDSGID(),
		NewProcessEnum(),
		NewCredHunter(),
		NewCronEnum(),
		NewServiceEnum(),
		NewDockerEnum(),
		NewEnvSecrets(),
		NewSSHEnum(),
		NewFSWritable(),
	}
}
