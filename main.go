package main

import (
    "math/rand"
    "os"
    "os/exec"
    "runtime"
    "fmt"
    "time"
)

type portalData struct {
    id int
    portalType string
    direction string
    dest_room_id int
}

type roomData struct {
    desc string
    portals []portalData
}

func main() {

    rand.Seed(time.Now().UnixNano())

    var name string

    clearScreen()
    fmt.Println("What is your name, brave adventurer?")
    fmt.Scanln(&name)

    clearScreen()
    fmt.Printf("Hello %v. Welcome to Eternia.\n\n", name)

    currentRoom := 1

    //portals := getPortalData()
    rooms := getRoomData()
    
    for {
        viewRoom(rooms[currentRoom])
        cmd := getCommand()
        fmt.Printf("%v\n", getResponse(name))
        fmt.Printf("We will: %v\n\n", cmd)
        if currentRoom == 1 {
            currentRoom = 2
        } else {
            currentRoom = 1
        }
    }
}

func clearScreen() {
    var cmd *exec.Cmd
    switch runtime.GOOS {
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

func getCommand() string {
    cmd := ""
    fmt.Printf("What do we do now?  ")
    fmt.Scanln(&cmd)
    clearScreen()
    return cmd
}

func getRoomData() (map[int]roomData) {
    r := make(map[int]roomData)
    r[1] = roomData {
        "You are in a simple room.  There is a desk and chair up against the wall, and a door to the north.",
        []portalData { portalData { 1, "door", "N", 2 },},
    }
    r[2] = roomData {
        "You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south.",
        []portalData { portalData { 1, "door", "S", 1 },},
    }
    return r
}

func getPortalData() map[int]bool {
    return map[int]bool { 0: false, 1: false, 2: false, }
}

func viewRoom(room roomData) {
    fmt.Println(room.desc)
    fmt.Printf("\nVisible exits: ")
    for _, portal := range room.portals {
        fmt.Printf(portal.direction)
    }
    fmt.Println("\n")
}
