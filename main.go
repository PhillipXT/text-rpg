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

type roomData struct {
    desc string
    portals []portalData
    items []itemData
}

type portalData struct {
    id int
    portalType string
    direction string
    dest_room_id int
}

type itemData struct {
    id int
    name string
    currentStatus int
    desc []descriptionData
    action []actionData
}

type descriptionData struct {
    status int
    desc string
}

type actionData struct {
    action string
    resultingStatus int
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
        room := rooms[currentRoom]
        viewRoom(room)
        fmt.Printf("What do we do now?  ")
        cmd := getCommand()
        clearScreen()
        r, move := processCommand(room, cmd)
        if move {
            currentRoom = r
        }
        fmt.Println(getResponse(name))
        fmt.Printf("We will: %v\n\n", cmd)
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

func processCommand(room roomData, cmd string) (int, bool) {
    switch cmd {
    case "n":
        return 1, true
    case "s":
        return 0, true
    case "smash vase":
        room.items[0].currentStatus = 1
    case "q", "quit":
        os.Exit(0)
    }
    return 0, false
}


func getRoomData() (map[int]roomData) {
    r := make(map[int]roomData)
    r[0] = roomData {
        "You are in a simple room.  There is a desk and chair up against the wall, and a door to the north.",
        []portalData { portalData { 0, "door", "N", 1 },},
        []itemData { itemData { 0, "Vase", 0,
            []descriptionData {
                descriptionData { 0, "There is a vase on the desk." },
                descriptionData { 1, "The vase on the desk has been smashed into pieces." }, },
            []actionData { actionData { "smash", 1 },},
        }, },
    }
    r[1] = roomData {
        "You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south.",
        []portalData { portalData { 0, "door", "S", 0 },},
        []itemData{},
    }
    return r
}

func getPortalData() map[int]bool {
    return map[int]bool { 0: false, 1: false, 2: false, }
}

func viewRoom(room roomData) {
    fmt.Printf(room.desc)
    for _, item := range room.items {
        for _, desc := range item.desc {
            if desc.status == item.currentStatus {
                fmt.Printf(" " + desc.desc)
            }
        }
    }
    fmt.Println()
    fmt.Println()
    fmt.Printf("Visible exits: ")
    for _, portal := range room.portals {
        fmt.Printf(portal.direction + " ")
    }
    fmt.Printf("\n\n")
}
