package main

import (
	"github.com/mslmio/cli/lib/complete"
	"github.com/mslmio/cli/lib/complete/predict"
)

var completions = &complete.Command{
	Sub: map[string]*complete.Command{
		"emailverify": completionsEmailVerify,
		"signup":      completionsSignup,
		"login":       completionsLogin,
		"completion":  completionsCompletion,
		"version":     completionsVersion,
	},
	Flags: map[string]complete.Predictor{
		"-v":        predict.Nothing,
		"--version": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
