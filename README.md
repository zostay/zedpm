# ZedPM

ZedPM aims to be a general purposes project management tool aimed at helping
developers maintain their various and sundry Golang projects with minimal effort
and thinking, and no Makefiles for common tasks, but with massive amounts of raw
power.

> **Dr. McKay:** ...And we'll need the Zed-P-M.  
> **Gen. O'Neill:** What?  
> **Dr. Jackson:** The ZPM. He's...he's Canadian.  
> **Gen. O'Neill:** (at McKay) I'm sorry.  
> **Dr. McKay:** The Zero Point Module, General. The ancient power source you
> recovered from Proclarush Toanas and that's now powering the outpost defenses.
> I've since determined it generates its enormous power from vacuum energy
> derived from a self-contained region of subspace time.  
> **Gen. O'Neill:** That was a waste of a perfectly good explanation.  
> 
> *Stargate Atlantis, Season 1, Episode 1, "Rising, Part 1"*

# Solving Common Problems

This tool aims at solving common problems in a way that lets developers run
self-documenting commands that do the right thing. Since any two developers will
likely disagree on the correct way to do these things and, depending on aspects
of a project like what kind of deliverables it produces, how it used, the 
target audience, and so on, TheRightWay™ varies from project to project. 
However, the goals you want to achieve are largely the same:

 * Every project needs to be checked for correctness: syntax needs to be 
   checked, anti-patterns need to be detected, and semantics need to be 
   verified.
 * Many projects need to generate some section of their code.
 * Any serious project needs a release process. Libraries need to be tagged so
   that godoc sites can index them. Client-side applications need to generate
   binaries that can be downloaded and installed. Server-side applications need
   to push container images. And so on.
 * Creating projects requires a common set of boilerplate that is tedious to
   setup every time. Typically a project needs to have git setup, go mod init
   run, a Changes file, a README, a license, maybe some common boilerplate code
   based on the kind of application or library, etc. It'd be nice if there was a
   tool that could at least get you past crawling and to the walking phase
   whenever starting a new prototype or project.

This tool aims at doing all of these things.

# Goals, Tasks, and Plugins, Oh my!

The tool operates by executing goals. Each goal is made up of a number of tasks.
These tasks are provided by plugins. The plugins are chosen and configured by
you. For complex projects, plugins can also be configured to perform the same
goal with multiple configurations by having configured targets.

## Common Goals

Here are listed the common goals that are built-in to zedpm. Any plugin may
define additional goals, if desired.

### Init (not yet implemented)

The init goal will initialize a new project and fill it with standard 
boilerplate and initialize files and configurations.

### Generate (not yet implemented)

The generate goal is will take source code in one form and generate it in
another. This is useful for building gRPC stubs, web sites built using
node-based tooling (or whatever), embedding files and resources, generating ORM
code, etc.

### Build (not yet implemented)

The build goal will check the syntax of the source code and ensure that it will
run. It will not guarantee that the code runs correctly and will not generate
source code from other source code.

### Lint

Linting runs programs that perform correctness checks against your code and
other files. This can run a static analyzer to verify that your code does not
implement certain anti-patterns that are prone to cause bugs.

### Test

Testing runs some set of test suites against the project to verify that the code
is working correctly.

### Request (not yet implemented)

Request (a.k.a., pull-request or merge-request) is the act of declaring that
some set of changes is ready to be considered for merger into some other target
branch (usually main or master) of the project.

### Install (not yet implemented)

The install goal will build local artifacts and install them. This is intended
to be a helper to allow binaries and other artifacts to be installed as part of
the development process.

### Release

Releasing is the process of tagging a set of changes and setting up the release
process.

### Deploy (not yet implemented)

Deploy will construct and deliver artifacts to a destination, such as Docker
Hub or another container repository.

### Info

Info retrieves information about the project and can be used to either query the
state of the application or provide zedpm plugin derived information to other
tooling.

## Built-in Plugins

The zedpm project has the following built-in plugins:

### zedpm-plugin-changelog

This provides tasks for linting the correctness of a changelog, for extracting
the changes related to a given release version, and for preparing and fixing the
changelog for release.

This is currently done according to an opinionated format that is subject to
change in the future.

### zedpm-plugin-git

This provides tasks for creating a release branch and tagging the release
according to a semantic version.

### zedpm-plugin-github

This provides tasks for creating pull requests during release, awaiting for
CI/CD tests to complete successfully, merging the pull request, and creating the
release.

### zedpm-plugin-go

This provides tools for accessing aspects of the go command for various zedpm
commands.

### zedpm-plugin-goals

This provides the definition for the zedpm built-in goals and related support
code and tasks.

# Inspiration

At this time, I think it would be wrong for me not to mention that this project
was inspired by other tools. I spent most of the past 20 years working in the
Perl programming language. When I came to Golang, one thing became immediately
apparent: while the `go` command itself does slightly more and useful things
than the `perl` command, it is vastly inadequate for performing the full range
of common dev tasks. The Perl community built rich project management tools 
aimed at enabling developers to quickly and easily build new prototypes, check 
the code for correctness, and release that code for others to use.

Tools like [ExtUtils::MakeMaker][mm], [Module::Build][mb], and [Dist::Zilla][dz]
provided much of the inspiration for this tool.  ExtUtils::MakeMaker is a very 
old tool that allows you to quickly generate an entire Makefile customized for 
your project based upon configuration contained in a simple Perl program. 
Module::Build extended this to alleviate the need for the Makefile when such was 
unnecessary. And Dist::Zilla provided additional tools for managing module 
dependencies, releasing and deploying software, and generating boilerplate, such 
as ensuring that all source files have the latest license attached or that each 
have complete and correctly formatted documentation.

[mm]: <https://metacpan.org/pod/ExtUtils::MakeMaker>
[mb]: <https://metacpan.org/pod/Module::Build>
[dz]: <https://metacpan.org/pod/Dist::Zilla>

These tools went a long way to ensuring that Perl developers did not waste a lot
of time remembering commands, mucking about with esoteric Makefile syntax, and
generally wasting their time on the more tedious aspects of developing a library
or application.

# Copyright & License

Copyright 2023 Andrew Sterling Hanenkamp.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the “Software”), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
