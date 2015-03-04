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

### The % operator

This is the modulus operator. It takes two arguments, a and b, and finds a mod b. If either argument is a floating point, this computes a - b*floor(a/b).

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

### The && operator

If none of the arguments are empty, this returns "true". Otherwise, it returns "".

### The || operator

This returns its first non-empty argument. If all arguments are empty, this returns "".

# I/O

## The console

### gets

This command reads a line from the console and takes no arguments. The newline is not included in the returned string.

### print

This command prints all of its arguments (separated by spaces) to the console. It does not print a newline, but it does flush the output.

### puts

This command prints all of its arguments (separated by spaces), followed by a newline, to the console.

## The web and filesystem

### read

This command takes one argument which is either a filepath or a URL. It returns a string representing the contents of the specified resource, or throws an exception.

### write

This command takes two arguments: first a path; second, some data to write to the path. It throws an exception if the data cannot be written.

## OS commands

### cmd

This command takes 1 or more argument and executes it as a command. It returns the combined output (stdout+stderr) of the command, or throws an error.

# Language functionality

### call

This command takes a command name and zero or more arrays to use as arguments. It executes the command with the specified arguments. For example, `call + 1\n2\n3` would yield 6.

### count

This command takes one argument and returns how many elements it contains as a newline-delimited list.

### eval

This command executes a block of pragmash code. The code which is executed will have complete access to the main script's variables. It will be able to throw exceptions. It will be able to print to the console. In essence, the code runs as if it were part of the main script. The code may use the "return" keyword to return values.

Example: `eval "puts hey"`

### exec

This command executes a pragmash file. The scirpt which is executed will have complete access to the main script's variables. It will be able to throw exceptions. It will be able to print to the console. In essence, the script runs as if it were part of the main script. The executed file can use the "return" keyword to return values.

Example: `exec (join $DIR /some_file.pragmash)`

### exit

This command exits the program. It takes an optional integer argument with a return value. If this argument is specified but is not a valid integer, the command throws an exception.

### get

This command takes one argument--a variable name--and returns its contents. It throws an exception if the variable is not defined.

### len

This command takes one string argument and returns its length, in bytes.

### pragmash

This command takes one or more arguments. This executes the script (specified by the first argument) in a new context and returns its return value. The script runs with a new set of variables (including $ARGV and $DIR), but it may still print to the console or exit the program.

For example, suppose this is the contents of a file "main.pragmash":

    pragmash foo.pragmash hey there
    puts unreachable

and this is the contents of the file "foo.pragmash":

    puts ([] $ARGV 1)
    exit 1

This would print "there" to the console and exit with status code 1. The "unreachable" print would not run.

### set

This command takes two arguments and sets a variable. The first argument is a variable name, the second is a value to give the variable.

### throw

This command raises an exception with the specified error message. It joins its arguments with spaces and uses them for the error message.

# Strings

### chars

This generates a newline-delimited list of strings which correspond to each character of the string. Newline characters are encoded as the two-character "\n" escape sequence. For example, `chars 12\n3` yields `"1\n2\n\\n\n3"`.

### echo

This command joins its arguments and inserts spaces between them.

### escape

This command replaces backslashes with double backslashes and newlines with "\\n". This is helpful for storing strings with newlines as elements in arrays.

### join

This command joins its arguments without inserting spaces between them.

### lowercase

This joins its arguments with spaces and converts the result to lower-case.

### match

This command takes a regular expression and a string. It returns an array of matches. Each sub-match is its own element in the array.

For example, `match "x([a-z])z" "abc xyz xwz xoz"` yields the array equivalent to `arr xyz y xwz w xoz o`.

### rep

This command takes three arguments. It replaces all occurances of the second argument with the third argument in the first argument.

For example, `rep heythere e E` yields "hEythErE".

### repreg

This command takes a regular expression, a string, and a replacement string. It replaces all occurances of the regular expression. Inside the replacement, `$1` can be used to refer to the first submatch, `$2` to the second, etc.

For example, the following code yields "X e X tX e Xe bro"

    repreg "[a-z](e)[a-z]" "hey there bro" "X $1 X"

### substr

This command takes three arguments and performs bytewise substring. The first is a string, the second is the starting index, and the third is the ending index.

For example, `substr yoyo 1 3` yields "oy".

### unescape

This command inverts the effect of the escape command.

### uppercase

This joins its arguments with spaces and converts the result to upper-case.

# Arrays

### arr

This command joins its arguments with newlines and throws away empty arguments. For example, `arr a b c` generates "a\nb\nc". As another example, `arr "" a ""` generates "a".

### contains

This command takes an array and a string and returns "true" if the array contains the string.

### delete

This command deletes an element from the array.

For example, `delete (arr a b c) 1` yields the array equivalent of `(arr a c)`

### insert

This command inserts an element into the array.

For example, `insert (arr a b c) 1 A` yields the array equivalent of `(arr a A b c)`

### range

This command generates a newline-delimited list of integers.

If the command is given one argument `N`, it will generate the ordered list of integers `i` such that `0 <= i < N`.

If the command is given two arguments `M` and `N`, it will generate the ordered list of integers `i` such that `M <= i < N`

If the command is given three arguments, it generates the ordered list of integers starting with the first argument going to the second argument, stepping by the third argument each time. For example, `range 10 5 -2` yields `10\n8\n6`.

### shuffle

This command takes an array and returns an array. The resulting array will be permuted in a random order.

### sort

This command takes an array and returns an array. The resulting array will be sorted alphabetically.

### sortnums

This command takes an array of numbers and returns the sorted array.

### subarr

This command takes a string and two indices. It returns a slice of an array.

For example, `subarr (arr a b c) 1 3` returns the array equivalent to `arr b c`.

### sum

This command takes zero or more arrays of numbers and returns the sum of all the numbers.

# Filesystem

### exists

This returns a boolean expression indicating whether or not a file exists. It may throw an error if it does not have permissions to check or if some other error occurs.

### filetype

This returns the type of a named file. This can be "file", "dir", "link", or "other".

### glob

This command takes any number of arguments and "globs" files by those names. 

For example, if my current directory includes the files "foo" "bar" and "foobar", `glob foo*`, it would return the array equivalent to `arr foo foobar`.

### mkdir

This command creates a directory at a given path. This will not create intermediate directories.

### path

This command takens any number of string arguments and joins them as path components.

### rm

This deletes a file or an empty directory.

### rmall

This deletes a file or directory recursively.

# Math

### abs

This takes the absolute value of its numerical argument. For example, `abs -2` yields `2`.

### ceil

This returns the greatest integer which is less than or equal to a floating-point number.

### cos

This computes the cosine of an angle in radians.

### exp

This takes an argument x and computes e^x. If no arguments are given, this returns the value of e.

### factorial

This takes a number and returns its factorial. If the number is not a positive integer, this uses the Gamma function to compute a fractional answer.

### floor

This returns the lowest integer which is greater than or equal to a floating-point number.

### log

This computes a logarithm. If you supply one argument, this computes log base 10 of the argument. If there are two arguments, the first argument is treated as the base. If either argument is invalid, this will throw an exception.

### pi

This takes no arguments and returns the value of pi.

### rand

This returns a random floating point between 0.0 and 1.0.

### round

This rounds a floating point number to the nearest integer.

### sin

This computes the sine of an angle in radians.

# Time

### sleep

This takes a numerical argument and sleeps for that many seconds.

### time

This returns the current UNIX epoch time as a floating point in seconds.
