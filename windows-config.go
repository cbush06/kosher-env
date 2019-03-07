package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/sys/windows/registry"
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

	doWindows(action, installationDir)
}

func doWindows(action string, installationDir string) {
	switch action {
	case "install":
		doWindowsInstall(installationDir)
	case "uninstall":
		doWindowsUninstall(installationDir)
	}
}

func doWindowsInstall(installationDir string) {
	var (
		environment registry.Key
		err         error
	)

	// set path
	if environment, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.READ|registry.WRITE); err != nil {
		log.Fatalf("error encountered while opening PATH registry key: %s", err)
	}
	defer environment.Close()

	path, _, _ := environment.GetStringValue("Path")
	path = strings.Join([]string{path, installationDir}, ":")
	environment.SetStringValue("Path", path)
}

func doWindowsUninstall(installationDir string) {
	var (
		environment registry.Key
		err         error
	)

	// set path
	if environment, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.READ|registry.WRITE); err != nil {
		log.Fatalf("error encountered while opening PATH registry key: %s", err)
	}
	defer environment.Close()

	path, _, _ := environment.GetStringValue("Path")

	r := regexp.MustCompile(fmt.Sprintf(`\:{0,1}%s`, strings.Replace(installationDir, `\`, `\\`, -1)))
	path = r.ReplaceAllString(path, "")

	environment.SetStringValue("Path", path)
}
