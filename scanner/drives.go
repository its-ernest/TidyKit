package scanner

import (
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

type DriveInfo struct {
	Device     string
	MountPoint string
	TotalBytes uint64
	FreeBytes  uint64
	UsedBytes  uint64
	FSType     string
}

func FetchDrives() ([]DriveInfo, error) {
	out, err := exec.Command("df", "-P").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var drives []DriveInfo

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		device := fields[0]
		mountPoint := fields[5]

		if !strings.HasPrefix(device, "/dev/") {
			continue
		}

		var stat unix.Statfs_t
		if err := unix.Statfs(mountPoint, &stat); err != nil {
			continue
		}

		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bavail * uint64(stat.Bsize)
		used := total - free

		drives = append(drives, DriveInfo{
			Device:     device,
			MountPoint: mountPoint,
			TotalBytes: total,
			FreeBytes:  free,
			UsedBytes:  used,
		})
	}

	return drives, nil
}
