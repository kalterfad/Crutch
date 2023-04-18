package crutchClient

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func StartClient() {
	sock := newSocket()
	err := sock.connectToSocket()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	path, err := checkArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	app := cli.NewApp()
	app.Name = "Crutch client"
	app.Usage = "This program is designed to scan directories and help you manage files. It allows you to " +
		"constantly monitor the appearance of new files and automatically sort them into folders.\n."

	crutch := &crutchService{path: path, sock: sock}
	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add dir to scan list",
			Action: func(c *cli.Context) {
				err = crutch.add(askUser())
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			},
		},
		{
			Name:  "rm",
			Usage: "rm dir from scan list",
			Action: func(c *cli.Context) {
				err = crutch.rm()
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			},
		},
	}
	_ = app.Run(os.Args)

	response, err := sock.readSock()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(response)

	sock.closeConn()
}
