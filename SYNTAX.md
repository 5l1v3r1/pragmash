# Syntax Guide

This is the pragmash syntax guide. It uses human-readable langugae and tries to be concise.

# Charset

A pragmash script is encoded as UTF-8. The term "whitespace" refers to all the characters which are considered whitespace in UTF-8. A newline is the character represented by the number `10` (i.e. "\n").

# Lines

There are two types of "lines" in a pragmash script. **Physical lines** are separated by newline characters and include all of the leading and trailing whitespace. A **logical line** is similar to a physical line, but leading and trailing whitespace is removed and it is possible to represent a logical line across multiple physical lines using backslashes at the end of each physical line. For example, this counts as one logical line but two physical lines:

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

The `{` and `}` tokens \(also known as curly braces\) are essential to pragmash because they define blocks of code. Note, however, that a piece of pragmash code may contain curly braces inside of strings, command names, or comments. In these cases, curly braces are treated as regular characters. Curly braces only denote blocks of code when they are situated correctly. In these cases, every `{` must be the last character on a logical line and every `}` must be the first character on a logical line. Blocks of code can also be nested, such as in this case:

    if a {
      while b {
        # something here
      }
    }

# Strings

Every argument to every command is a string. Variables contain strings and their names are themselves strings. Command names, too, are themselves strings. Thus, it is very important to be able to express strings concisely in pragmash code.

A string can be expressed without quotes, provided it has no unescaped whitespace. In addition, strings can be surrounded by `'` \(single quotes\) and `"` \(double quotes\). In these cases, unescaped whitespace can be present.

An escaped character within a string begins with a backslash. To escape a space, use `\ `. To escape a backslash, use `\\`. Inside of single quotes, it is necessary to escape single quotes as `\'`. Likewise, you should use `\"` if you want to include quotes inside of a quoted string.

Escapes can also be used to write out characters which a programmer cannot easily enter from their keyboard. These are the available escape sequences, not including the ones discussed above:

| Sequence   | Description         | Hex Code (if applicable) |
|------------|---------------------|--------------------------|
| \?         | Question mark       | 0x3f                     |
| \a         | Bell                | 0x07                     |
| \b         | Backspace           | 0x08                     |
| \f         | New page            | 0x0c                     |
| \n         | New line            | 0x0a                     |
| \r         | Carriage return     | 0x0d                     |
| \t         | Tab                 | 0x09                     |
| \v         | Vertical tab        | 0x0b                     |
| \nnn       | 3-digit octal value |                          |
| \xnn       | 2-digit hex value   |                          |
| \unnnn     | Unicode value       |                          |
| \Unnnnnnnn | Unicode value       |                          |
