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

	PrintUser(current, input, output)

out:
	return
}

func PrintUser(user *user.User, input *Input, output abstract.OutputInterface) {
	if input.PrintAll {
		output.WriteFormatted("%s (%s)\n%s\n%s [%s]\n",
						user.Username,
						user.Name,
						user.HomeDir,
						user.Uid,
						user.Gid,
					    )
	} else {
		if input.PrintUsername {
			output.WriteFormatted("%s\n", user.Username)
		}

		if input.PrintName {
			output.WriteFormatted("%s\n", user.Name)
		}

		if input.PrintHome {
			output.WriteFormatted("%s\n", user.HomeDir)
		}

		if input.PrintUid {
			output.WriteFormatted("%s\n", user.Uid)
		}

		if input.PrintGid {
			output.WriteFormatted("%s\n", user.Gid)
		}
	}
}
