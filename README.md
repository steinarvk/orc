# Orc

Orc (short for "orchestrator" or something along those
lines) is a Golang library to manage initialization.

Specifically, it is designed to manage "singleton" libraries
that add functionality to a command-line or server binary
after being initialized.

These libraries can manage their own flags (using the
Cobra library) and may have dependencies in between
themselves and require a certain initialization order.
The orc library is the glue code that calls them
at startup to allow them to set up their flags, and
figures out the initialization order.

## Goal

The goal of Orc is to make it possible to write very
short Go binaries -- especially server binaries -- that
come with an ever-expanding set of batteries included,
making them pleasant to debug and work with.

The Orc libraries handle all the boilerplate and
essentially initialize themselves, letting the application
code focus on the actual application logic.

## Don't use this in its current state

Orc is a highly opinionated framework, and one that is
still taking form and in pretty sketchy shape at the moment.

I'm releasing it to be able to use it as glue code for
other hobby projects. I'm working out the kinks required
for those projects as I go, but I don't recommend anyone
else use it at the moment. Expect usability quirks, bugs,
and API instability. At this point, even the expectations
for code using Orc may change.

Using this library in your own code is strongly
recommended against at the moment.

## Dealing with Orc code

See STANDARDS.md for the peculiarities of Orc code.
Don't use Orc without reading that document. (In fact,
at the moment, don't use Orc at all. See above.)

## Authorship

Orc was written by me, Steinar V. Kaldager.

This is a hobby project written primarily for my own usage.
Don't expect support for it. It was developed in my spare
time and is not affiliated with any employer of mine.

It is released as open source under the MIT license. Feel
free to make use of it however you wish under the terms
of that license.

## License

This project is licensed under the MIT License - see the
LICENSE file for details.
