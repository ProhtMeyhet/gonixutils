package whoami

import(
	"os/user"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func WhoAmI(input *Input) (exitCode uint8) {
	output := abstract.NewOutput(input.Stdout, input.Stderr)
	exitCode = CurrentUser(input, output)
	output.Done(); output.Wait(); return
}

func CurrentUser(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	current, e := user.Current(); if output.WriteE(e) {
		exitCode = ERROR_USER_NOT_FOUND; goto out
	}

	output.WriteFormatted("%s\n", current.Username)

out:
	return
}
