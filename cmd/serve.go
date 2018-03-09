package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aflog/assignment-messagebird/message"
	"github.com/aflog/assignment-messagebird/message/client"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyHost   = "host"
	keyPort   = "port"
	keyAPIKey = "apikey"
)

var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Runs assignment API",
	Long: `MessageBird assignment API
The server listens for incoming requests for sending a message, reads the json
data and redirects it to Message Bird API. In case the message exceeds 160
caracters it will send it as a concatenated message.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString(keyAPIKey)
		if apiKey == "" {
			log.Fatalf("value of %s is not set", keyAPIKey)
		}
		fmt.Println(apiKey)
		mbc := client.NewMessageBird(apiKey)
		defer mbc.Close()

		s, err := message.NewHandler(mbc)
		if err != nil {
			log.Fatalf("unable to create message handler: %v", err)
			return
		}

		r := mux.NewRouter()
		r.HandleFunc("/send", s.Procces).Methods(http.MethodPost)

		h := viper.GetString(keyHost)
		p := viper.GetString(keyPort)
		addr := h + ":" + p

		log.Println("Starting Service")
		srv := &http.Server{
			Handler:      r,
			Addr:         addr,
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  30 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	cmdRoot.AddCommand(cmdServe)
	cmdServe.Flags().StringP(keyHost, "H", "", "server address")
	cmdServe.Flags().StringP(keyPort, "p", "8080", "port address")
	cmdServe.Flags().String(keyAPIKey, "", "Message Bird API key")
	err := viper.BindPFlags(cmdServe.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
