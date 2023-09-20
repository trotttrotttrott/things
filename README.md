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

### Done and Pause

Add `done: true` to a thing to mark it as done. It will be removed from the
list.

Add `pause: true` to pause it. This just dims its color a little to indicate you
can skip it for now.

### Time Tracking

The amount of time a thing is open in your editor is tracked with csv files in
`~/.things/time/`. Each thing has a file here as well. A row is added with the
open and close time every time you open a thing in order to calculate the
cumulative time you've spent working on a thing.

## Actions

n = open new thing in $EDITOR

enter = open thing in $EDITOR

k = cursor up

j = cursor down

g = set cursor to first thing

G = set cursor to last thing

d = toggle show done things

\# = toggle line numbers

q = quit
