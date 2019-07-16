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
		environment   registry.Key
		err           error
		bfcache32     registry.Key
		bfcache32Path = `SOFTWARE\Microsoft\Internet Explorer\Main\FeatureControl\FEATURE_BFCACHE`
		bfcache64     registry.Key
		bfcache64Path = `SOFTWARE\Wow6432Node\Microsoft\Internet Explorer\Main\FeatureControl\FEATURE_BFCACHE`
		bfcacheErr    = `error encountered while creating/opening [HKEY_LOCAL_MACHINE\%s]: %s`
		iexploreErr   = `error encountered while setting subkey [iexplore.exe] within [%s]: %s`
	)

	// set path
	if environment, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.READ|registry.WRITE); err != nil {
		log.Fatalf("error encountered while opening PATH registry key: %s", err)
	}
	defer environment.Close()

	path, _, _ := environment.GetStringValue("Path")
	path = strings.Join([]string{path, installationDir}, ";")
	environment.SetStringValue("Path", path)

	// set BFCACHE (for IE Web Driver) on 32-bit machines
	if bfcache32, _, err = registry.CreateKey(registry.LOCAL_MACHINE, bfcache32Path, registry.WRITE); err != nil {
		log.Fatalf(bfcacheErr, bfcache32Path, err)
	}
	defer bfcache32.Close()
	if err = bfcache32.SetDWordValue("iexplore.exe", 0); err != nil {
		log.Fatalf(iexploreErr, bfcache32Path, err)
	}

	// set BFCACHE (for IE Web Driver) on 64-bit machines
	if bfcache64, _, err = registry.CreateKey(registry.LOCAL_MACHINE, bfcache64Path, registry.WRITE); err != nil {
		log.Fatal(bfcacheErr, bfcache64Path, err)
	} else {
		defer bfcache64.Close()
		if err = bfcache64.SetDWordValue("iexplore.exe", 0); err != nil {
			log.Fatalf(iexploreErr, bfcache64Path, err)
		}
	}
}

func doWindowsUninstall(installationDir string) {
	var (
		environment   registry.Key
		err           error
		bfcache32     registry.Key
		bfcache32Path = `SOFTWARE\Microsoft\Internet Explorer\Main\FeatureControl\FEATURE_BFCACHE`
		bfcache64     registry.Key
		bfcache64Path = `SOFTWARE\Wow6432Node\Microsoft\Internet Explorer\Main\FeatureControl\FEATURE_BFCACHE`
		bfcacheErr    = `error encountered while opening [HKEY_LOCAL_MACHINE\%s]: %s`
		iexploreErr   = `error encountered while removing subkey [iexplore.exe] within [%s]: %s`
	)

	// set path
	if environment, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.READ|registry.WRITE); err != nil {
		log.Fatalf("error encountered while opening PATH registry key: %s", err)
	}
	defer environment.Close()

	path, _, _ := environment.GetStringValue("Path")

	r := regexp.MustCompile(fmt.Sprintf(`;{0,1}%s`, strings.Replace(installationDir, `\`, `\\`, -1)))
	path = r.ReplaceAllString(path, "")

	environment.SetStringValue("Path", path)

	// remove [iexplore.exe] from BFCACHE on 32-bit machines
	if bfcache32, _, err = registry.CreateKey(registry.LOCAL_MACHINE, bfcache32Path, registry.WRITE); err != nil {
		log.Fatalf(bfcacheErr, bfcache32Path, err)
	}
	defer bfcache32.Close()
	if err = bfcache32.DeleteValue("iexplore.exe"); err != nil {
		log.Fatalf(iexploreErr, bfcache32Path, err)
	}

	// remove [iexplore.exe] from BFCACHE on 64-bit machines
	if bfcache64, _, err = registry.CreateKey(registry.LOCAL_MACHINE, bfcache64Path, registry.WRITE); err != nil {
		log.Fatalf(bfcacheErr, bfcache64Path, err)
	} else {
		defer bfcache64.Close()
		if err = bfcache64.DeleteValue("iexplore.exe"); err != nil {
			log.Fatalf(iexploreErr, bfcache64Path, err)
		}
	}
}
