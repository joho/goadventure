# Go Adventure

![I have no idea what I'm doing](http://media.tumblr.com/tumblr_me0qrvGwDP1r3z80e.jpg)

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

## Wat?

One day this will be the finest twitter based novelty text adventure engine in all the lands.

Right now it's a git log of embarrassment as I learn a new programming language. When it's a bit further along, I'll ask for some help/feedback, but for right now just let me struggle along myself.

## Build Status

[![wercker status](https://app.wercker.com/status/c5ab59c3a612589b5b804fb6814e535c/m "wercker status")](https://app.wercker.com/project/bykey/c5ab59c3a612589b5b804fb6814e535c)

## Licence

I haven't decided yet, for right now consider it plain old copyright.

## TODO

  - [x] main game loop with input and output goroutines
  - [x] read & write from twitter API
  - [x] have "interactive" mode to mock twitter for dev
  - [x] persistently store tweet reply status (for avoiding spam)
  - [x] persistently store game state
  - [x] data structures for game -> scenes -> choices
  - [ ] file loading/parsing for game data
  - [ ] robust command parsing/handle synonyms etc
  - [ ] Handle twitter "duplicate" errors (code 187)
  - [ ] Write the most baller choose your own adventure script in all the lands.

## Code review notes

  * tomb package for proper stopping on control channels
  * using this/self is kinda discouraged (rob pike, single chars, fuck that)
  * go fmt -r (renames stuff by magic)
  * start developing against 1.1 (compatible and faster)
    - do not set go root
    - build from source (or homebrew)
    - use abs path
  * check out "go check" or https://github.com/remogatto/prettytest
  * if you're writing a library do executable examples
  * read https://gobyexample.com/
