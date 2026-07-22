package main

import bootstrapadmin "cointrade/internal/bootstrap/admin"

func main() {
	bootstrapadmin.Run(bootstrapadmin.OptionsFromEnv())
}
