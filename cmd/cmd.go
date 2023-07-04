package cmd

import "juicer/gameserver"

func Run() error {
	return gameserver.RunServer()
}
