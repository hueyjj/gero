# gero


## Why is gero slightly slower than the website?
Apparently, the RSS feed could be 50% slower than sending a query through the website. 
It is possible that requesting for the HTML document and parsing it could be faster than
requesting for the RSS feed and parsing that. The slight delay does affect the main usage of the
program, but will consider alternatives to speed things up.


## No pagination supported. 
We are retrieving the information through a rss endpoint, which doesn't support pagination for some reason. The only other way at the moment would be to do GET request for the html and parse that; rewrite seems inevitable if it is actually necessary.

## Known issues
- Jumping to end of line sometimes jumps too far -- strongly believe this a gocui issue. 
To replicate: open terminal, decrease window size of terminal significantly, jump to end of a line that has a long line (>200 char). 