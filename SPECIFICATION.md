# Abstract

This document describes every feature of pragmash. It uses human-readable langugae and tries to be concise.

# Charset

A pragmash script is encoded as UTF-8. The term "whitespace" refers to all the characters which are considered whitespace in UTF-8. A newline is the character represented by the number `10` \(i.e. "\n"\).

# Lines

There are two types of "lines" in a pragmash script. **Physical lines** are separated by newline characters and include all of the leading and trailing whitespace. The only thing to note is that carriage returns at the ends of physical lines are ignored. A **logical line** is similar to a physical line, but leading and trailing whitespace is removed and it is possible to represent a logical line across multiple physical lines using backslashes at the end of each physical line. For example, this counts as one logical line but two physical lines:

    this is \
    one line of code

The newline after the backslash is not counted as part of the logical line.

# Comments

If a logical line begins with a `#`, it is considered a comment. Comments are not evaluated and do not affect the behavior of the code. Note that logical lines have no trailing or leading whitespace, so the `#` can be preceded by whitespace on its physical line.

Since logical lines can be continued across multiple physical lines, this code is entirely ignored:

    # this is a long\
    comment across multiple\
    physical lines.

However, it is recommended that you use separate logical lines and precede each one by `#`:

    # this is a long
    # comment across multiple
    # physical lines.

Using one physical line for every line of a comment makes it easier to tell that each line is commented.

# Blocks

The `{` and `}` tokens \(also known as curly braces\) are used in pragmash to define blocks of code. Note, however, that a pragmash script may contain curly braces inside of strings, command names, or comments. In these cases, curly braces are treated as regular characters. Curly braces only denote blocks of code when they are situated correctly \(i.e. when they appear on lines in the proper places after the proper keywords\). In these cases, every `{` must be the last character on a logical line and every `}` must be the first character on a logical line. Blocks of code can be nested, such as in this case:

    if a {
      while b {
        # something here
      }
    }

If text comes after a `}` on a logical line, there must be at least one whitespace character before the text after the `}`. On a line which begins a block, the trailing `{` must be preceded by at least one whitespace character.

Here are some examples of **invalid** block syntax:

    # Invalid because no space before {
    if (get a){
    }

    # Invalid because no space after }
    if $a {
    }else {
    }

    # Invalid because { is on a different logical line than the if token.
    if $a
    {
        puts hey
    }

# Strings

Every argument to every command is a string. Variables contain strings and their names are themselves strings. Command names, too, are strings and every command returns a string.

But what exactly is a string? In pragmash, a string is an array of bytes. While strings can often be treated as a piece of UTF-8 text, they are not limited to that. Strings can also store binary data. The contents of a string does not have to be valid text in any character encoding.

A string can be expressed without quotes as long as it has no unescaped whitespace or special characters. It is also possible to surround strings with `'` \(single quotes\) or `"` \(double quotes\). In these cases, unescaped whitespace \(and some special characters\) can be present.

An escaped character within a string begins with a backslash. To escape a space, use a backspace before it such as in `escaped\ whitespace`. To escape a backslash, use `\\`. Inside of single quotes, it is necessary to escape single quotes as `\'`. Likewise, you should use `\"` if you want to include quotes inside of a quoted string.

Escapes can also be used to write out characters which a programmer cannot easily enter from their keyboard. These are the available escape sequences, not including the ones discussed above:

| Sequence    | Description         | Hex Code \(if applicable\) |
|-------------|---------------------|----------------------------|
| \\(         | Parenthesis         | 0x28                       |
| \\)         | Parenthesis         | 0x29                       |
| \\?         | Question mark       | 0x3f                       |
| \\'         | Single quote        | 0x27                       |
| \\"         | Double quote        | 0x22                       |
| \\a         | Bell                | 0x07                       |
| \\b         | Backspace           | 0x08                       |
| \\f         | New page            | 0x0c                       |
| \\n         | New line            | 0x0a                       |
| \\r         | Carriage return     | 0x0d                       |
| \\t         | Tab                 | 0x09                       |
| \\v         | Vertical tab        | 0x0b                       |
| \\nnn       | 3-digit octal value |                            |
| \\xnn       | 2-digit hex value   |                            |
| \\unnnn     | Unicode value       |                            |
| \\Unnnnnnnn | Unicode value       |                            |

Here are some examples of strings:

 * `someText`
 * `some\ text\ with\ spaces`
 * `"some text with spaces"`
 * `'some text with spaces'`
 * `"some text which contains \"quotes,\" so to speak"`

# Commands

Many logical lines in a typical pragmash program will be commands. A command starts with a command name and is followed by zero or more arguments. Both the command name and the arguments are either strings or nested commands, as explained later in this section. Here are some examples of commands which are comprised entirely of string arguments:

 * `puts hello world!`
 * `write file.txt "some text"`
 * `rm file\ with\ space\ in\ name.txt`

Commands can be nested so that the output of one command can be used as an input to another. Instead of quotes, a nested command is wrapped in `(` and `)` \(parentheses\). Nested commands are evaluated first so that they can be passed as string arguments to an outer command. Here are some examples of nested commands:

 * `write file.txt (read 'other file.txt')`
 * `write file.txt (+ 1 (rand))`
 * `set x (+ $y 2)`

Note that a bare string cannot contain a `)`. If it does, the parser will mistake it for the end of the code block. There are two distinct ways of dealing with this:

    # Bad! This code is not valid.
    puts (read filename_ending_in_parenthesis))

    # Better, but still ugly.
    puts (read filename_ending_in_parenthesis\))

    # Good!
    puts (read 'filename_ending_in_parenthesis)')

Note that arguments to commands are separated by whitespace. Arguments must be separated by whitespace; the code is invalid otherwise. Here are examples of **invalid code**:

    ############################################
    # NOTE: this code is all invalid.          #
    # Please don't just read the code samples  #
    # or else you'll look at invalid code. :-P #
    ############################################
    puts "string 1""string 2"
    puts "string 1"'string 2'
    puts (read file.txt)string
    puts string(read file.txt)
    puts 'string'(read file.txt)
    puts 'string'string
    puts /Path/to/"File with spaces"/foobar

Note however that this code is valid and would print `string'string'`:

    puts string'string'

This means that a bareword can contain double or single quotes as long as it does not begin with one.

# Concise variable syntax

Whenever a bare string begins with a $, it is not treated as a bare string. Instead, the interpreter treats something like `$xyz` as `(get xyz)` \(that is, as a nested command\). In a standard pragmash environment, `get` accesses a variable. Thus, the $ makes it easier to access variables.

# Concise mutation syntax

It is often helpful to modify the contents of an existing variable. For example, you may wish to add 3 to a counter, like so:

    set a (+ $a 3)

This code seems very verbose for such a simple operation. Luckily, you can also do this:

    # This expands to set a (+ $a 3)
    += a 3

Whenever you execute a command whose name ends with "=", the command is rewritten according to a template. The call `command= y args...` is rewritten as `set y (command $y args...)`. This is not just syntactic sugar. This transformation occurs at runtime, so it is possible to use nested commands for the variable name and command name:

    # This will be equivalent to doing set variableName (+ $variableName 3)
    (echo +=) (echo variableName) 3

# Conditions

A condition looks like a set of arguments but is evaluated to a boolean expression. Conditions are used in `if` blocks and `while` loops. In addition, some commands may process their arguments as if the arguments were a condition.

There are three types of conditions.

 * An empty condition \(one with zero arguments\) is always *true*.
 * A condition with one argument is *true* if and only if the argument is not an empty string. If the argument is empty, the condition is *false*.
 * A condition with mulitple arguments is *true* if and only if all of the arguments are equal. Otherwise, it is false.

# Control blocks: an overview

In pragmash, certain keywords at the beginning of logical lines indicate that the line is not a command but rather the beginning of a control block. These keywords are `if`, `for`, `while`, `try` and `def`.

# *if* blocks

An `if` block makes it possible to run different pieces of code depending on the outcome of a condition. An if block can have multiple branches, denoted by `if`, `else if`, and `else` depending on their context. Here are some example `if` blocks:

    if $name Alex {
      puts Hello.
    }

    if $name Joe {
      puts Hi, joe old pal.
    } else {
      puts Who are you and what have you done with joe?
    }

    if $count 3 {
      puts It's three.
    } else if $count 4 {
      puts It's four.
    } else {
      puts I don't know exactly what it is.
    }

An `else` branch can only appear as the last branch of the if statement. An `else if` branch can never appear as the first branch of an if statement.

The not keyword can go directly after `if` or `else if` to negate the result of the expression before evaluating it. For instance, you can do this:

    if not $name Alex {
        puts You are not Alex
    } else if not $age 18 {
        puts If you were really alex, you would have known your own age.
    }

Since not is a special keyword, it musn't be in quotes and musn't be written in terms of escape sequences. This makes it possible for you to use the string "not" in a condition within an if statement:

    if "not" $token {
        puts 'The token was "not"'
    }

As long as the word "not" appears in quotes, it will not be mistaken for the not keyword.

# *while* blocks

A `while` block executes a piece of code again and again until a condition is false. Here are some examples of while blocks:

    while {
      puts This is an infinite loop!
    }

    set a 0
    while (< $a 10) {
      puts This is a loop and the counter is is currently $a
      ++ a
    }

The `not` keyword can be placed directly after the `while` keyword to negate the result of the condition. This works very similarly to if statements:

    set a 0
    while not (>= 10 $a) {
      puts This is a loop and the counter is is currently $a
      ++ a
    }

The `break` command can exit a loop. If you pass `break` an optional integer argument, you can specify how many levels of loops to break out of. By default, this is 1. Here are some examples:

    # This loop will output "Hello world." and then end.
    while {
      puts Hello world.
      break
    }

    # This loop will never output "Foo" or "Bar".
    while {
      while {
        while {
          puts Hey there.
          break 2
          puts Foo
        }
        puts Bar
      }
    }

The `continue` command skips to the next iteration of a loop. It takes an optional integer argument which is 1 by default. If the argument is more than one, it specifies the index of the outer loop to continue. For instance, `continue 2` breaks the innermost loop and continues the loop which contains it. Here are some examples:

    # This will print the odd numbers between 0 and 100.
    set a 0
    while (< $a 100) {
      ++ a
      if 0 (% a 2) {
        continue
      }
      puts $a
    }

    # This will print "hello" in a loop and never print "foo" or "bar".
    while {
      while {
        puts hello
        continue 2
        puts foo
      }
      puts bar
    }

# Errors \(a.k.a. exceptions\)

In pragmash, it is difficult if not impossible to represent errors as return values. To show why this is so, suppose I decided that a return value of "ERROR" indicated that a command failed. Then, if you tried to read a file with the contents "ERROR" you would mistakenly think that there was an error reading the file. Other languages solve this problem using multiple return values, null values, or different datatypes. However, these approaches can result in clutter and redundant code.

Instead of using return values for errors, pragmash has a built-in mechanism for handling errors. When an error occurs, the error is propagated to the nearest `try` block \(as described in the next section\) or causes the program to terminate. Errors which arise this way include both a context \(i.e. a line number\) and an error message.

The `throw` command allows you to throw an error from your code. The arguments are joined with spaces and together form the error message. Here are some examples of calls to `throw`:

 * `throw file not found: $path`
 * `throw unknown error`
 * `throw "unknown error"`

Generally, error messages should not be capitalized and should not end with punctuation \(i.e. periods\). This is because an error message will usually be used within a broader sentence \(for example, "Error at line 12: unknown error." wherein the actual message was "unknown error"\).

# *try* blocks

The `try` block makes it possible to catch errors and handle them gracefully. These blocks have two distinct parts. The first part of a `try` block is the code which will run normally. If an error occurs within the first part of the try block, the "catch" part of the `try` block will be executed. Within the catch block, you can optionally access the error message and error context via variables whose names you specify.

Here are some examples of `try` blocks:

    # Catch an error and do nothing in the handler.
    # This will print "Foo" but not "Bar"
    try {
      puts Foo
      throw An error.
      puts Bar
    } catch {
    }

    # This will print "Foo" and "Bar" but not "Baz".
    try {
      puts Foo
      throw Bar
      puts Baz
    } catch e {
      puts $e
    }

    # This is a real life example. You can capture both the context and the
    # error in a try block.
    try {
      set a (read file1.txt)
      set b (read file2.txt)
      write files1and2.txt (join $a $b)
    } catch c e {
      puts Error joining files: $e. Error was at: $c
    }

# Lists

Strings are the only real datatype in pragmash. However, in certain contexts, a string can be treated as a newline-separated list. In these contexts, the string `"1\n2\n3\n4"` would represent the array `[1, 2, 3, 4]`. One unusual thing about newline-separated lists is that the empty string corresponds to the empty array, whereas one might expect it to correspond to an array with one empty string element.

In many situations, newline-separated lists are enough. In some rare cases, however, it might be necessary to store arrays of strings which contain newlines themselves. In these cases, it is necessary to escape the elements of the list.

# *for* blocks

The `for` block runs a block of code for each element in a newline-separated list. A for block can take one, two, or three arguments depending on how much information the block needs for each iteration.

If a for block has one argument, that argument is treated as an array. The actual elements of the array are not relevant; all that matters is the length of the array:

    for "1\n2\n3\n4" {
      puts This will print out 4 times.
    }

If a for block has two arguments, the first argument is treated as a variable name and is re-assigned to the value of each element for each iteration:

    puts Here are some fruits:
    for fruitName "Pickels\nApples\nPeaches" {
      puts $fruitName
    }

If a for block has three arguments, the first and second arguments are treated as variable names. The first argument is assigned the index for each iteration \(starting at zero\). The second argument is assigned the value of each element for each iteration:

    for i fruitName "Pickels\nApples\nPeaches" {
      puts Fruit number $i is $fruitName
    }

The `break` and `continue` keywords work the same for for blocks as they do for while blocks. Note that both for and while blocks count as "loops", so `break n` will break n loops regardless of whether those loops are while loops or for loops.

# *def* blocks

A `def` block defines a new command or overrides an existing one.

Here is an example of the basic syntax to define a command which takes no arguments:

    def foobar {
      # Put the command body here.
    }

`def` blocks may be nested within other blocks. For instance, it is possible to define a command within an `if` block:

    if $name Alex {
      def whoami {
        puts alex
      }
    } else {
      def whoami {
        throw Invalid user.
      }
    }
    whoami

Within a command body, the `return` command exits the command and sets its result. The `return` command takes any number of arguments and joins them with spaces. If you don't pass any arguments to `return`, it will yield the empty string. To demonstrate, this code will print "hello world":

    def hello_text {
      return hello world
    }
    print (hello_text)

The command's name is a regular token just like arguments to commands. For example, you can use a nested command as the command name:

    def (get functionName) {
        # Command body here.
    }

This example further highlights this feature:

    # Note that "Name" is capitalized; this makes it a global variable.
    set Name "hello"

    # Define a new command called "hello"
    def $Name {
        puts $Name
    }

    # Print out "bonjour"
    set Name bonjour
    hello

    # Print out "goodbye"
    set Name goodbye
    hello

A command can take arguments. On the line of a command definition, tokens after the command name are treated as variable names to which arguments are bound. For example, this command takes two arguments:

    def difference a b {
      return (- $a $b)
    }

Calling `difference 1 2` would yield `-1`; calling `difference 5 2` would yield `3`. Calling the `difference` command with more or less than 2 arguments would trigger an error.

Some commands take a variable argument list \(commonly known as varargs\). This can be achieved using "..." \(ellipsis\) after the last argument's name like so:

    def summation numbers... {
      set sum 0
      for number $numbers {
        set sum (+ $sum $number)
      }
      return $sum
    }

Note that the varargs argument is a newline-separated list of arguments. Arguments in the varargs list are escaped by adding backslashes before backslashes and turning newlines into `\n` escape sequences.

# Variable scope

Variables which begin with a capital letter are global. All other variables are local to their context. Every time a command is called, it is given a new context which shares global variables with its parent context but has a completely new set of local variables.

Note that code written directly inside a pragmash script outside of a command definition has its own local context. In essence, it is as if top-level code is within an invisible command definition.

Keep in mind that nested commands do not share the same context. For instance, this will trigger a "variable not found" error:

    def func1 {
      set a 3
      def func2 {
        # The variable "a" is not defined within this context.
        puts $a
      }
      func2
    }

This brief example demonstrates how variable scope works:

    # If f1 didn't get a new context every time it was called, the variable
    # "a" would be overwritten for each call.
    def f1 a {
      if (> $a 0) {
        f1 (- $a 1)
        puts $a
      }
    }
    # This will print: "1\n2\n3"
    f1 3
