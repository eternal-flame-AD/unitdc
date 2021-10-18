package localize

import "fmt"
import "os"
import "os/exec"
import "strings"

func getLang() string {
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		langLocRaw := strings.TrimSpace(string(output))
		langLoc := strings.Split(langLocRaw, "-")
		return langLoc[0]
	}
	fmt.Fprintf(os.Stderr, "could not get locale: %v", err)
	return ""
}
