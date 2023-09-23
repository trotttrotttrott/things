# Things

A terminal application to help keep track of all the things you're juggling at
work or life in general.

## Thing Files

Things are files stored in `~/.things/things/`.

When you create a new thing, a file is created and it's opened with `EDITOR`.

```
---
title: Thing
type: chore
priority: 3
---

Respond to that email from so and so.
```

### Types

Things have types that you define in `~/.things/types/`.

For example, you could define `~/.things/types/chore.md` as:

```markdown
---
color: '#33ffc1'
---

Random, small task. Something otherwise untracked.
```

### Priority

You can assign things a positive integer to represent priority - 0 being highest
priority.

### Done, Pause, Today

Add `done: true` to a thing to mark it as done. It will be removed from the
list.

Add `pause: true` to pause it. This just dims its color a little to indicate you
can skip it for now.

Add `today: true` to indicate that thing needs to be addressed today. These will
be in bold.

You can filter by each of these. This is documented below.

### Time Tracking

The amount of time a thing is open in your editor is tracked with csv files in
`~/.things/time/`. Each thing has a file here as well. A row is added with the
open and close time every time you open a thing in order to calculate the
cumulative time you've spent working on a thing.

## Actions

```
// navigation

k = cursor up

j = cursor down

ctrl+u = cursor up 5

ctrl+d = cursor down 5

g = set cursor to first

G = set cursor to last

// filter

A = clear filter (default)

D = done: true

P = pause: true

T = today: true

// sort

a = sort things by age

p = sort things by priority (default)

t = sort things by type and priority

// display

# = toggle line numbers

// edit

n = open new thing in $EDITOR

enter = open thing in $EDITOR

// quit

q = quit
```
