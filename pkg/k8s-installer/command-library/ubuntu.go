package command_library

const (
	dpkg command = "dpkg-query"
	grep command = "grep"

	//nolint:lll // complex line to be handled by lll
	dpkgKeys = `-W -f='${Status} {"name":"${Package}","realname":"${Source}","version":"${Version}","main":"${Maintainer}","arch":"${Architecture}"}` + "\n" + `'`
)

var (
	newUbuntuGatherer = withBook(ubuntuCommandBook)
	ubuntuCommandBook = merge(unixCommandBook, pureUbuntuCommandBook)

	pureUbuntuCommandBook = map[gatheringStage][]commandAndParser{
		os:      {},
		network: {},
		software: {
			{
				command:   dpkg.WithArgs(dpkgKeys).Pipe(grep.WithArgs(`"^install ok installed"`)),
				parser:    UnmarshalDebianSoftware,
				condition: anyway,
			},
			// TODO: add apt
		},
		hardware: {},
	}
)
