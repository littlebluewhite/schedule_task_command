package main

import (
	"flag"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	config "schedule_task_command/util/config"
	"schedule_task_command/util/migrate"
)

func main() {
	var up, down, to, t bool
	var version int
	flag.BoolVar(&up, "up", false, "up to newest")
	flag.BoolVar(&t, "test", false, "for test")
	flag.BoolVar(&down, "down", false, "down to oldest")
	flag.BoolVar(&to, "to", false, "to version")
	flag.IntVar(&version, "version", -1, "version")
	flag.Parse()
	var c config.DBConfig
	if t {
		c = config.NewConfig[config.DBConfig](".", "env", "db_test")
	} else {
		c = config.NewConfig[config.DBConfig](".", "env", "db")
	}

	client := migrate.New(c)
	if up {
		client.Up()
	}
	if down {
		client.Down()
	}
	if to {
		if version != -1 {
			client.To(uint(version))
		}
	}
}
