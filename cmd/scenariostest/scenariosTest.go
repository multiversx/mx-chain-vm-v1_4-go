package main

import (
	scenclibase "github.com/multiversx/mx-chain-scenario-go/clibase"
	mc "github.com/multiversx/mx-chain-scenario-go/scenario/io"

	vmscenario "github.com/multiversx/mx-chain-vm-v1_4-go/scenario"
	cli "github.com/urfave/cli/v2"
)

var _ scenclibase.CLIRunConfig = (*vm14Flags)(nil)

func main() {
	scenclibase.ScenariosCLI("VM 1.4 internal", &vm14Flags{})
}

type vm14Flags struct{}

func (*vm14Flags) GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "force-trace-gas",
			Aliases: []string{"g"},
			Usage:   "overrides the traceGas option in the scenarios`",
		},
	}
}

func (*vm14Flags) ParseFlags(cCtx *cli.Context) scenclibase.CLIRunOptions {
	runOptions := &mc.RunScenarioOptions{
		ForceTraceGas: cCtx.Bool("force-trace-gas"),
	}

	vmBuilder := vmscenario.NewScenarioVMHostBuilder()

	return scenclibase.CLIRunOptions{
		RunOptions: runOptions,
		VMBuilder:  vmBuilder,
	}
}
