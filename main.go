package main

import (
    "os"
    "os/exec"
    "runtime"
    "fmt"
)

func main() {

    opsys := runtime.GOOS

    var name string

    clearScreen(opsys)
    fmt.Println("What is your name, brave adventurer?")
    fmt.Scanln(&name)

    clearScreen(opsys)
    fmt.Printf("Hello %v. Welcome to Eternia.\n", name)

    for {
        cmd := ""
        fmt.Println("What do we do now?")
        fmt.Scanln(&cmd)
        clearScreen(opsys)
        fmt.Printf("Okay %v, I'll give that a shot.\n", name)
        fmt.Println()
    }
}

func clearScreen(opsys string) {
    var cmd *exec.Cmd
    switch opsys {
    case "windows":
        cmd = exec.Command("cmd", "/c", "cls")
    default:
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    cmd.Run()
}
