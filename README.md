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

This is my personal to-do list.

 * Create file system manipulation functions
   * Chmod
   * Touch
   * Link
   * FileType
 * Create time manipulation functions
   * Sleep
   * Get UNIX epoch time
   * Parse/format dates
 * Strings
   * Is letter
   * Is digit
   * Escape/unescape
 * Math
   * Random
   * Trig
   * Sqrt
 * Arrays
   * Shuffle
   * Sum
 * Add `break` built-in for leaving loops.
 * Add `and` and `or` commands for conditions.
 * Support up arrow in REPL
