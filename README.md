# gero
Gero is a terminal based UI application for interfacing nyaa.si. It allows you to 
open a torrent or batches at a time (by marking) from the command line. Gero is **not** a torrent client. 

gero uses [gocui](https://github.com/jroimartin/gocui) to create its console user interfaces.

![Example](/screenshots/example.png?raw=true)
![Help](/screenshots/help.png?raw=true)


## Why is gero slightly slower than the website?
Gero uses nyaa's rss feed instead of parsing the webpage directly. Evidently, the RSS feed could be 50% slower than sending a query through the website (around 500 ms slower).
It is possible that requesting for the HTML document and parsing it could be faster than
requesting for the RSS feed and parsing that. The slight delay does affect the main usage of the
program, but it is not a priority at the moment.

## No pagination supported. 
Gero is retrieving the information through a rss endpoint, which doesn't support pagination for various reasons. The only other way at the moment would be to get the actual webpage and parse it and do our own pagination with the url. A rewrite to using the webpage instead is inevitable if pagination is a priority.

## Known issues
- Jumping to the end of line sometimes jumps too far -- high possiblity this a gocui issue. 
To replicate: open terminal, decrease window size of terminal significantly, jump to end of a line that has a long line (>200 char). 