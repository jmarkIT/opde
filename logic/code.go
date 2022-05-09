package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type Vault struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Content_Version int    `json:"content_id"`
	Groups          []Group
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

func (v *Vault) SetGroups(account string) {
	var cmd exec.Cmd
	vault := v.Name
	if account != "" {
		cmd = *exec.Command("op", "--format", "json", "--account", account, "vault", "group", "list", vault)
	} else {
		cmd = *exec.Command("op", "--format", "json", "vault", "group", "list", vault)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var groups []Group
	err = json.Unmarshal([]byte(out), &groups)
	if err != nil {
		log.Fatal(err)
	}

	v.Groups = groups
}

func (g *Group) SetMembers(account string) {
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

func (g *Group) GetManagers(account string) []GroupMember {
	var managers []GroupMember
	for _, user := range g.Members {
		if user.Role == "MANAGER" && user.State == "ACTIVE" {
			managers = append(managers, user)
		}
	}

	return managers
}

func (g *Group) GetMembers(account string) []GroupMember {
	var members []GroupMember
	for _, user := range g.Members {
		if user.State == "ACTIVE" {
			members = append(members, user)
		}
	}

	return members
}

type GroupMember struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
	State string `json:"state"`
	Role  string `json:"role"`
}

func GetVaultGroups(vault string, account string) []Group {
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

func PrintGroupManagers(group Group, managers []GroupMember, csv bool) {
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

func PrintGroupMembers(group Group, members []GroupMember, csv bool) {
	if csv {
		for _, member := range members {
			fmt.Printf("%s,%s,%s\n", group.Name, member.Name, member.Email)
		}
	} else {
		if len(members) > 0 {
			fmt.Println(group.Name)
			for _, manager := range members {
				fmt.Printf("%s\t%s\n", manager.Name, manager.Email)
			}
			fmt.Println()
		}
	}
}
