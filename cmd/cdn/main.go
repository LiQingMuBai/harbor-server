package main

import bootstrapcdn "cointrade/internal/bootstrap/cdn"

func main() {
	bootstrapcdn.Run(bootstrapcdn.OptionsFromEnv())
}
