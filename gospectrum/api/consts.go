package api

type SpectrumCommand string

const (
	SpectrumAPIPrefix                                  = "/rest"
	SpectrumCommandAuth                SpectrumCommand = "auth"
	SpectrumCommandLsEventLog          SpectrumCommand = "lseventlog"
	SpectrumCommandLsFcMap             SpectrumCommand = "lsfcmap"
	SpectrumCommandLsSystem            SpectrumCommand = "lssystem"
	SpectrumCommandLsSystemStats       SpectrumCommand = "lssystemstats"
	SpectrumCommandLsEnclosure         SpectrumCommand = "lsenclosure"
	SpectrumCommandLsEnclosureCanister SpectrumCommand = "lsenclosurecanister"
	SpectrumCommandLsNodeCanister      SpectrumCommand = "lsnodecanister"
	SpectrumCommandLsMDisk             SpectrumCommand = "lsmdisk"
	SpectrumCommandLsVDisk             SpectrumCommand = "lsvdisk"
	SpectrumCommandLsArray             SpectrumCommand = "lsarray"
	SpectrumCommandLsDrive             SpectrumCommand = "lsdrive"
	SpectrumCommandLsHost              SpectrumCommand = "lshost"
	SpectrumCommandLsHostVDiskMap      SpectrumCommand = "lshostvdiskmap"
)

func (_command SpectrumCommand) String(id string) string {
	if id == "" {
		return SpectrumAPIPrefix + "/" + string(_command)
	} else {
		return SpectrumAPIPrefix + "/" + string(_command) + "/" + id
	}
}
