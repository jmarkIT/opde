package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const usage = `Usage:
	opde [options] (vault)
Options:
	-a, --account ACCOUNT	Shortname for 1Password account to use
	-c, --csv		Print output in csv format`

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
	Members     []GroupMember
}

func (g *Group) setMembers(account string) {
	var cmd exec.Cmd
	group := g.Name
	if account != "" {
		cmd = *exec.Command("op", "--format", "json", "--account", account, "group", "user", "list", group)
	} else {
		cmd = *exec.Command("op", "--format", "json", "group", "user", "list", group)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var users []GroupMember
	err = json.Unmarshal([]byte(out), &users)
	if err != nil {
		log.Fatal(err)
	}

	g.Members = users
}

func (g *Group) getManagers(account string) []GroupMember {
	var managers []GroupMember
	for _, user := range g.Members {
		if user.Role == "MANAGER" && user.State == "ACTIVE" {
			managers = append(managers, user)
		}
	}

	return managers
}

type GroupMember struct {
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

func printOutput(group Group, managers []GroupMember, csv bool) {
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

	flag.BoolVar(&csv, "c", false, "Print in CSV format")
	flag.BoolVar(&csv, "csv", false, "Print in CSV format")
	flag.StringVar(&account, "a", "", "1Password account shorthand")
	flag.StringVar(&account, "account", "", "1Password account shorthand")
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()
	if len(flag.Args()) <= 0 {
		fmt.Println("Please include the name or id of a vault")
		os.Exit(1)
	} else {
		vault = flag.Args()[0]
	}

	groups := getVaultGroups(vault, account)
	for _, group := range groups {
		group.setMembers(account)
		managers := group.getManagers(account)
		printOutput(group, managers, csv)
	}
}
