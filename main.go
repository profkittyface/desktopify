package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"bytes"

	"github.com/fatih/color"
)

var installPath string = "/opt/AppImage"
var iconPath string = "/usr/share/icons/AppImage"

func main() {
	// Show help if program run without paramters
	if len(os.Args) == 1 {
		ShowHelp()
		os.Exit(0)
	}

	switch args := os.Args[1]; args {
	case "install":
		if len(os.Args) == 2 {
			fmt.Println("Please specify an AppImage to install")
			os.Exit(1)
		}

		if CheckForRoot() == false {
			fmt.Println("Need root for install operation. Please rerun desktopify with sudo")
			os.Exit(1)
		}
		
		CheckInstallDirectory()
		appImageSrc := TransformFile(os.Args[2])
		appImageDest := GetDestinationPath(appImageSrc)
		CopyFile(appImageSrc, appImageDest)
		iconLocation := ExtractIcon(appImageSrc)
		GenerateShortcut(os.Args[2], iconLocation)
		fmt.Println("Installed AppImage to", appImageDest)
	case "list":
		ListAppImages()
	case "remove":
		if CheckForRoot() == false {
			fmt.Println("Need root for removal operation. Please rerun desktopify with sudo")
			os.Exit(1)
		}
		fmt.Printf("This will remove %s from the install directory. Confirm y/n \n", os.Args[2])
		var s string
		_, err := fmt.Scan(&s)
		CheckError(err)
		if s == "y" || s == "yes" {
			RemoveAppImage(os.Args[2])
			fmt.Printf("%s removed from system\n", os.Args[2])
		} else {
			fmt.Println("Aborted")
			os.Exit(1)
		}
	case "debug":
		ExtractIcon(os.Args[2])
	default:
		ShowHelp()
	}
}

func CheckInstallDirectory() {
	_, err := os.Stat(installPath)
	if os.IsNotExist(err) {
		fmt.Println("Installation directory /opt/AppImage created")
		_ = os.Mkdir(installPath, 0775)
	}
	_, err = os.Stat(iconPath)
	if os.IsNotExist(err) {
		fmt.Println("Icon directory /usr/share/icons created")
		_ = os.Mkdir(iconPath, 0775)
	}
}

func CheckForRoot() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	CheckError(err)
	if i == 0 {
		return true
	} else {
		return false
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TransformFile(file string) string {
	splitPath := strings.Split(file, "/")
	var returnText string
	if len(splitPath) == 1 {
		cwd, _ := os.Getwd()
		fullDir := cwd + "/" + file
		returnText = fullDir
	} else {
		returnText = file
	}
	return returnText
}

func GetDestinationPath(file string) string {
	splitPath := strings.Split(file, "/")
	result := installPath + "/" + splitPath[len(splitPath)-1]
	return result
}

func CopyFile(src string, dst string) {
	source, err := os.Open(src)
	destination, err := os.Create(dst)
	_, err = io.Copy(destination, source)
	CheckError(err)
	source.Close()
	destination.Close()
	os.Chmod(dst, 0775)
}

func ListAppImages() {
	files, err := ioutil.ReadDir(installPath)
	CheckError(err)
	for _, file := range files {
		if strings.Contains(file.Name(), "AppImage"){
				fmt.Println(file.Name())
		}
	}
}

func RemoveAppImage(name string) {
	// Remove AppImage
	file := installPath + "/" + name
	err := os.Remove(file)
	CheckError(err)
	// Remove desktop shortcut
	s := strings.Split(name, "/")
	AppImageName := s[len(s)-1]
	programName := strings.Split(AppImageName, "-")
	desktopFile := "/home/" + os.Getenv("SUDO_USER") + "/.local/share/applications/" + programName[0] + ".desktop"
	err = os.Remove(desktopFile)
	CheckError(err)

	files, err := ioutil.ReadDir(iconPath)
	CheckError(err)

	var iconFile string

	for _, file := range files {
		if strings.Contains(file.Name(), programName[0]) {
			iconFile = file.Name()
		}
	}

	err = os.Remove(iconPath + "/" + iconFile)
	CheckError(err)
}

type Shortcut struct {
	Name string
	Icon string
	ExecFile string
	DesktopFile string
}

func GenerateShortcut(file string, iconLocation string) {
	fileName := TransformFile(file)
	st := Shortcut{}
	st.DesktopFile =`[Desktop Entry]
Name={{.Name}}
Comment={{.Name}}
Icon={{.Icon}}
Exec={{.ExecFile}}
Terminal=false
Type=Application
`
	s := strings.Split(fileName, "/")
	name := s[len(s)-1]
	programName := strings.Split(name, "-")
	st.Name = programName[0]
	st.Icon = iconLocation
	st.ExecFile = GetDestinationPath(fileName)

	var b bytes.Buffer
	t, err := template.New("").Parse(st.DesktopFile)
	CheckError(err)
	err = t.Execute(&b, st)
	desktopFile := "/home/" + os.Getenv("SUDO_USER") + "/.local/share/applications/" + programName[0] + ".desktop"
	by := []byte(b.String())
	f, err := os.Create(desktopFile)
	f.Write(by)
	f.Close()
	os.Chmod(desktopFile, 0775)
	CheckError(err)
}

func ExtractIcon(appImageFile string) string {
	var srcIconFile string
	var destIconFile string

	os.Chmod(appImageFile, 0775)
	cmd := exec.Command(appImageFile, "--appimage-extract")
	cmd.Dir = "/tmp"
	cmd.Run()
	files, err := ioutil.ReadDir("/tmp/squashfs-root")
	CheckError(err)

	var iconFile string
	for _, file := range files {
		if strings.Contains(file.Name(), ".png") {
			iconFile = file.Name()
		}
		if strings.Contains(file.Name(), ".svg") {
			iconFile = file.Name()
		}
	}

	iconFileSplit := strings.Split(iconFile, ".")

	s := strings.Split(appImageFile, "/")
	name := s[len(s)-1]
	programName := strings.Split(name, "-")

	srcIconFile = "/tmp/squashfs-root/" + iconFile
	destIconFile = iconPath + "/" + programName[0] + "." + iconFileSplit[1]

	err = os.Rename(srcIconFile, destIconFile)
	CheckError(err)
	return destIconFile
}

func ShowHelp() {
	color.Cyan("Desktopify - Your AppImage Installer")
	fmt.Printf("\n")
	fmt.Println("List images")
	color.Green("desktopify list")
	fmt.Printf("\n")
	fmt.Println("Install image")
	color.Green("desktopify install Software.AppImage")
	fmt.Printf("\n")
	fmt.Println("Remove image")
	color.Green("desktopify remove Software.AppImage")
}
