package main

import (
    "math/rand"
    "os"
    "os/exec"
    "runtime"
    "fmt"
    "time"
    "bufio"
    "strings"
    "log"
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

    currentRoom := 0

    //portals := getPortalData()
    rooms := getRoomData()
    
    for {
        viewRoom(rooms[currentRoom])
        fmt.Printf("What do we do now?  ")
        cmd := getCommand()
        clearScreen()
        processCommand(cmd)
        fmt.Println(getResponse(name))
        fmt.Printf("We will: %v\n\n", cmd)
        if currentRoom == 0 {
            currentRoom = 1
        } else {
            currentRoom = 0
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
    reader := bufio.NewReader(os.Stdin)
    cmd, err := reader.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    return strings.TrimSuffix(cmd, "\n")
}

func processCommand(cmd string) {
    switch cmd {
    case "q", "quit":
        os.Exit(0)
    }
}


func getRoomData() (map[int]roomData) {
    r := make(map[int]roomData)
    r[0] = roomData {
        "You are in a simple room.  There is a desk and chair up against the wall, and a door to the north.",
        []portalData { portalData { 1, "door", "N", 1 },},
    }
    r[1] = roomData {
        "You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south.",
        []portalData { portalData { 1, "door", "S", 0 },},
    }
    return r
}

func getPortalData() map[int]bool {
    return map[int]bool { 0: false, 1: false, 2: false, }
}

func viewRoom(room roomData) {
    fmt.Println(room.desc)
    fmt.Println()
    fmt.Printf("Visible exits: ")
    for _, portal := range room.portals {
        fmt.Printf(portal.direction + " ")
    }
    fmt.Printf("\n\n")
}
