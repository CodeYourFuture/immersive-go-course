<!--forhugo
+++
title="File Parsing"
+++
forhugo-->

In this project we'll practice parsing data from files in different formats.

Many larger projects require reading configuration data or files, and often times it's also convenient to write a small program to understand or process some data, and being comfortable quickly doing so can be very helpful.

Timebox: 2 days

## Objectives:

- Write Go code to parse data in non-standard formats.
- Write tests for Go code.
- Become comfortable leveraging libraries to parse standard formats in slightly non-standard ways.

## Project

The `examples` sub-directory in this directory contains a number of data files, each of which contains the same data in a different format. Descriptions of each format can be found below. Each data file contains a data-set of player names, and their high scores, for some game.

We are going to write a program to analyse the data, and print out the names of the players with the highest and lowest scores.

Don't forget to write tests for your program, too.

### Data formats

### JSON

json.txt contains an array of JSON objects, each with a "name" and "high_score" key.

### Repeated JSON

repeated-json.txt contains lines of data stored in JSON format. Each line contains exactly one record, stored as an object. Lines starting with # are comments and should be ignored.

### Comma Separated Value

data.csv is a standard Comma Separated Value ("CSV") file. The format is well-documented online, and there are many libraries which support parsing it.

### Custom Binary

There are two files in a custom binary serialisation format.

> :memo: This section refers to a concept called "endianness". You can learn more about endianness [in this article](https://www.freecodecamp.org/news/what-is-endianness-big-endian-vs-little-endian/).

> :memo: This section refers to a character encoding called UTF-8. You can learn more about UTF-8 [on wikipedia](https://en.wikipedia.org/wiki/UTF-8). For this exercise, it's sufficient to know that UTF-8 is a way of encoding strings as bytes, and that if you read the bytes of a UTF-8 string in Go, you can use that value as a `string` without needing to change it.

> :memo: This section refers to a "null terminating character". When encoding a piece of data with variable length, we need to know how big it is. There are a few ways we can typically do that; the one we're going to use is that we'll specify "you'll know when the string is over, when you see a byte which is all zeros - the byte before that one was the last byte of the string".

The format is as follows:
* First two bytes of the file indicate endianness of numbers. If the bytes are FE FF, numbers in the file are stored in big endian byte order. If the bytes are FF FE, numbers in the file are stored in little endian byte order.
* Each record contains exactly four bytes representing the score as a signed 32-bit integer, in the above described endian format, then the name of the player stored in UTF-8 which may not contain a null character, followed by a null terminating character.

The tool `od` (which you can learn more about [here](https://man7.org/linux/man-pages/man1/od.1.html)) can be useful for exploring binary data. For instance, we can run:

```console
> od -t x1 projects/file-parsing/examples/custom-binary-le.bin
0000000    ff  fe  0a  00  00  00  41  79  61  00  1e  00  00  00  50  72
0000020    69  73  68  61  00  ff  ff  ff  ff  43  68  61  72  6c  69  65
0000040    00  19  00  00  00  4d  61  72  67  6f  74  00                
0000054
```

This prints each byte of the file, one at a time, represented as hexidecimal digits.

We can see in this example that the first byte is ff and the second is fe - according to our file format specification, that suggests the numbers here are stored in little endian byte order.

We can see the next four bytes contain 0a then three 00s, then three non-null bytes, then a null byte.


## Extra things to consider

### Variable-length data encoding

In the description of the binary serialisation format, we mentioned that there were different ways of encoding variable-length data.

The example that we used for strings was to use a terminating character - we specify that there's a character (or several characters in sequence) which isn't allowed to appear in our data (e.g. a null byte, one where all the bits are 0), and that we'll add that to the end of the data so you know it's the end.

We also used another technique, for storing the scores - we specified that the score is stored in a fixed number of bytes (specifically 4). 4 bytes is probably more space than we actually need to store our game's scores (4 bytes can store 4294967296 different values - really big scores!), and in fact most of our scores fit into 1 byte (256 different values), but specifying exactly 4 bytes is a simple rule, and gives us flexibility in case scores increase in the future.

Yet another technique is that we can write down the length of the variable length data in a fixed amount of memory before the data; i.e. to say "These 4 bytes say that the string after them will be 100 bytes long".

Each of these three approaches has different trade-offs - benefits they bring, and drawbacks they add.

Consider the trade-offs of each approach. What makes it a good approach? What makes it a bad approach? What kind of data and use-cases is each well-suited for?

Some things to consider:
* Are any more or less efficient in terms of how much space they use? Do any waste space?
* What limits do they apply to the kind of data we can actually store?
* Do any of the approaches make it easier/harder or faster/slower to parse data stored in that format?

## Avoiding writing code

While it's useful to be comfortable putting together ad-hoc programs to parse some data (and you should practice this!), one of the advantages of using existing formats of data is that there are often tools which can help us to do some parsing or analysis without even needing to write a program at all.

Two examples of this are [`jq`](https://stedolan.github.io/jq/) (which allows you to parse JSON using a custom query language), and [`fx`](https://github.com/antonmedv/fx) which allows you to write JavaScript snippets to manipulate JSON.

For example, you can use `jq` to answer the question "Who had the highest score" without needing to write a whole program:

```console
> jq -r '. | max_by(.high_score).name' file-parsing/examples/json.txt
Prisha
```

Or use `fx` to do the same, but using more familiar JavaScript as the query language:

```console
> fx file-parsing/examples/json.txt '.sort((l, r) => r.high_score - l.high_score)[0].name'
Prisha
```

Similarly, a program called [`csvq`](https://github.com/mithrandie/csvq) can be used to query CSV files in a SQL-like query language:

```console
> cat examples/data.csv | csvq 'SELECT * ORDER BY `high score` DESC LIMIT 1'
+--------+------------+
|  name  | high score |
+--------+------------+
| Prisha | 30         |
+--------+------------+
```

Spend some time experimenting with these tools:
* Write some interesting queries over the data.
* Try to work out what the limits are of using these pre-existing tools, and when you're more likely going to want to just write a custom program yourself.
