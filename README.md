# Go Adventure

![I have no idea what I'm doing](http://media.tumblr.com/tumblr_me0qrvGwDP1r3z80e.jpg)

## Wat?

One day this will be the finest twitter based novelty text adventure engine in all the lands.

Right now it's a git log of embarrassment as I learn a new programming language. When it's a bit further along, I'll ask for some help/feedback, but for right now just let me struggle along myself.

## Licence

I haven't decided yet, for right now consider it plain old copyright.

## TODO

  - [x] main game loop with input and output goroutines
  - [x] read & write from twitter API
  - [x] have "interactive" mode to mock twitter for dev
  - [ ] persistently store tweet reply status (for avoiding spam)
  - [ ] persistently store game state
  - [x] data structures for game -> scenes -> choices
  - [ ] file loading/parsing for game data
  - [ ] robust command parsing/handle synonyms etc

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
