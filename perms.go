package main

import (
	"fmt"
)

func userPerms(uID string, gID string) int64 {
	member, err := client.GuildMember(gID, uID)
	if err != nil {
		fmt.Printf("member not found guild : %s user id : %s\n", gID, uID)
		return 0
	}
	var total int64
	for _, roleID := range member.Roles {
		role, err := client.State.Role(gID, roleID)
		if err != nil {
			fmt.Println("role not found", roleID)
			continue
		}
		total |= role.Permissions
	}
	return total
}
