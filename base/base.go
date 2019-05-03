package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	server "github.com/yamadashi/EscaTSGen"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env file failed : %s", err)
	}
}

func main() {
	var (
		addr   = flag.String("addr", ":8080", "addr to bind")
		dbconf = flag.String("dbconf", "dbconfig.yml", "database configuration file.")
		env    = flag.String("env", "development", "application envirionment (production, development etc.)")
	)
	flag.Parse()
	s := server.New()
	s.Init(*dbconf, *env)
	s.Run(*addr)
}
