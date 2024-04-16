package utils

import "github.com/urfave/cli/v3"

func CommandBoolSetAndValue(cmd *cli.Command, name string) (bool, bool) {
	return cmd.IsSet(name), cmd.Bool(name)
}
