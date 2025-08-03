package main

import (
    "bufio"
    "fmt"
    "os"
)

func getPlayerName(player *Player) {

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Printf("What is your name, brave adventurer?  (5 to 12 characters)  ")
        scanner.Scan()
        name := scanner.Text()
        if ok, errors := validatePlayerName(name); ok {
            player.name = name
            return
        } else {
            for _, err := range errors {
                fmt.Println(err)
            }
        }
    }
}

func validatePlayerName(name string) (bool, []string) {
    
    var errors []string

    ok := true
    
    if len(name) < 5 {
        errors = append(errors, "Your name must be at least 5 characters.")
        ok = false
    } else if len(name) > 12 {
        errors = append(errors, "Your name must be less than 12 characters.")
        ok = false
    }
    
    return ok, errors
}
