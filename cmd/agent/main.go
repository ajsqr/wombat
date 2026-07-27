package main

var (
	commit  = "commit"
	version = "version"
)

func main() {

	agentCli := cli{
		version: version,
		commit:  commit,
	}

	agentCli.run()
}
