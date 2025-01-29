/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"livekit-keysrv/internal/livekit"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TokenStruct struct {
	Token string `json:"token"`
}

// serveTokenCmd represents the serveToken command.
var serveTokenCmd = &cobra.Command{
	Use:   "serveToken",
	Short: "A brief description of your command",
	Long:  `Serves a 24h livekit access token via http.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serveToken called")

		apiKey := viper.GetString("API_KEY")
		if apiKey == "" {
			log.Fatalf("API_KEY not set")
		}

		apiSecret := viper.GetString("API_SECRET")
		if apiSecret == "" {
			log.Fatalf("API_SECRET not set")
		}

		tkngen := livekit.NewDefaultTokenGenerator(apiKey, apiSecret)

		route := mux.NewRouter()

		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		})

		route.Handle("/token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			//if r.Method != http.MethodGet {
			//	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			//	w.WriteHeader(http.StatusMethodNotAllowed)
			//	return
			//}
			identity := r.URL.Query().Get("identity")
			if identity == "" {
				http.Error(w, "identity not set", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			room := r.URL.Query().Get("room")
			if room == "" {
				http.Error(w, "room not set", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Printf("generating token for identity %s in room %s", identity, room)
			token, err := tkngen.GenerateToken(identity, room)
			if err != nil {
				return
			}

			tokenObject := TokenStruct{Token: token}

			jsonToken, err := json.Marshal(tokenObject)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(jsonToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		})).Methods(http.MethodGet)

		if err := http.ListenAndServe(":3030", cors.Handler(route)); err != nil {
			log.Fatalf("error starting http server: %v\n", err.Error())
		}
		fmt.Println("server shutdown")
	},
}

func init() {
	rootCmd.AddCommand(serveTokenCmd)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveTokenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveTokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
