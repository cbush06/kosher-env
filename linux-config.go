package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("2 arguments required")
	}

	action := os.Args[1]
	if action != "install" && action != "uninstall" {
		log.Fatal("argument must be 'install' or 'uninstall'")
	}

	installationDir := os.Args[2]

	doLinux(action, installationDir)
}

func doLinux(action string, installationDir string) {
	switch action {
	case "install":
		doLinuxInstall(installationDir)
	case "uninstall":
		doLinuxUninstall(installationDir)
	}
}

func doLinuxInstall(installationDir string) {
	log.Printf("Linux Install: %s\n", installationDir)

	f, _ := os.OpenFile("/etc/profile.d/kosher.sh", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	defer f.Close()

	f.WriteString(fmt.Sprintf(`export PATH="%s:$PATH"`, installationDir))
}

func doLinuxUninstall(installationDir string) {
	os.Remove("/etc/profile.d/kosher.sh")
}
