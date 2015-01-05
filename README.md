# pragmash

**pragmash** is a simple, flexible scripting language. It is made to be dropped into other systems (i.e. build systems, scripting environments, etc.).

The name "pragmash" is an abbreviation of the words "pragmatic" and "shell" put together. I want pragmash to be as practical as possible: it doesn't have to be elegant, just useful. I also want it to feel like writing for a shell (with barewords support and a simple command architecture).

# Syntax

Scripts in pragmash run a series of commands. Each command has a string output.

## Commands

A command is simply a name followed by space-delimited arguments. Here's some examples:

    read http://aqnichol.com
    cp /path1 /path2
    cp "/my path" /newpath
    wget http://aqnichol.com "/Users/alex/Desktop/downloaded files/foo"
    ls /foobar
    ls /this/path\ has\ spaces\ without\ quotes
    ls /this/path\\has\\backslashes\nand\nnewlines

In addition, a command's output can be used as an argument to another command:

    read `replace http://aqnichol.com aqnichol google`

The backticks wrap a command whose output can be fetched. Backticks may be nested, too:

    read `replace http://aqnichol.com `echo aqnichol` google`

**NOTE**: the commands shown in this section are just examples of commands that *could* exist.

## Variables

Variables exist in a global scope, just like environment variables in Bash scripts.

The `set` pseudo-command sets a variable:

    set x `read http://aqnichol.com`

The `get` pseudo-command gets a variable:

    write ./home.html `get x`

Other pseudo-commands like `for` will also set a variable.

## If conditionals

If a command outputs an empty string, it is considered a false output. Thus, you can use a basic if statement to check if a command outputs a non-empty string like this:

    if `read /might/be/empty` {
        puts The file wasn't empty.
    }

Certain commands might be crafted to output some sort of boolean result in this manner. For instance, an `exists` command might return "" if a file doesn't exist and "true" if it does:

    if `exists /some/path` {
        puts The file exists.
    }

But checking if a command is an empty string only goes so far. You can also check if a command outputs a certain value like this:

    if "Not found." `read http://google.com` {
        puts The page couldn't be found.
    }

## Loops

There are no array types in pragmash; instead, arrays are represented as strings with newline delimiters. You can loop over the lines in a string like this:

    for x `ls /foo/bar` {
        puts Found file called `get x`
    }
