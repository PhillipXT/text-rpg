package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const maxLineLength = 80

type iVisible interface {
	makeVisible()
}

type iPrintable interface {
	printData()
}

type Player struct {
	name string
}

type roomData struct {
	desc    string
	portals []portalData
}

type portalData struct {
	id           int
	portalType   string
	direction    []string
	dest_room_id int
}

// Add some flags here: canTake, canBreak, etc. to save on adding actions
type itemData struct {
	name       string       // Simple name for the item
	room       int          // Current position of the item
	status     int          // Current status of the item
	visible    bool         // Whether or not the item can be seen
	statusList []statusItem // List of descriptions available
	actionList []actionItem // Actions available to be performed
}

func (i *itemData) makeVisible() {
	i.visible = true
}

func (item itemData) printData() {
	fmt.Printf("Item: [%v] Status:[%v] Room: [%v] Visible: [%v]\n", item.name, item.status, item.room, item.visible)
}

type statusItem struct {
	status int
	desc   string
}

type actionItem struct {
	action          string
	requiredStatus  int
	resultingStatus int
	desc            string
	triggerList     []triggerItem
}

type triggerItem struct {
	name  string
	param string
}

func main() {

	player := Player{}

	clearScreen()
	getPlayerName(&player)

	clearScreen()
	displayWelcome(player.name)

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
		fmt.Println(getResponse(player.name))
		//fmt.Printf("We will: %v\n\n", cmd)
		fmt.Printf(message + "\n\n")
	}
}

func displayWelcome(name string) {
	printLine("", fmt.Sprintf("Welcome back, %v. You've been 'asleep' for some time.\n", name))
	printLine("", fmt.Sprintf("I can only imagine the headache you have right now, given the size of that lump on your forehead.\n"))
	printLine("", fmt.Sprintf("Take your time. Look around. Figure out where you are.\n"))
	printLine("", fmt.Sprintf("Try to remember... WHO you are....\n"))
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
	responses := map[int]string{
		0: fmt.Sprintf("Okay %v, I'll give that a shot.\n", name),
		1: "Are you serious? Hey, it's your funeral.\n",
		2: "Ugh, do we have to?\n",
	}
	return responses[rand.Intn(len(responses))]
}

func getRoomData() map[int]roomData {
	arr := make(map[int]roomData)
	arr[0] = roomData{
		"This is a really nice backpack, with lots of space to put things.", nil,
	}

	arr[1] = roomData{
		"You are in a simple room.  There is a desk and chair up against the wall, and a door to the north.",
		[]portalData{
			portalData{0, "door", []string{"north", "n"}, 2},
		},
	}
	arr[2] = roomData{
		"You are in another room.  There is a bookshelf that is mostly empty. There is a door to the south.",
		[]portalData{
			portalData{0, "door", []string{"south", "s"}, 1},
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

	arr = append(arr, itemData{"vase", 1, 0, true,
		[]statusItem{
			statusItem{0, "There is a vase on the desk."},
			statusItem{1, "The vase on the desk has been smashed into pieces."},
		},
		[]actionItem{
			actionItem{"smash", 0, 1, "You smash the vase into tiny pieces.",
				[]triggerItem{
					triggerItem{"makeVisible", "key"},
				},
			},
			actionItem{"look", 0, 0, "This drab coloured vase looks empty.", nil},
			actionItem{"look", 1, 1, "It's not much of a vase any more. Shards of it are all over the table.", nil},
			actionItem{"examine", 0, 0, "It's too dark to see inside, and too narrow for your hand, but you can hear something jingle inside.", nil},
			actionItem{"examine", 1, 1, "There's really nothing to examine anymore.", nil},
			actionItem{"smash", 1, 1, "The vase has already been smashed.", nil},
		},
	})

	arr = append(arr, itemData{"key", 1, 0, false,
		[]statusItem{
			statusItem{0, "There is a key lying amongst the shards of the vase."},
			statusItem{1, "The key is in your backpack."},
			statusItem{2, "There is a key on the ground."},
		},
		[]actionItem{
			actionItem{"look", 0, 0, "This key is metallic and sturdy, and looks like it would open a door.", nil},
			actionItem{"take", 0, 1, "You take the key.", nil},
			actionItem{"take", 2, 1, "You take the key.", nil},
			actionItem{"drop", 1, 2, "You drop the key on the ground.", nil},
		},
	})

	return arr
}

func getPortalData() map[int]bool {
	return map[int]bool{0: false, 1: false, 2: false}
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

	if len(prefix)+len(text) <= maxLineLength {
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
			fmt.Println(prefix + text[start:start+last_space])
			start = start + last_space + 1
		}
		if end >= len(text) {
			break
		}
	}
}
