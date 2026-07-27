package main

var (
	commit  = "commit"
	version = "version"
)

func main() {

	serverCli := cli{
		commit:  commit,
		version: version,
	}

	serverCli.run()
}
