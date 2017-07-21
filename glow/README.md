# JET/Glow

A dataflow engine in Go, inspired by [Pure Data](http://puredata.info) (Pd).

[![GoDoc](https://godoc.org/github.com/jeelabs/jet/glow?status.svg)](https://godoc.org/github.com/jeelabs/jet/glow)
[![Build Status](https://travis-ci.org/jeelabs/jet.svg?branch=master)](https://travis-ci.org/jeelabs/jet)
[![license](https://img.shields.io/github/license/jeelabs/jet.svg)](http://unlicense.org)

Glow combines building blocks (called _gadgets_) into a runnable system (called
a _circuit_). Each gadget can have _inlets_ which accept messages and _outlets_
which emit messages. Explicit connections between them determine the message
flow and processing order. A message can currently be an integer, a string, nil,
or a vector of these. Gadgets have to be implemented in Go, but Circuits can be
used as additional building blocks for convenient nesting. Circuits can also be
instantiated from a text description, called a _design_.

These terms were chosen to resemble the vocabulary of electronics ("chip" was
rejected in favour of "gadget"). Note that in Pd, a gadget is called an
"object", a circuit is a "patch", and a design is an "abstraction".

Running the code:

    $ go test ./tests
    ok    github.com/jeelabs/jet/glow/tests    0.009s
    $ 
