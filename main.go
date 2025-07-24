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
}

type portalData struct {
    id int
    portalType string
    direction string
    dest_room_id int
}

type itemData struct {
    name string
    room int
    status int
    desc []descData
    action []actionData
}

type descData struct {
    status int
    desc string
}

type actionData struct {
    action string
    requiredStatus int
    resultingStatus int
    desc string
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
    items := getItemData()
    
    for {
        room := rooms[currentRoom]
        viewRoom(currentRoom, room, items)
        fmt.Printf("What do we do now?  ")
        cmd := getCommand()
        clearScreen()
        r, move := processCommand(room, items, cmd)
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

func processCommand(room roomData, items []itemData, cmd string) (int, bool) {
    switch cmd {
    case "n":
        return 1, true
    case "s":
        return 0, true
    case "smash vase":
        items[0].status = 1
    case "q", "quit":
        os.Exit(0)
    }
    return 0, false
}


func getRoomData() map[int]roomData {
    arr := make(map[int]roomData)
    arr[0] = roomData {
        "You are in a simple room.  There is a desk and chair up against the wall, and a door to the north.",
        []portalData {
            portalData { 0, "door", "N", 1 },
        },
    }
    arr[1] = roomData {
        "You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south.",
        []portalData {
            portalData { 0, "door", "S", 0 },
        },
    }
    return arr
}

func getPortals() map[int]portalData {
    arr := make(map[int]portalData)
    return arr
}

func getItemData() []itemData {
    
    arr := [] itemData {}

    arr = append(arr, itemData { "Vase", 0, 0,
        []descData {
            descData { 0, "There is a vase on the desk." },
            descData { 1, "The vase on the desk has been smashed into pieces." },
        },
        []actionData {
            actionData { "smash", 0, 1, "You smash the vase into tiny pieces." },
            actionData { "smash", 1, 1, "The vase has already been smashed." },
        },
    })

    return arr
}

func getPortalData() map[int]bool {
    return map[int]bool { 0: false, 1: false, 2: false, }
}

func viewRoom(currentRoom int, room roomData, items []itemData) {
    fmt.Printf(room.desc)
    for _, item := range items {
        if item.room == currentRoom {
            for _, desc := range item.desc {
                if desc.status == item.status {
                    fmt.Printf(" " + desc.desc)
                }
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
