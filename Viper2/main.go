package main

import (
	"log"

	"viper2/cmd"
)

func main() {
	// Print an ASCII banner on startup (logged via standard logger)
	printBanner()

	if err := cmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

// printBanner logs a small Bee Message ASCII art banner.
func printBanner() {
	log.Println(`
	   .--.
	  /    \   Bee Message
	 /  /\  \
	/  /  \  \
   (  (    )  )
	\  \  /  /   "Hello!"
	 '--\/--'  
	  (o  o)
	  /|/\|\
	 /_/  \_\
	`)
}
