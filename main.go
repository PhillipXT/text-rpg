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
    "regexp"
    "slices"
)

const maxLineLength = 80

type iVisible interface {
    makeVisible()
}

type iPrintable interface {
    printData()
}

type roomData struct {
    desc string
    portals []portalData
}

type portalData struct {
    id int
    portalType string
    direction []string
    dest_room_id int
}

// Add some flags here: canTake, canBreak, etc. to save on adding actions
type itemData struct {
    name string                 // Simple name for the item
    room int                    // Current position of the item
    status int                  // Current status of the item
    visible bool                // Whether or not the item can be seen
    statusList []statusItem     // List of descriptions available
    actionList []actionItem     // Actions available to be performed
}

func (i *itemData) makeVisible() {
    i.visible = true
}

func (item itemData) printData() {
    fmt.Printf("Item: [%v] Status:[%v] Room: [%v] Visible: [%v]\n", item.name, item.status, item.room, item.visible)
}

type statusItem struct {
    status int
    desc string
}

type actionItem struct {
    action string
    requiredStatus int
    resultingStatus int
    desc string
    triggerList []triggerItem
}

type triggerItem struct {
    name string
    param string
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
    items := getItemData()

    for {
        room := rooms[currentRoom]
        viewRoom(currentRoom, room, items)
        fmt.Printf("What do we do now?  ")
        cmd := getCommand()
        clearScreen()
        message, new_room, did_move := processCommand(currentRoom, room, items, cmd)
        if did_move {
            currentRoom = new_room
        }
        fmt.Println(getResponse(name))
        //fmt.Printf("We will: %v\n\n", cmd)
        fmt.Printf(message + "\n\n")
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
    responses := map[int]string {
        0: fmt.Sprintf("Okay %v, I'll give that a shot.\n", name),
        1: "Are you serious? Hey, it's your funeral.\n",
        2: "Ugh, do we have to?\n",
    }
    return responses[rand.Intn(len(responses))]
}

func getCommand() string {
    reader := bufio.NewReader(os.Stdin)
    cmd, err := reader.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    return strings.TrimSuffix(cmd, "\n")
}

func processCommand(currentRoom int, room roomData, items []itemData, cmd string) (string, int, bool) {

    cmd = strings.ToLower(cmd)

    // Directions could be programmatically determined from the portals arrays
    // in roomData.  That way, the list could be dynamically generated.
    directions := []string {
        "north", "n",
        "south", "s",
        "east", "e",
        "west", "w",
        "up", "u",
        "down", "d",
    }

    quit := []string { "quit", "q" }

    reg, err := regexp.Compile("[^a-z ]+")
    if err != nil {
        log.Printf("Error in regexp: %w", err)
        return "", 0, false
    }

    cmd = reg.ReplaceAllString(cmd, "")
    tokens := strings.Fields(cmd)

    /*
    for i, token := range tokens {
        fmt.Printf("Token %v: [%v]\n", i, token)
    }
    fmt.Println("Press a key to continue...")
    fmt.Scanln()
    */

    if len(tokens) == 0 {
        return "Stop wasting my time and give me something to do.", 0, false
    }

    if slices.Contains(directions, tokens[0]) {
        for _, portal := range room.portals {
            if slices.Contains(portal.direction, tokens[0]) {
                return fmt.Sprintf("You go %v", portal.direction[0]), portal.dest_room_id, true
            }
        }
        return "You can't go that direction.", 0, false
    } else if tokens[0] == "i" {
        fmt.Println("Your backpack contains:")
        for _, item := range items {
            if item.room == 0 {
                fmt.Println("    " + item.name)
            }
        }
    } else if tokens[0] == "look" {
        if len(tokens) == 1 {
            return "Look at what?", 0, false
        } else {
            item := itemData{}
            for _, i := range items {
                if i.name == tokens[1] {
                    item = i
                    break
                }
            }
            if item.name != "" {
                if !item.visible {
                    return fmt.Sprintf("I don't see a %v anywhere.", item.name), 0, false
                }
                for _, action := range item.actionList {
                    if action.action == tokens[0] && action.requiredStatus == item.status {
                        item.status = action.resultingStatus
                        return action.desc, 0, false
                    }
                }
            }
            return "I don't see anything special about that.", 0, false
        }
    } else if tokens[0] == "smash" {
        // BUG:  Make sure we're in the same room as the object first
        if len(tokens) == 1 {
            return "Smash what?", 0, false
        } else if tokens[1] == "vase" {
            var text string
            for _, action := range items[0].actionList {
                if action.requiredStatus == items[0].status {
                    text = action.desc
                    items[0].status = action.resultingStatus
                    for _, trigger := range action.triggerList {
                        if trigger.name == "makeVisible" {
                            items[1].makeVisible()
                        }
                    }
                    break
                }
            }
            return text, 0, false
        }
    } else if tokens[0] == "take" {
        if len(tokens) == 1 {
            return "Take what?", 0, false
        } else if tokens[1] == "key" {
            if items[1].visible == false || items[1].room != currentRoom {
                return "I don't see a key anywhere.", 0, false
            }
            items[1].room = 0
            var text string
            for _, action := range items[1].actionList {
                if action.requiredStatus == items[1].status {
                    text = action.desc
                }
            }
            items[1].status = 1
            return text, 0, false
        }
    } else if tokens[0] == "drop" {
        if len(tokens) == 1 {
            return "Drop what?", 0, false
        } else {
            item := itemData{}
            for _, i := range items {
                if i.name == tokens[1] && i.room == 0 {
                    item = i
                    break
                }
            }
            //fmt.Printf("Item: [%v], Action: [%v]\n", item.name, tokens[0])
            if item.name != "" {
                for _, action := range item.actionList {
                    if action.action == tokens[0] {
                        items[1].room = currentRoom
                        items[1].status = action.resultingStatus
                        fmt.Println(item)
                        return action.desc, 0, false
                    }
                }
            }
            return "I don't have an item like that.", 0, false
        }
    } else if slices.Contains(quit, tokens[0]) {
        os.Exit(0)
    }
    return "", 0, false
}

func getRoomData() map[int]roomData {
    arr := make(map[int]roomData)
    arr[0] = roomData {
        "This is a really nice backpack, with lots of space to put things.", nil,
    }

    arr[1] = roomData {
        "You are in a simple room.  There is a desk and chair up against the wall, and a door to the north.",
        []portalData {
            portalData { 0, "door", []string{ "north", "n" }, 2 },
        },
    }
    arr[2] = roomData {
        "You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south.",
        []portalData {
            portalData { 0, "door", []string{ "south", "s" }, 1 },
        },
    }
    return arr
}

func getPortals() map[int]portalData {
    arr := make(map[int]portalData)
    return arr
}

func getItemData() []itemData {
    
    arr := []itemData{}

    arr = append(arr, itemData { "vase", 1, 0, true,
        []statusItem {
            statusItem { 0, "There is a vase on the desk." },
            statusItem { 1, "The vase on the desk has been smashed into pieces." },
        },
        []actionItem {
            actionItem { "smash", 0, 1, "You smash the vase into tiny pieces.",
                []triggerItem {
                    triggerItem { "makeVisible", "key" },
                },
            },
            actionItem { "look", 0, 0, "This drab coloured vase looks empty.", nil },
            actionItem { "look", 1, 1, "It's not much of a vase any more. Shards of it are all over the table.", nil },
            actionItem { "examine", 0, 0, "It's too dark to see inside, and too narrow for your hand, but you can hear something jingle inside.", nil },
            actionItem { "examine", 1, 1, "There's really nothing to examine anymore.", nil },
            actionItem { "smash", 1, 1, "The vase has already been smashed.", nil },
        },
    })

    arr = append(arr, itemData { "key", 1, 0, false,
        []statusItem {
            statusItem { 0, "There is a key lying amongst the shards of the vase." },
            statusItem { 1, "The key is in your backpack." },
            statusItem { 2, "There is a key on the ground." },
        },
        []actionItem {
            actionItem { "look", 0, 0, "This key is metallic and sturdy, and looks like it would open a door.", nil },
            actionItem { "take", 0, 1, "You take the key.", nil },
            actionItem { "take", 2, 1, "You take the key.", nil },
            actionItem { "drop", 1, 2, "You drop the key on the ground.", nil },
        },
    })

    return arr
}

func getPortalData() map[int]bool {
    return map[int]bool { 0: false, 1: false, 2: false, }
}

func viewRoom(currentRoom int, room roomData, items []itemData) {
    printLine("", strings.Repeat("=", maxLineLength))
    room_desc := room.desc
    for _, item := range items {
        item.printData()
        if item.room == currentRoom && item.visible {
            for _, status := range item.statusList {
                if status.status == item.status {
                    room_desc += " " + status.desc
                }
            }
        }
    }
    printLine("> ", room_desc)
    printLine(">", "")
    printLine(">", "")
    exits := "Visible exits: "
    for _, portal := range room.portals {
        exits += portal.direction[1] + " "
    }
    printLine("> ", exits)
    printLine("", strings.Repeat("=", maxLineLength))
    printLine("", "")
}

func printLine(prefix string, text string) {

    if len(prefix) + len(text) <= maxLineLength {
        fmt.Println(prefix + text)
        return
    }

    start := 0

    for {
        end := start + maxLineLength - len(prefix)
        if end > len(text) {
            end = len(text)
        }
        substring := text[start:end]
        last_space := strings.LastIndexByte(substring, ' ')
        if last_space == -1 || end == len(text) {
            fmt.Println(prefix + substring)
            start = start + maxLineLength - len(prefix)
        } else {
            fmt.Println(prefix + text[start:start + last_space])
            start = start + last_space + 1
        }
        if end >= len(text) {
            break
        }
    }
}
