//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package module

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gvallee/go_exec/pkg/advexec"
)

// ToEnv returns environment variables specified in 'envVars' that the list of 'modules' sets.
// It returns a map where the keys are the environment variables that are requested and the values a slice of string representing each element being added to the environment variable that is set by the module(s).
func ToEnv(envVars []string, modules []string) (map[string][]string, error) {
	// Get the environment before and after loading the module
	// and then get the difference.

	moduleEnvVars := make(map[string][]string)
	initialEnvVars := make(map[string]string)
	for _, v := range envVars {
		initialEnvVars[v] = os.Getenv(v)
	}

	// Generate the script to get the env vars after loading the modules
	fd, err := ioutil.TempFile("", "module-to-env_")
	if err != nil {
		return nil, err
	}

	filepath := fd.Name()
	err = os.Chmod(filepath, 0777)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filepath)

	_, err = fd.WriteString("#!/bin/bash\n")
	if err != nil {
		return nil, err
	}
	_, err = fd.WriteString("#\n\n")
	if err != nil {
		return nil, err
	}
	_, err = fd.WriteString("module load " + strings.Join(modules, " ") + "\n")
	if err != nil {
		return nil, err
	}
	for _, envvar := range envVars {
		_, err = fd.WriteString("echo " + envvar + "=$" + envvar + "\n")
		if err != nil {
			return nil, err
		}
	}
	fd.Close()

	var moduleCmd advexec.Advcmd
	moduleCmd.BinPath = filepath
	res := moduleCmd.Run()

	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			// It is most certainly some output from the module loading script that we do not care about
			continue
		}
		currentEnvVar := tokens[0]
		currentEnvVarValue := tokens[1]
		if _, ok := initialEnvVars[currentEnvVar]; !ok {
			// This is not a variable we were looking for
			continue
		}
		moduleEnvVars[currentEnvVar] = strings.Split(strings.Replace(currentEnvVarValue, initialEnvVars[currentEnvVar], "", 1), ":")
	}

	return moduleEnvVars, nil
}
