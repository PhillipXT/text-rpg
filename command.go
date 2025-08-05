package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "slices"
    "strings"
)

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
