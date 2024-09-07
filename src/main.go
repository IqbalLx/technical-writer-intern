package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IqbalLx/technical-writer-intern/src/modules/chat"
	"github.com/IqbalLx/technical-writer-intern/src/modules/document"
	"github.com/IqbalLx/technical-writer-intern/src/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {
	env := utils.NewEnv()

	// initialize database
	connString := env.Read("POSTGRES_CONNSTRING")
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, connString) // pgx.Connect(ctx, connOrDSNString)
	if err != nil {
		log.Fatalln("open: %w", err)
	}
	defer dbpool.Close()

	slackApi := slack.New(env.Read("SLACK_BOT_TOKEN"))
	slackSigningSecret := env.Read("SLACK_SIGNING_SECRET")
	slackBotID := "B07HBRXUQ83"

	http.HandleFunc("POST /slack/event", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sv, err := slack.NewSecretsVerifier(r.Header, slackSigningSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}

		w.WriteHeader(http.StatusOK)
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:

				// dont respond to bot's own message
				if ev.BotID == slackBotID {
					return
				}

				go func() {
					answer, err := chat.DoAnswerUserChat(context.Background(), dbpool, ev.User, ev.Text)
					if err != nil {
						slackApi.PostMessage(ev.Channel, slack.MsgOptionText("Bentar, error bang, coba lagi ya ntaran", false))
						return
					}

					slackApi.PostMessage(ev.Channel, slack.MsgOptionText(answer, false))
				}()
			}
		}
	})

	http.HandleFunc("POST /slack/slash-command/catetin", func(w http.ResponseWriter, r *http.Request) {
		sv, err := slack.NewSecretsVerifier(r.Header, slackSigningSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(io.TeeReader(r.Body, &sv))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		go func() {
			err = document.DoInsertNewDocument(context.Background(), dbpool, s.Text, s.UserName)
			if err != nil {
				errorResp := fmt.Sprintf("Waduh error bang <@%s>, coba lagi ya ntaran", s.UserID)
				slackApi.PostMessage(s.ChannelID, slack.MsgOptionText(errorResp, false))
				return
			}

			successResp := fmt.Sprintf("Okidoki bang <@%s>!, udah ku catet", s.UserID)
			slackApi.PostMessage(s.ChannelID, slack.MsgOptionText(successResp, false))
		}()

		// Create the response data
		response := map[string]string{
			"response_type": "in_channel",
			"text":          "Siap bang! Tak proses sik yo, abis ini kukabarin",
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Marshal the response data to JSON
		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		// Write the JSON data to the response
		w.Write(jsonData)
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}
