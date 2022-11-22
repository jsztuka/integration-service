package snapshot

import (
	"fmt"
)

type EnvVarPair struct {

	// Name is the environment variable name
	Name string `json:"name"`

	// Value is the environment variable value
	Value string `json:"value"`
}

type Binder struct {
	Env []EnvVarPair `json:"env,omitempty"`
}

func main() {
	//section from Env
	env1 := EnvVarPair{"kubeconf", "/home/neco"}
	env2 := EnvVarPair{"gobin", "/bin/go"}
	envy := []EnvVarPair{env1, env2}
	blind := Binder{envy}

	//section from IntegrationTestScenario
	envE1 := EnvVarPair{"kubeconf", "/home/neco/ahoj"}
	envE2 := EnvVarPair{"tknbin", "/bin/tkn"}
	envyE := []EnvVarPair{envE1, envE2}
	blindI := Binder{envyE}

	//here magic happen
	// If ITS is empty and Env has something
	// If env is empty and ITS has something

	//print array of envs for ITS
	fmt.Println(blind)
	//print array of envs for ENV
	fmt.Println(blindI)
	fmt.Println("=====================================")

	fmt.Println("DATA BEFORE CHANGE")
	fmt.Println(blind)

	blindIZer := Binder{envyE}
	blindEZer := Binder{}

	for v := range blindI.Env {
		match := false
		if len(blindIZer.Env) == 0 {
			break
		} else if len(blindEZer.Env) == 0 {
			blindEZer.Env = blindI.Env
			break
		}
		for p := range blind.Env {
			if blindI.Env[v].Name == blind.Env[p].Name {
				match = true
				blind.Env[p].Value = blindI.Env[v].Value
				fmt.Println(blind.Env[p].Value)
			}
			if !match && (p == len(blind.Env)-1) {
				blind.Env = append(blind.Env, EnvVarPair{blindI.Env[v].Name, blindI.Env[v].Value})
			}
		}
	}

	fmt.Println("DATA AFTER CHANGE")
	fmt.Println(blind)
	fmt.Println("EZER")
	fmt.Println(blindEZer)

	blindT := Binder{}
	fmt.Println(blindT)

}
