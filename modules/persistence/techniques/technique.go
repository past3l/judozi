package techniques

// Technique represents a persistence mechanism
type Technique interface {
	// Name returns the technique name
	Name() string
	
	// Description returns a short description
	Description() string
	
	// StealthLevel returns stealth rating: "low", "medium", "high"
	StealthLevel() string
	
	// SurvivesReboot returns true if persistence survives system reboot
	SurvivesReboot() bool
	
	// Install deploys the persistence mechanism
	Install() error
	
	// Remove cleans up the persistence mechanism
	Remove() error
	
	// Details returns additional information after installation
	Details() string
}

// GetAll returns all available persistence techniques
func GetAll() []Technique {
	return []Technique{
		NewSSHKey(),
		NewCronJob(),
		NewSystemdService(),
		NewSUIDBinary(),
		NewLDPreload(),
		NewPAMBackdoor(),
		NewShellRC(),
		NewSetuidShell(),
	}
}
