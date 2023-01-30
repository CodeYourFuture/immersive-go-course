<!--forhugo
+++
title="File Parsing"
+++
forhugo-->

In this project we'll practice parsing data from files in different formats. Often times it's convenient to write a small program to understand or process some data, and being comfortable quickly doing so can be very helpful.

Timebox: 2 days

## Objectives:

- Write code to parse data in non-standard formats.
- Become comfortable leveraging libraries to parse standard formats in slightly non-standard ways.

## Project

The `examples` sub-directory in this directory contains a number of data files, each of which contains the same data in a different format. Descriptions of each format can be found below. Each data file contains a data-set of player names, and their high scores, for some game.

We are going to write a program to analyse the data, and print out the names of the players with the highest and lowest scores.

### Data formats

### JSON

json.txt contains an array of JSON objects, each with a "name" and "high_score" key.

### Repeated JSON

repeated-json.txt contains lines of data stored in JSON format. Each line contains exactly one record, stored as an object. Lines starting with # are comments and should be ignored.

### CSV

data.csv is a standard CSV file. The format is well-documented online, and there are many libraries which support parsing it.

### Custom Binary

There are two files in a custom binary serialisation format. The format is as follows:
* First two bytes of the file indicate endianness of numbers. If the bytes are FE FF, numbers in the file are stored in big endian byte order. If the bytes are FF FE, numbers in the file are stored in little endian byte order.
* Each record contains exactly four bytes representing the score as a signed 32-bit integer, in the above described endian format, then the name of the player stored in UTF-8 which may not contain a null character, followed by a null terminating character.

The tool `od` can be useful for exploring binary data. For instance, we can run:

```console
% od -t x1 projects/file-parsing/examples/custom-binary-le.bin
0000000    ff  fe  0a  00  00  00  41  79  61  00  1e  00  00  00  50  72
0000020    69  73  68  61  00  ff  ff  ff  ff  43  68  61  72  6c  69  65
0000040    00  19  00  00  00  4d  61  72  67  6f  74  00                
0000054
```

This prints each byte of the file, one at a time, represented as hexidecimal digits.

We can see in this example that the first byte is ff and the second is fe - according to our file format specification, that suggests the numbers here are stored in little endian byte order.

We can see the next four bytes contain 0a then three 00s, then three non-null bytes, then a null byte.
