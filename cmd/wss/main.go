package main

import bootstrapwss "cointrade/internal/bootstrap/wss"

func main() {
	bootstrapwss.Run(bootstrapwss.OptionsFromEnv())
}
