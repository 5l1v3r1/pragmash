# pragmash

**pragmash** is a simple, flexible scripting language. It is made to be dropped into other systems (i.e. build systems, scripting environments, etc.).

The runtime&mdash;if you can call it that&mdash;is up and running. I am now working to build up the standard library until the language is usable.

# Installing

Pragmash has no external dependencies and can be installed very easily:

    go get github.com/unixpickle/pragmash
    go install github.com/unixpickle/pragmash/pragmash

# Learning

To learn the syntax of pragmash, checkout [SYNTAX.md](SYNTAX.md).

To see the commands you can currently use in a pragmash program, see [COMMANDS.md](COMMANDS.md).

To see some pre-written example programs, see [demo](demo).

# TODO

Here are the commands that I'd like to add:

 * File system manipulation functions
   * Chmod
   * Touch
   * Link
 * Time manipulation functions
   * Parse/format dates
 * Strings
   * Is letter
   * Is digit
   * Add escmatch command for match with escapes
 * Math
   * Inverse trig
   * Sqrt
 * Networking
   * Cookies
   * POST
 * Add option to get current index in for loop

Here are improvements I'd like to make to the tools and/or implementation:

 * Support up arrow in REPL
 * Unify standard variables
