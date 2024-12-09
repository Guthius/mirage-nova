package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/config"
)

// IsBanned checks if a player is banned
func IsBanned(ipAddr string) bool {
	file, err := os.OpenFile("banlist.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("failed to open banlist.txt (%s)", err)
		}
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ipAddr) {
			return true
		}
	}

	return false
}

// BanPlayerBy bans a player from the server
func BanPlayerBy(p *PlayerData, bannedBy string) {
	if !BanIP(p.Connection.RemoteAddr(), bannedBy) {
		return
	}

	if p.Character != nil {
		SendGlobalMessage(fmt.Sprintf("%s has been banned from %s by %s!", p.Character.Name, config.GameName, bannedBy), color.White)

		log.Printf("%s has been banned by %s", p.Character.Name, bannedBy)
	}

	SendAlert(p, fmt.Sprintf("You have been banned by %s!", bannedBy))
}

// BanPlayer bans a player from the server
func BanPlayer(p *PlayerData) {
	if !BanIP(p.Connection.RemoteAddr(), "Server") {
		return
	}

	if p.Character != nil {
		SendGlobalMessage(fmt.Sprintf("%s has been banned from %s by the Server!", p.Character.Name, config.GameName), color.White)

		log.Printf("%s has been banned by the Server.", p.Character.Name)
	}

	SendAlert(p, "You have been banned by the Server!")
}

// BanIP bans an IP address from the server
func BanIP(ipAddr string, bannedBy string) bool {
	if IsBanned(ipAddr) {
		return false
	}

	file, err := os.OpenFile("banlist.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open banlist.txt (%s)", err)
		return false
	}

	defer file.Close()

	_, err = fmt.Fprintf(file, "%s;%s\n", ipAddr, bannedBy)
	if err != nil {
		log.Printf("failed to write to banlist.txt (%s)", err)
		return false
	}

	return true
}
