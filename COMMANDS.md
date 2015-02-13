# Overview

A standard pragmash environment has a standard library of functions. These functions include file manipulation commands, mathematical operations, and a host of other things. This file lists each command and explains it in a decent amount of detail.

# Operators

Unlike most other programming languages, pragmash does not include built-in operators. Instead, commands serve to provide basic operations.  Here's a list of commands and brief examples for each one.

### The + and * operators

These operators both take 0 or more arguments and perform arithmetic on either big integers on floating points. Here are some examples:

 * `* 3 5` yields "15"
 * `* 1 2 3 4` yields "24"
 * `+ 1 2 3` yields "6"
 * `* 1.5 2.5` yields "3.7500000000"

### The / operator

This operator takes two arguments, a and b, and returns a/b. Both arguments are parsed as floating points. For example, `/ 3 2` yields "1.5000000000".

### The - operator

This operator takes two arguments, a and b, and returns a-b. The arguments can be big integers or floating points.

### The [] operator

The subscript operator (denoted with two brackets `[]`) is used to access an element in a newline-delimited list. The first argument is the list, the second is the index. For example:

 * `[] "hey\nthere" 0` yields "hey"
 * `[] "hey\nthere" 1` yields "there"
 * `[] "hey\nthere" 2` throws an exception.

### The <= operator

This operator checks if the first numerical argument is less than or equal to the second. It returns "true" in such a case, and "" otherwise.

### The >= operator

This operator checks if the first numerical argument is greater than or equal to the second. It returns "true" in such a case, and "" otherwise.

### The < operator

This operator checks if the first numerical argument is less than the second. It returns "true" in such a case, and "" otherwise.

### The > operator

This operator checks if the first numerical argument is greater than the second. It returns "true" in such a case, and "" otherwise.

### The = operator

This returns "true" if all its arguments are equal (when compared as strings). Returns "" otherwise.

# I/O

## The console

### The "gets" command

This command reads a line from the console and takes no arguments. The newline is not included in the returned string.

### The "print" command

This command prints all of its arguments (separated by spaces) to the console. It does not print a newline, but it does flush the output.

### The "puts" command

This command prints all of its arguments (separated by spaces), followed by a newline, to the console.

## The web and filesystem

### The "read" command

This command takes one argument which is either a filepath or a URL. It returns a string representing the contents of the specified resource, or throws an exception.

### The "write" command

This command takes two arguments: first a path; second, some data to write to the path. It throws an exception if the data cannot be written.

# Basic functionality

### The "count" command

This command takes one argument and returns how many elements it contains as a newline-delimited list.

### The "exit" command

This command exits the program. It takes an optional integer argument with a return value. If this argument is specified but is not a valid integer, the command throws an exception.

### The "get" command

This command takes one argument--a variable name--and returns its contents. It throws an exception if the variable is not defined.

### The "len" command

This command takes one string argument and returns its length, in bytes.

### The "set" command

This command takes two arguments and sets a variable. The first argument is a variable name, the second is a value to give the variable.

### The "throw" command

This command raises an exception with the specified error message. It joins its arguments with spaces and uses them for the error message.

# Strings

### The "echo" command

This command joins its arguments and inserts spaces between them.

### The "join" command

This command joins its arguments without inserting spaces between them.

### The "match" command

This command takes a regular expression and a string. It returns an array of matches. Each sub-match is its own element in the array.

For example, `match "x([a-z])z" "abc xyz xwz xoz"` yields the array equivalent to `arr xyz y xwz w xoz o`.

# Arrays

### The "arr" command

This command joins its arguments with newlines. For example, `arr a b c` generates "a\nb\nc".

### The "range" command

This command generates a newline-delimited list of integers.

If the command is given one argument `N`, it will generate the ordered list of integers `i` such that `0 <= i < N`.

If the command is given two arguments `M` and `N`, it will generate the ordered list of integers `i` such that `M <= i < N`

If the command is given three arguments, it generates the ordered list of integers starting with the first argument going to the second argument, stepping by the third argument each time. For example, `range 10 5 -2` yields `10\n8\n6`.

# Filesystem

### The "glob" command

This command takes any number of arguments and "globs" files by those names. 

For example, if my current directory includes the files "foo" "bar" and "foobar", `glob foo*`, it would return the array equivalent to `arr foo foobar`.
