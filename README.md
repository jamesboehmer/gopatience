[![License](https://img.shields.io/badge/license-GPLv3.0+-blue)](https://opensource.org/licenses/GPL-3.0)

# Gopatience - A Collection of Patience Solitaire Games in Go

This project aims to be a small collection of console solitaire card games built on a common framework. The first such 
game is Klondike (or "Classic Solitaire").  Requires Go 1.14+.  To run type `go run cmd/klondike/klondike.go` from the 
root directory.

This project is a work in progress to port [pytience](https://github.com/jamesboehmer/pytience) from Python to Go.

## Klondike

The game is played via text command according to standard [rules](https://en.wikipedia.org/wiki/Klondike_%28solitaire%29).

### Starting the game

Run `klondike` in your terminal.  The screen will be cleared and replaced with a view of the foundation, tableau, and 
game stats.  If you've played before, it will automatically load the last state.  The tableau is displayed as fanned 
down piles, with the "bottom" card at the top of the column, and each card below it built on top of the previous.  The 
foundation piles are squared. 

The previous command is displayed in the prompt, and hitting return will execute it again.


```text
Score: 0
Stock: 24
Waste: []
Foundation: [♠]  [♦]  [♣]  [♥]
Tableau:
0    1    2    3    4    5    6
---  ---  ---  ---  ---  ---  ---
6♣   #    #    #    #    #    #
     A♦   #    #    #    #    #
          4♠   #    #    #    #
               5♦   #    #    #
                    7♥   #    #
                         4♦   #
                              8♠

klondike[type ? for help]>
```

### Deal a card from the stock

To deal a card, run `deal` (or just `d`).  A single card will be taken from the top of the stock, revealed, and placed 
fanned right on top of the waste pile.  The stock count will decrease by 1.  If you deal until the stock is empty, the 
waste will be returned to the stock and recycled.

```text
Score: 0
Stock: 23
Waste: [8♥]
Foundation: [♠]  [♦]  [♣]  [♥]
Tableau:
0    1    2    3    4    5    6
---  ---  ---  ---  ---  ---  ---
6♣   #    #    #    #    #    #
     A♦   #    #    #    #    #
          4♠   #    #    #    #
               5♦   #    #    #
                    7♥   #    #
                         4♦   #
                              8♠

klondike[deal]>
```

### Move tableau cards

To choose a card (or pile of cards) by using the `tableau` (or `t`) command:

`tableau <from_pile> [<card_num> [to_pile]]`

Arguments for pile and card numbers are zero-indexed.  If a chosen card has cards on top of it, that card and the cards 
 on top of it are moved together.  If `card_num` is omitted, it's assumed to be the top card in the pile.  If `to_pile` 
 is omitted, it will try to put single cards in the foundation, and then seek a pile in the tableau.
 
The `tableau` or `t` command can be omitted entirely, and integers will be assumed to be `tableau` arguments.  

Some examples: 

* `tableau 2 2 3` - Move the 3rd card from third pile to the fourth pile.  `4♠` will move below `5♦`, and the next card 
will be revealed.
* `t 1` - Move the top card from the third pile to the first available foundation or tableau pile.  The `A♦` will be 
moved to the foundation, and the next card will be revealed.  Moving a card from the tableau to the foundation will 
earn you 15 points.
* `0` - Move the top card from the first pile.  The `6♣` will find a home below `7♥`, and the first column will be 
empty.  Only a king may be placed there.
* `4 4` - Move the 5th card from the 5th column.  The `7♥` (and any cards on top of it) will move below the `8♠`

If all of the above commands are run in order, the tableau might look like this:

```text
Score: 15
Stock: 23
Waste: [8♥]
Foundation: [♠]  A♦  [♣]  [♥]
Tableau:
0    1    2    3    4    5    6
---  ---  ---  ---  ---  ---  ---
[ ]  2♥   #    #    #    #    #
          J♦   #    #    #    #
               #    #    #    #
               5♦   A♣   #    #
               4♠        #    #
                         4♦   #
                              8♠
                              7♥
                              6♣

klondike[tableau 4 4]>
```

### Move waste cards

Choose the top card from the waste using the `waste` (or `w`) command:

`waste [<tableau pile>]`

If the tableau pile argument is omitted, it will attempt to put the top waste card in the foundation (worth 10 points), 
or else the tableau (worth 5 points).  If there is a spot in the foundation _and_ the tableau, it would be advantageous 
to specify the tableau pile and make a second move from the tableau to the foundation, earning you 20 points combined.

### Move foundation cards

Sometimes you may need to move a foundation card to the tableau to open a spot for a waste card or another move.

`foundation <c(lubs)|d(diamonds|s(pades)|h(earts)> [<tableau pile num>]`

Specify the suit to pull the foundation card from, and optionally the tableau pile number.  If the pile number is 
omitted, it will seek a fit.  Making a move from the foundation will penalize you 15 points.

### Start a new game

Throw out the current game and create a new one with `new`, or `n`.

### Undo

Every move you make will be recorded.  You can undo all of them, one at a time, using the `undo` (or `u`) command.

### Solve

If the stock and waste are empty, and all of the tableau cards are revealed, `solve` will move a single card to the
foundation for you.  This is a convenience, since there are no moves left that can make the game unsolvable.

### Save

Each move you make will save the current game to `~/.gopatience/klondike.save`.  You can specify your own save file with:

`save [filename]`

### Load

Using the command `load [filename]`, you may load a previously saved game.

### Help

Using `help` or `?`, you can see all of the above commands, their syntax, and their descriptions.

### Quit

Quit the game using `quit`, `q`, or `ctrl-d`.

### Seeing the game state

The only undocumented command is `_dump` (formerly `_state`).  Use it to see a JSON representation of the game state.
