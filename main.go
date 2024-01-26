package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/joho/godotenv"
)

var botToken string

type TelegramChatMember struct {
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews"`
	CanBeEdited           bool   `json:"can_be_edited"`
	CanChangeInfo         bool   `json:"can_change_info"`
	CanDeleteMessages     bool   `json:"can_delete_messages"`
	CanEditMessages       bool   `json:"can_edit_messages"`
	CanInviteUsers        bool   `json:"can_invite_users"`
	CanJoinGroups         bool   `json:"can_join_groups"`
	CanPinMessages        bool   `json:"can_pin_messages"`
	CanPostMessages       bool   `json:"can_post_messages"`
	CanPromoteMembers     bool   `json:"can_promote_members"`
	CanReadMessages       bool   `json:"can_read_messages"`
	CanRestrictMembers    bool   `json:"can_restrict_members"`
	CanSendMediaMessages  bool   `json:"can_send_media_messages"`
	CanSendMessages       bool   `json:"can_send_messages"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages"`
	CanSendPolls          bool   `json:"can_send_polls"`
	Status                string `json:"status"`
	UntilDate             int64  `json:"until_date"`
	User                  struct {
		FirstName string `json:"first_name"`
		ID        int64  `json:"id"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"user"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken = os.Getenv("TOKEN")
	if botToken == "" {
		panic("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewCommand("start", start))

	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	updater.Idle()
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	channelID := int64(-1002013563867)
	userID := int64(ctx.EffectiveUser.Id)

	chatMember, err := b.GetChatMember(channelID, userID, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	chatMemberJSON, jsonErr := json.Marshal(chatMember)
	if jsonErr != nil {
		log.Println("Error:", jsonErr)
		return jsonErr
	}

	var telegramChatMember TelegramChatMember

	reader := bytes.NewReader(chatMemberJSON)

	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&telegramChatMember)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	userName := telegramChatMember.User.Username
	status := telegramChatMember.Status

	log.Printf("Username: %s | Status: %s", userName, status)

	_, sendErr := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("<b>Username:</b> @%s\n<b>Status:</b> %s\n", userName, status), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if sendErr != nil {
		return fmt.Errorf("failed to send start message: %w", sendErr)
	}
	return nil
}
