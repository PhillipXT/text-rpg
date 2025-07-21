package main

import (
    "errors"
    "math/rand"
    "os"
    "os/exec"
    "runtime"
    "fmt"
    "time"
)

func main() {

    opsys := runtime.GOOS
    rand.Seed(time.Now().UnixNano())

    var name string

    clearScreen(opsys)
    fmt.Println("What is your name, brave adventurer?")
    fmt.Scanln(&name)

    clearScreen(opsys)
    fmt.Printf("Hello %v. Welcome to Eternia.\n", name)

    position := 0

    for {
        cmd := ""
        roomDesc, err := getRoomData(int(position))
        if err != nil {
            fmt.Println("We're not on the map anymore.")
        }
        if position == 0 {
            position = 1
        } else {
            position = 0
        }
        fmt.Println(roomDesc)
        fmt.Println()
        fmt.Println("What do we do now?")
        fmt.Scanln(&cmd)
        clearScreen(opsys)
        fmt.Printf(getResponse(name))
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

func getResponse(name string) string {
    var response string
    switch rand.Intn(3) {
    case 0:
        response = fmt.Sprintf("Okay %v, I'll give that a shot.\n", name)
    case 1:
        response = "Are you serious? Hey, it's your funeral.\n"
    case 2:
        response = "Ugh, do we have to?\n"
    }
    return response
}

func getRoomData(room int) (string, error) {
    var desc string
    switch room {
    case 0:
        desc = "You are in a simple room.  There is a desk and chair up against the wall, and a door to the north."
    case 1:
        desc = "You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south."
    default:
        return "", errors.New("Room id not found")
    }
    return desc, nil
}
