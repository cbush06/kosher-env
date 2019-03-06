package main

import (
	"log"
	"os"
	"runtime"

	registry "golang.org/x/sys/windows/registry"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("argument required")
	}

	action := os.Args[1]
	if action != "install" && action != "uninstall" {
		log.Fatal("argument must be 'install' or 'uninstall'")
	}

	system := runtime.GOOS
	switch system {
	case "windows":
		doWindows(action)
	case "linux":
		doLinux(action)
	default:
		log.Fatalf("unsupported OS [%s] detected", runtime.GOOS)
	}
}

func doWindows(action string) {
	switch action {
	case "install":
		doWindowsInstall()
	case "uninstall":
	}
}

func doWindowsInstall() {
	// set path
	if registryEnvEntry, err := registry.OpenKey(registry.LOCAL_MACHINE, `HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment\Path`, registry.QUERY_VALUE|registry.READ|registry.WRITE); err != nil {
		log.Fatalf("error encountered while opening PATH registry key: %s", err)
	}
}

func doWindowsUninstall() {

}

func doLinux(action string) {
	switch action {
	case "install":
	case "uninstall":
	}
}

func doLinuxInstall() {

}

func doLinuxUninstall() {

}
