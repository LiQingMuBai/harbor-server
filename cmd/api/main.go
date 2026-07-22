package main

import bootstrapapi "cointrade/internal/bootstrap/api"

func main() {
	bootstrapapi.Run(bootstrapapi.OptionsFromEnv())
}
