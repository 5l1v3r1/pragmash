# pragmash

**pragmash** is a simple, flexible scripting language. It is made to be dropped into other systems (i.e. build systems, scripting environments, etc.).

The name "pragmash" is an abbreviation of the words "pragmatic" and "shell" put together. I want pragmash to be as practical as possible: it doesn't have to be elegant, just useful. I also want it to feel like writing for a shell (with barewords support and a simple command architecture).

# Syntax

I am trying to keep the syntax minimal. Here are all the built-in constructs the language will support.

## Commands

Scripts in pragmash run a series of commands. Each command has a string output. Similar to Bash and many other shells, a command is a name followed by space-delimited arguments. If an argument needs to have a space in it, the spaces can be escaped or the entire argument can be surrounded by quotes. Other escapes (like `\n`, `\"`, `\\`, etc.) are allowed as well. Here's some examples:

    cp /path1 /path2
    cp "/my path" /newpath
    ls /this/path\ has\ spaces\ without\ quotes
    ls /this/path\\has\\backslashes\nand\nnewlines

In addition, a command's output can be used as an argument to another command using backticks:

    read `replace http://aqnichol.com aqnichol google`

No whitespace is allowed directly after an opening backtick; the command's name must directly follow the opening backtick. Conversely, a closing backtick must be followed by a space, a newline, or an EOF. This makes it possible to nest backticks:

    read `replace http://aqnichol.com `cat old_domain.txt` google`

## Variables

Variables exist in a global scope, just like environment variables in Bash scripts.

The `set` pseudo-command sets a variable:

    set x `read http://aqnichol.com`

The `get` pseudo-command gets a variable:

    write ./home.html `get x`

As a shorthand for `get`, you can use a `$` followed by the variable name:

    write ./home.html $x

## If conditionals

An empty string is considered "false" in pragmash. Thus, you can use a basic `if` statement to check if a command outputs a non-empty string like this:

    if `read /might/be/empty` {
        puts The file wasn't empty.
    }

Certain commands might be crafted to output some sort of boolean result in this manner. For instance, an `exists` command might return "" if a file doesn't exist and "true" if it does:

    if `exists /some/path` {
        puts The file exists.
    }

But checking if a command outputs an empty string only goes so far. You can also check if any number of arguments are equal:

    if "Not found." `read http://google.com` {
        puts The page couldn't be found.
    }
    ...
    if "a" `echo a` a {
        puts It works!
    }

Finally, you can also use the `elif` and `else` keywords as expected:

    if "a" $x {
        echo It's A.
    } elif "b" $x {
        echo It's B.
    } else {
        echo It's not A or B.
    }

## While loops

A `while` loop repeats a block as long as a condition remains true. Conditions for `while` loops work exactly the same way as conditions for `if` statements.

    set x 0
    while `less $x 10` {
        set x `add $x 1`
        puts $x
    }
    puts Lift-off!

## For loops

There are no array types in pragmash; instead, arrays are represented as strings with newline delimiters. You can loop over the lines in a string like this:

    for x `ls /foo/bar` {
        puts Found file called $x
    }

This could be used for other purposes as well, such as iterating through a small range of numbers:

    for x `range 1 11` {
        puts $x
    }
    puts Lift-off!
