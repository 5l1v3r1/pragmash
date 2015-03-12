# Overview

A standard pragmash environment should have a basic set of commands. These commands offer everything from network I/O to string manipulation. This file lists these commands.

These commands are divided into several categories. Here are the categories:

 * [Operators](#api-operators)
 * [I/O](#api-io)
 * [Language essentials](#api-language)
 * [Strings](#api-strings)
 * [Arrays](#api-arrays)
 * [Filesystem](#api-filesystem)
 * [Math](#api-math)
 * [Time](#api-time)

<a name="api-operators"></a>
# Operators

Unlike most other programming languages, pragmash does not include built-in operators. In their place, it provides symbolically-named commands.  Here's the list of these commands.

### The + command

This takes 0 or more arguments and returns their sum. If no numbers are provided, "0" is returned. The numbers can be floating-points or big integers.

Examples:

 * `+` yields "0"
 * `+ 1 2 3` yields "6"
 * `+ 3 -2 1` yields "2"
 * `+ 1.5 2.5` yields "4"

### The * command

This takes 0 or more arguments and returns their product. If no numbers are provided, "1" is returned. The numbers can be floating-points or big integers.

 * `* 3 5` yields "15"
 * `* 1 2 3 4` yields "24"
 * `* 1.5 2.5` yields "3.75"

### The / command

This takes exactly two arguments and returns the first divided by the second. If both numbers are big integers, this may return a big integer if the quotient is a whole number. If the denominator is 0, this throws an exception.

Examples:

 * `/ 2 3` yields "0.6666666666666666"
 * `/ 9 3` yields "3"
 * `/ 1 0` throws an exception
 * `/ 1000000000000000000000000000000000000000000000010 10` yields "100000000000000000000000000000000000000000000001"

### The - command

This takes exactly two arguments and returns the first minus the second.

Examples:

 * `- 3 2` yields "1"
 * `- 2 3` yields "-1"
 * `- 2 1.5` yields "0.5"

### The % command

This takes two arguments and returns the first modulo the second. If either argument is not an integer, this computes `a - b*floor(a/b)` where *a* is the first argument and *b* is the second.

Examples:

 * `% 3 2` yields "1"
 * `% -10 3` yields "2"
 * `% 30.5 10` yields "0.5"

### The [] command

This is used to access an element in a list which is delimited by newlines. The first argument is the list, the second is the index. For example:

 * `[] "hey\nthere" 0` yields "hey"
 * `[] "hey\nthere" 1` yields "there"
 * `[] "hey\nthere" 2` throws an exception.

### The &lt;= operator

This takes two numerical arguments and checks if the first is less than or equal to the second. It returns "true" in such a case, and "" otherwise.

### The &gt;= operator

This takes two numerical arguments and checks if the first is greater than or equal to the second. It returns "true" in such a case, and "" otherwise.

### The &lt; operator

This takes two numerical arguments and checks if the first is less than the second. It returns "true" in such a case, and "" otherwise.

### The &gt; operator

This takes two numerical arguments and checks if the first is greater than the second. It returns "true" in such a case, and "" otherwise.

### The = operator

This takes zero or more arguments and returns "true" if and only if all its arguments are equal when compared as strings. Otherwise, this returns "".

### The &amp;&amp; operator

This takes zero or more arguments and returns "true" if none of the arguments are empty. Otherwise, this returns "".

### The || operator

This takes zero or more arguments and returns its first non-empty argument. If all arguments are empty or no arguments were supplied, this returns "".

<a name="api-io"></a>
# I/O

## The console

### gets

This reads a line from the console and returns it. A newline character is not included in the resulting string.

### print \[string...\]

This prints all of its arguments to the console separated by spaces. It does not print a newline, but it does flush the output.

### puts \[string...\]

This prints all of its arguments to the console separated by spaces. It follows this output with a newline character.

## The web and filesystem

### httpCookiesOff

This disables cookie saving for httpGet httpPost. This will delete all existing cookies.

### httpCookiesOn

This enables cookie saving for httpGet and httpPost. This will delete all existing cookies.

### httpGet &lt;url&gt;

This runs an HTTP get request. This uses cookies if cookies are enabled.

### httpPost &lt;url&gt; &lt;content-type&gt; &lt;body&gt;

This runs an HTTP post request. This uses cookies if cookies are enabled.

### read &lt;resource&gt;

This takes one argument which is either a file path or a URL. It returns a string representing the contents of the specified resource, or throws an exception if the resource cannot be read.

This does not use or save cookies if the argument is a URL.

### write &lt;path&gt; &lt;data&gt;

This writes a string to a file. It throws an exception if the data cannot be written.

## OS commands

### cmd &lt;name&gt; \[arguments...\]

This executes a command on the system. On UNIX-based systems, this is similar to running a command in a shell. It returns the combined output (stdout+stderr) of the command. This throws an exception if the command cannot be executed or if it fails in some platform-specific way.

<a name="api-language"></a>
# Language essentials

### call &lt;name&gt; \[arrays...\]

This takes a command name and zero or more arrays to use as arguments. It executes the command with the specified arguments.

Examples:

 * `call + 1\n2\n3` yields "6".
 * `call echo (arr a b c)` yields "a b c".
 * `call call echo (arr a\nb\nc d\ne\nf)` yields "a b c d e f"

### eval &lt;code&gt;

This executes a block of pragmash code. The code which is executed will have complete access to the main script's variables. It will be able to throw exceptions. It will be able to print to the console. In essence, the code runs as if it were part of the main script. The only difference is that the code may use the "return" keyword to return values.

Examples:

 * `eval "return test"` yields "test"
 * `eval "print test"` prints "test" to the console

### exec &lt;file&gt;

This executes a pragmash file. `exec <file>` is almost exactly equivalent to `eval (read <file>)`. The only difference is that exceptions generated from the exec'd script include the script's filename.

### exit \[exit code\]

This exits the program. If the exit code is specified, it will be used as the numerical return value of the pragmash executable. If the exit code is not a valid number, an exit code of 1 is used.

### get &lt;variable&gt;

This returns the contents of a variable. It throws an exception if the variable is not defined.

### pragmash &lt;path&gt; \[arguments...\]

This executes a pragmash script in a new context and returns its return value. The script runs with a new set of variables (including the built-in ones), but it may still print to the console or exit the parent script. The optional arguments after the script path determine the child script's ARGV variable. The child script's DIR and SCRIPT variables will be based on the path of the child script.

For example, suppose this is the contents of a file "main.pragmash":

    pragmash foo.pragmash arg1 arg2
    puts unreachable

and this is the contents of the file "foo.pragmash":

    puts ([] $ARGV 1)
    exit 1

This would print "arg2" to the console and exit with status code 1. The string "unreachable" would not be printed to the screen.

### set &lt;variable&gt; &lt;value&gt;

This assigns a value to a given variable.

### throw \[string...\]

This throws an exception. It joins its arguments with spaces and uses the result as the error message.

<a name="api-strings"></a>
# Strings

### chars &lt;string&gt;

This generates a newline-delimited list of strings which correspond to each character of the argument. Newline characters are encoded as the two-character "\\\\n" escape sequence. For example, `chars 12\n3` yields `"1\n2\n\\n\n3"`.

### echo \[string...\]

This joins its arguments with spaces and returns the result.

Examples:

 * `echo` yields ""
 * `echo arg` yields "arg"
 * `echo arg ument` yields "arg ument"
 * `echo a  b    "c  d"` yields "a b c  d"

### escape \[string...\]

This replaces backslashes with double backslashes and newlines with "\\\\n". This makes it easier to represent array elements which contain newlines.

### join \[string...\]

This joins its arguments without inserting spaces between them.

### len &lt;string&gt;

This returns the length of a string in bytes.

### lowercase \[string...\]

This joins its arguments with spaces and converts the result to lower-case.

### match &lt;regexp&gt; &lt;haystack&gt;

This matches a string against a regular expression. It returns an array of matches. Each sub-match is its own element in the array.

For example, `match "x([a-z])z" "abc xyz xwz xoz"` yields the array equivalent to `arr xyz y xwz w xoz o`.

### rep &lt;haystack&gt; &lt;needle&gt; &lt;replacement&gt;

This performs a global find-and-replace operation. It replaces all occurances of a "needle" inside a "haystack" with a "replacement" string.

For example, `rep abcdcba a A` yields "AbcdcbA".

### repreg &lt;haystack&gt; &lt;regexp&gt; &lt;replacement&gt;

This performs a global find-and-replace operation with regular expressions. Inside the replacement string, `$1` can be used to refer to the first submatch, `$2` to the second, etc.

Examples: "X e X tX e Xe bro"

 * `repreg "Alex Nichol" [A-Z] _` yields "\_lex \_ichol"
 * `repreg "Alex Nichol" "([A-Z])([a-z])" "$2$1"` yields "lAex iNchol"
 * `repreg "10.50 20 30" "([0-9\\.]*)" "$$$1"` yields "$10.50 $20 $30"

### substr &lt;string&gt; &lt;start&gt; &lt;end&gt;

This takes three arguments and performs bytewise substring. The first is a string, the second is the starting index, and the third is the ending index.

For example, `substr yoyo 1 3` yields "oy".

### unescape &lt;string&gt;

This inverts the effect of the escape command.

### uppercase \[string...\]

This joins its arguments with spaces and converts the result to upper-case.

<a name="api-arrays"></a>
# Arrays

### arr \[arrays...\]

This joins its arguments with newlines and throws away empty arguments.

Examples:

 * `arr a b c` yields "a\nb\nc"
 * `arr "" a ""` yields "a"

### contains &lt;array&gt; &lt;element&gt;

This takes an array and a string and returns "true" if the array contains the string. Otherwise, it returns "".

### count &lt;array&gt;

This takes a newline-delimited list and returns the number of elements it contains. If the argument is "", this returns 0.

### delete &lt;array&gt; &lt;index&gt;

This deletes an element from an array.

Examples:

 * `delete (arr a b c) 1` yields "a\nc"
 * `delete (arr a b c) 0` yields "b\nc"
 * `delete "" 0` throws an exception

### insert &lt;array&gt; &lt;index&gt; &lt;element&gt;

This inserts an element into an array.

Examples:

 * `insert (arr a b c) 1 A` yields "a\nA\nb\nc"
 * `insert (arr a b c) 3 d` yields "a\nb\nc\nd"
 * `insert (arr a b c) 4 d` throws an exception

### range \[start\] &lt;end&gt; \[count\]

This generates a newline-delimited list of integers.

If the command is given one argument `N`, it will generate the ordered list of integers `i` such that `0 <= i < N`.

If the command is given two arguments `M` and `N`, it will generate the ordered list of integers `i` such that `M <= i < N`

If the command is given three arguments, it generates the ordered list of integers starting with the first argument going to the second argument, stepping by the third argument each time. For example, `range 10 5 -2` yields `10\n8\n6`.

### shuffle &lt;array&gt;

This takes an array and returns an array with the same elements in a random order.

### sort &lt;array&gt;

This takes an array and returns an alphabetically sorted version.

### sortnums &lt;array&gt;

This takes an array of numbers and returns the sorted array.

### subarr &lt;array&gt; &lt;start&gt; &lt;end&gt;

This takes an array and two indices. It returns a portion of the original array.

Examples:

 * `subarr (arr a b c) 1 3` yields "b\nc"
 * `subarr (arr a b c) 0 2` yields "a\nb"

### sum \[arrays...\]

This takes zero or more arrays of numbers and returns the sum of all the numbers.

<a name="api-filesystem"></a>
# Filesystem

### exists &lt;path&gt;

This returns "true" if a file exists or "" if it does not. It may throw an exception if the file is inaccessible.

### filetype &lt;path&gt;

This returns the type of a named file. This can be "file", "dir", "link", or "other". This will throw an exception if the file's type cannot be determined or if it does not exist.

### glob \[globs...\]

This takes any number of arguments and "globs" files by those names. 

For example, if my current directory includes the files "foo" "bar" and "foobar", `glob foo*`, it would return the array equivalent to `arr foo foobar`.

### mkdir &lt;dir&gt;

This creates a directory at a given path. This will not create intermediate directories. This may return an exception if the directory cannot be created.

### path \[comps...\]

This takes any number of string arguments and joins them as path components.

### rm &lt;path&gt;

This deletes a file or an empty directory. This will throw an exception if the file or directory cannot be deleted.

### rmall &lt;path&gt;

This deletes a file or directory recursively. This will throw an exception if the file or directory cannot be deleted.

### touch \[path...\]

This creates one or more files or updates their timestamps to the current time. This will throw an exception if the file's timestamp cannot be changed or if the file cannot be created.

<a name="api-math"></a>
# Math

### abs &lt;number&gt;

This takes the absolute value of its numerical argument.

Examples:

 * `abs -2` yields "2"
 * `abs 2` yields "2"

### acos &lt;number&gt;

This returns the inverse cosine of a value. This will throw an exception if the number is less than -1 or greater than 1.

### asin &lt;number&gt;

This returns the inverse sine of a value. This will throw an exception if the number is less than -1 or greater than 1.

### atan &lt;number&gt;

This returns the inverse tangent of a value.

### atan2 &lt;y&gt; &lt;x&gt;

This returns the inverse tangent using an x and y coordinate.

### ceil &lt;number&gt;

This returns the greatest integer which is less than or equal to a given floating-point number.

### cos &lt;angle&gt;

This computes the cosine of an angle in radians.

### exp \[exponent\]

This takes an argument x and computes e^x where e is Euler's constant. If no arguments are given, this returns the value of Euler's constant.

### factorial &lt;number&gt;

This takes a number and returns its factorial. If the number is not a positive integer, this uses the [Gamma function](http://en.wikipedia.org/wiki/Gamma_function) to compute a fractional answer.

### floor &lt;number&gt;

This returns the lowest integer which is greater than or equal to a given floating-point number.

### log \[base\] &lt;number&gt;

This computes a logarithm. If you supply one argument, this computes log base 10 of its argument. If there are two arguments, the first argument is treated as the base. If either argument is invalid, this will throw an exception.

### pi

This takes no arguments and returns the value of pi.

### rand

This returns a random floating-point number between 0.0 and 1.0.

### round &lt;number&gt;

This rounds a floating-point number to the nearest integer.

### sin &lt;angle&gt;

This computes the sine of an angle in radians.

<a name="api-time"></a>
# Time

### sleep &lt;seconds&gt;

This takes a numerical argument and sleeps for that many seconds. The argument may be a floating-point number.

### time

This returns the current UNIX epoch time as a floating-point in seconds.
