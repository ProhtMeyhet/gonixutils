package ls

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// UNIX ls
func Ls(input *Input) (exitCode uint8) {
	output := abstract.NewSortedTabbedOutput(input.Stdout, input.Stderr)
	if input.Lines {
		output = abstract.NewOutput(input.Stdout, input.Stderr)
	} else if input.NoSort && !input.Detail {
		output = abstract.NewTabbedOutput(input.Stdout, input.Stderr)
	}

	exitCode = List(input, input.Paths, output)

	output.Done(); output.Wait()

	return
}

// Last week I walked into a local "home style cookin' restaurant/watering hole"
// to pick up a take out order. I spoke briefly to the waitress behind the counter,
// who told me my order would be done in a few minutes.
//
// So, while I was busy gazing at the farm implements hanging on the walls, I was approached
// by two, uh, um... well, let's call them "natives".
//
// These guys might just be the original Texas rednecks -- complete with ten-gallon hats,
// snakeskin boots and the pervasive odor of cheap beer and whiskey.
//
// "Pardon us, ma'am. Mind of we ask you a question?"
//
// Well, people keep telling me that Texans are real friendly, so I nodded.
//
// "Are you a Satanist?"
//
// And that, kids, is why the FreeBSD logo is a cute satan.
