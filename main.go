package main

import (
  "fmt"
  "os"
  "strconv"
  "os/exec"
  "log"
  "strings"
)

func main() {
  fmt.Println(os.Args)
  switch args := os.Args[1]; args {
  case "install":
    fmt.Println("Install was selected")
    fmt.Println(TransformFile(os.Args[2]))
    if CheckForRoot() == false {
      fmt.Println("Need root for directory creation. Please rerun desktopify with sudo")
    }
    CheckInstallDirectory()
  case "list":
    fmt.Println("List was selected")
  case "remove":
    fmt.Println("Remove was selected")
  default:
    fmt.Println("Unrecognized input")
}
}

func CheckInstallDirectory(){
  installDirectory := "/opt/AppImage"
  err := os.Mkdir(installDirectory, 666)
  fmt.Println(err)
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

func CheckError(err error){
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
