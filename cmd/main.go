package main

import (
	issuectl "github.com/janekbaraniewski/issuectl/pkg"
)

func main() {
	if err := issuectl.StartWorkingOnIssue("test-automated-issue"); err != nil {
		issuectl.Log.Infof("%v", err)
	}
}
