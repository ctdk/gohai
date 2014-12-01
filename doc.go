/*
Gohai is a golang implementation of ohai (https://github.com/opscode/ohai). It uses a combination of external shell programs and standard golang libraries to gather data about a system and print it out in JSON format.

Installation

Assuming you have Go installed and configured, just run

   go get github.com/go-chef/gohai

to install gohai.

Usage

   gohai

This will print out a bunch of information about the host in JSON format.

Cross-compiling

To build the binary file for several platforms, we use goxc:

	$ go get github.com/laher/goxc
	$ goxc -bc='linux,darwin,windows' -d=[BUILD_DIR]

Plugins

Gohai supports plugins. They can be written in any language, so long as they can print a JSON hash to standard out. The structure of the JSON output must include all levels of the keys it's supposed to have. If your plugin provides "toplevel", "foo/bar" and "foo/quuz", it should look something like this:

	{
		"toplevel": "hi!",
		"foo": {
			"bar": {
				"beep": "boop"
			},
			"quuz": [
				1,
				2,
				3
			]
		}
	}

The gohai plugins must be placed in the plugin directory. The default (and currently only) location for the plugin directory is /var/lib/gohai/plugins. Making this configurable is on the TODO list.

Submodule interface

Adding new data collector submodules is easy. The new module has to satisfy the "collector" interface defined in gohai.go - Name(), which returns a string naming the submodule, Collect(), which returns an interface{} filled with the collected data and an error, and Provides(), which returns an array of strings describing the data collected by the submodule. The `cpu` module is a good example of how they work.

If you're making a submodule "foo", it should go in a directory named foo, with its first file being `foo/foo.go`. This file should contain the functions for the interface. The Collect() function can perform any platform independent work, but platform specific functionality should be put in a separate function named something like `getFoo()`. Then, in separate files, define the platform specific versions of that function, using build tags to only build that version of the function for the right platform. After that's all done, add the new data collector to the `collectors` slice defined in `gohai.go`.

See Also

* facter (https://github.com/puppetlabs/facter)

* ohai (https://github.com/opscode/ohai)


TODO

There's a lot that needs to be done with gohai still. Here is a partial list:

* Only OS X data collection is well defined right now. It needs to support otherplatforms like Linux, Windows, etc.

* Should be able to build easily with gccgo. This would make it easier to extend gohai to run on stranger platforms that have gccgo available, but not the standard golang compiler.

* There's still a few bits of data ohai has by default that need to be added into gohai, and some others ought to be shifted around a bit.

* Needs a --help flag, and needs to be able to at least specify a non-default plugin directory.

* Presumably even more things that are escaping me right now.

Authors

Kentaro Kuribayashi (http://kentarok.org/)
Jeremy Bingham (http://time.to.pullthepl.ug)
Various DataDog people (fill in their names later)

Gohai was originally forked from Kentaro Kubayashi's verity (https://github.com/kentaro/verity) by folks at DataDog, but has been diverging from that since then.

LICENSE

MIT. See the LICENSE file for details.

*/
package main
