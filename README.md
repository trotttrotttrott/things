# Things

A terminal application to help keep track of all the things you're juggling at
work or life in general.

## Things Dir

By default it's `~/.things/`.

You can override this with `THINGS_DIR`.

## Thing Files

Things are files in `THINGS_DIR/things/`.

When you create a new thing, a file is created and it's opened with `EDITOR`.

```
---
title: Thing
type: chore
priority: 3
---

Do a thing!
```

### Priority

You can assign things a positive integer to represent priority - 0 being highest
priority.

### Done, Pause, Today

Add `done: true` to a thing to mark it as done. It will be removed from the
default list.

Add `pause: true` to pause it. This just dims its color a little to indicate you
can skip it for now.

Add `today: true` to indicate that thing needs to be addressed today. These will
be in bold.

You can filter by each of these. This is documented below.

### Pin

Add `pin: true` to pin it to the top of the list.

## Type Files

Things have types that you define in `THINGS_DIR/types/`.

For example, you could define `THINGS_DIR/types/chore.md` as:

```markdown
---
color: '#33ffc1'
---

Random, small task. Something otherwise untracked.
```

## Time Files

The amount of time a thing is open in your editor is tracked with csv files in
`~/.things/time/`. Each thing has a file here as well. A row is added with the
open and close time every time you open a thing in order to calculate the
cumulative time you've spent working on a thing.

## Thing Directories

For complex things that require extended research, multiple files, or
collaboration with AI tools, you can create a directory for it.

Directories are stored in `THINGS_DIR/things-deep/{thing-id}/`.

Things that have directories show a `*` indicator in the list.

## Actions

```
// mode

> = switch between "thing" and "type" modes

/ = search ("thing" mode only)

? = toggle help

// navigation

k = cursor up

j = cursor down

ctrl+u = cursor up 5

ctrl+d = cursor down 5

g = set cursor to first

G = set cursor to last

// filter ("thing" mode only)

C = current, done: false (default)

D = done: true

A = all, no filter

P = pause: true

T = today: true

// sort ("thing" mode only)

a = sort things by age

p = sort things by priority (default)

t = sort things by type and priority

// display

# = toggle line numbers

// edit

n = open new thing in $EDITOR ("thing" mode only)

enter = open thing or type in $EDITOR

E = open thing directory in $EDITOR ("thing" mode only)

ctrl+e = open thing time file in $EDITOR ("thing" mode only)

ctrl+x = delete thing ("thing" mode only)

// quit

ctrl+c, q = quit
```
