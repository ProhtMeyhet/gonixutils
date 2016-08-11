package main

import(

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/rm"
)

func main() {
	flags := NewRmFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	gonixutils.Rm(input)
}
