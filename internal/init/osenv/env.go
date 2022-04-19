package osenv

import (
	"flag"
	"log"
	"os"
)

var (
	MONGO_HOST              string
	MONGO_PORT              string
	TEST_MODE               bool
	CHORME_HANDLESS_ADDRESS string
	CHORME_HANDLESS_PORT    string
)

func init() {

	flag.BoolVar(&TEST_MODE, "test", false, "test mode args")
	flag.Parse()
	if TEST_MODE {
		MONGO_HOST = "127.0.0.1"
		MONGO_PORT = "27017"
		return
	}

	// Mongo_host
	MONGO_HOST = os.Getenv("MONGO_HOST")
	if MONGO_HOST == "" {
		log.Fatal("env: mongo_host is empty")
	}

	// mongo_port
	MONGO_PORT = os.Getenv("MONGO_PORT")
	if MONGO_PORT == "" {
		log.Fatal("env: mongo_port is empty")
	}

	// CHORME_HANDLESS_ADDRESS
	CHORME_HANDLESS_ADDRESS = os.Getenv("CHORME_HANDLESS_ADDRESS")
	if CHORME_HANDLESS_ADDRESS == "" {
		log.Fatal("env: CHORME_HANDLESS_ADDRESS is empty")
	}

	// CHORME_HANDLESS_PORT
	CHORME_HANDLESS_PORT = os.Getenv("CHORME_HANDLESS_PORT")
	if CHORME_HANDLESS_PORT == "" {
		log.Fatal("env: CHORME_HANDLESS_ADDRESS is empty")
	}

}
