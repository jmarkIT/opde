package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Vault struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Content_Version int    `json:"content_id"`
}

type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	State       string   `json:"state"`
	Created_At  string   `json:"created_at"`
	Permissions []string `json:"permissions"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
	State string `json:"state"`
	Role  string `json:"role"`
}

func getVaultGroups(vault string, account string) []Group {
	var cmd exec.Cmd
	if account != "" {
		cmd = *exec.Command("op", "--format", "json", "--account", account, "vault", "group", "list", vault)
	} else {
		cmd = *exec.Command("op", "--format", "json", "vault", "group", "list", vault)
	}
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(out)
		log.Fatal(err)
	}
	var groups []Group
	err = json.Unmarshal([]byte(out), &groups)
	if err != nil {
		log.Fatal(err)
	}
	return groups
}

func getManagers(group string, account string) []User {
	var cmd exec.Cmd
	if account != "" {
		cmd = *exec.Command("op", "--format", "json", "--account", account, "group", "user", "list", group)
	} else {
		cmd = *exec.Command("op", "--format", "json", "group", "user", "list", group)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var users []User
	err = json.Unmarshal([]byte(out), &users)
	if err != nil {
		log.Fatal(err)
	}
	var managers []User
	for _, user := range users {
		if user.Role == "MANAGER" && user.State == "ACTIVE" {
			managers = append(managers, user)
		}
	}

	return managers
}

func printOutput(group Group, managers []User, csv bool) {
	if csv {
		for _, manager := range managers {
			fmt.Printf("%s,%s,%s\n", group.Name, manager.Name, manager.Email)
		}
	} else {
		if len(managers) > 0 {
			fmt.Println(group.Name)
			for _, manager := range managers {
				fmt.Printf("%s\t%s\n", manager.Name, manager.Email)
			}
			fmt.Println()
		}
	}
}

func main() {

	var account, vault string
	var csv bool

	flag.BoolVar(&csv, "csv", false, "Print in CSV format")
	flag.StringVar(&account, "account", "", "1Password account shorthand")
	flag.Parse()
	if len(flag.Args()) <= 0 {
		fmt.Println("Please include the name or id of a vault")
		os.Exit(1)
	} else {
		vault = flag.Args()[0]
	}

	groups := getVaultGroups(vault, account)
	for _, group := range groups {
		managers := getManagers(group.Name, account)
		printOutput(group, managers, csv)
	}
}
