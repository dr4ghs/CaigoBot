package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID          uint64 `gorm:"primaryKey"`
	DiscordUser *DiscordUser
}

type DiscordUser struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null"`
	DiscordID string `gorm:"type:varchar(18)"`
}

func (DiscordUser) TableName() string {
	return "discord.users"
}

type DiscordGuild struct {
	ID               uint64 `gorm:"primaryKey"`
	DiscordID        string
	CommandChannelID string
	Members          []*DiscordUser    `gorm:"many2many:discord.guild_members"`
	Commands         []*DiscordCommand `gorm:"many2many:discord.guild_commands"`
	StreamRooms      []*DiscordStreamRoom
}

func (DiscordGuild) TableName() string {
	return "discord.guilds"
}

type DiscordMember struct {
	ID         uint64 `gorm:"primaryKey"`
	GuildID    uint
	MemberID   uint
	Roles      []*DiscordRole     `gorm:"many2many:discord.member_roles"`
	StreamRoom *DiscordStreamRoom `gorm:"references:OwnerID"`
}

func (DiscordMember) TableName() string {
	return "discord.members"
}

type DiscordRole struct {
	ID           uint64 `gorm:"primaryKey"`
	GuildID      uint64
	DiscordID    string
	Members      []*DiscordMember `gorm:"many2many:discord.member_roles"`
	StreamRoomID int
}

func (DiscordRole) TableName() string {
	return "discord.roles"
}

type DiscordCommand struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string
}

func (DiscordCommand) TableName() string {
	return "discord.commands"
}

type DiscordStreamRoom struct {
	ID                uint64 `gorm:"primaryKey"`
	GuildID           uint64
	OwnerID           uint64
	StaffRoleID       *DiscordRole
	SubscribersRoleID *DiscordRole
	FollowersRoleID   *DiscordRole
	AcceptPolicy      int
	JoinPolicy        int
}

func (DiscordStreamRoom) TableName() string {
	return "discord.stream_rooms"
}

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", host, port, name, user, pass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&User{},
		&DiscordUser{},
		&DiscordGuild{},
		&DiscordMember{},
		&DiscordCommand{},
	)
}
