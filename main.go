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
	// "bufio"
	"text/template"
	"bytes"
)

var installPath string = "/opt/AppImage"

func main() {
	// reader := bufio.NewReader(os.Stdin)

	switch args := os.Args[1]; args {
	case "install":
		if CheckForRoot() == false {
			fmt.Println("Need root for install operation. Please rerun desktopify with sudo")
		}
		CheckInstallDirectory()
		appImageSrc := TransformFile(os.Args[2])
		appImageDest := GetDestinationPath(appImageSrc)
		CopyFile(appImageSrc, appImageDest)
		fmt.Println("Installed AppImage to", appImageDest)
	case "list":
		ListAppImages()
	case "remove":
		fmt.Printf("This will remove %s from the install directory. Confirm y/n \n", os.Args[2])

		var s string
		_, err := fmt.Scan(&s)
		CheckError(err)
		// response, _ := reader.ReadString('\n')
		// RemoveAppImage(os.Args[2])
		if s == "y" {
			RemoveAppImage(os.Args[2])
		} else {
			os.Exit(1)
		}
	case "debug":
		GenerateShortcut("Program")
	default:
		fmt.Println("Unrecognized input")
	}
}

func CheckInstallDirectory() {
	_, err := os.Stat(installPath)
	if os.IsNotExist(err) {
		fmt.Println("Installation directory /opt/AppImage created")
		_ = os.Mkdir(installPath, 0665)
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
	defer source.Close()
	destination, err := os.Create(dst)
	defer destination.Close()
	_, err = io.Copy(destination, source)
	CheckError(err)
}

func ListAppImages() {
	files, err := ioutil.ReadDir(installPath)
	CheckError(err)
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func RemoveAppImage(name string) {
	file := installPath + "/" + name
	fmt.Println(file)
	err := os.Remove(file)
	CheckError(err)
}

type Shortcut struct {
	Name string
	Icon string
	ExecFile string
	DesktopFile string
}

func GenerateShortcut(file string) {
	st := Shortcut{}
	st.DesktopFile =`
[Desktop Entry]
Name={{.Name}}
Comment={{.Name}}
Icon={{.Icon}}
Exec={{.ExecFile}}
Terminal=false
Type=Application
`
	st.Name = "Demo"
	st.Icon = "Demo"
	st.ExecFile = "Demo"

	fmt.Println(st.DesktopFile)
	var b bytes.Buffer
	t, err := template.New("").Parse(st.DesktopFile)
	CheckError(err)
	err = t.Execute(&b, st)
	fmt.Println(b.String())
}
