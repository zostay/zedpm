package goals

const (
	NameBuild    = "build"
	NameDeploy   = "deploy"
	NameGenerate = "generate"
	NameInfo     = "info"
	NameInit     = "init"
	NameInstall  = "install"
	NameLint     = "lint"
	NameRelease  = "release"
	NameRequest  = "request"
	NameTest     = "test"
)

// DescribeBuild describes the build goal, which is primarily aimed at checking
// syntax of source files.
func DescribeBuild() *GoalDescription {
	return &GoalDescription{
		name:  NameBuild,
		short: "Syntax check and prepare for development.",
	}
}

// DescribeDeploy describes the deploy goal, which is aimed at saving or pushing
// artifacts prior to release.
func DescribeDeploy() *GoalDescription {
	return &GoalDescription{
		name:  NameDeploy,
		short: "Deploy software to a remote server.",
	}
}

// DescribeGenerate describes the generate goal, which handles generated source
// code from other source code annotations and templates.
func DescribeGenerate() *GoalDescription {
	return &GoalDescription{
		name:  NameGenerate,
		short: "Perform code generation tasks.",
	}
}

// DescribeInfo describes the info goal, which is used to extract data from
// zedpm for informational purposes or for use with outside tooling.
func DescribeInfo() *GoalDescription {
	return &GoalDescription{
		name:  NameInfo,
		short: "Describe information about the project.",
	}
}

// DescribeInit describes the init goal, which performs functions to initialize
// new projects.
func DescribeInit() *GoalDescription {
	return &GoalDescription{
		name:  NameInit,
		short: "Initialize a new project directory.",
	}
}

// DescribeInstall describes the install goal, which installs artifacts onto the
// developer's local system.
func DescribeInstall() *GoalDescription {
	return &GoalDescription{
		name:  NameInstall,
		short: "Install software and assets locally.",
	}
}

// DescribeLint describes the lint goal, which performs analysis of files to
// check for errors and identify anti-patterns.
func DescribeLint() *GoalDescription {
	return &GoalDescription{
		name:    NameLint,
		short:   "Check files and data for errors and anti-patterns.",
		aliases: []string{"analyze"},
	}
}

// DescribeRequest describes the request goal, which provides tooling for
// creating pull-request/merge-request type operations.
func DescribeRequest() *GoalDescription {
	return &GoalDescription{
		name:    NameRequest,
		short:   "Request the merger of a code patch.",
		aliases: []string{"pull-request", "pr", "merge-request", "mr"},
	}
}

// DescribeRelease describes the release goal, which performs the action of
// releasing a new software version.
func DescribeRelease() *GoalDescription {
	return &GoalDescription{
		name:  NameRelease,
		short: "Mint and publish a release.",
	}
}

// DescribeTest describes the test goal, which performs the action of testing
// that the software is working correctly.
func DescribeTest() *GoalDescription {
	return &GoalDescription{
		name:  NameTest,
		short: "Run tests.",
	}
}
