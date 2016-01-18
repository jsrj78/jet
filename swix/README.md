# JET/Swix

Embedded micro-controllers are flexible and low-cost. Their behaviour is
determined by the software they run. The predominant programming language for
this is C/C++. Some alternatives exist, such as [Forth][af] (AVR), [BasCom][bc]
(AVR), [Maximite][mm] (PIC), and [Espruino][es] (JavaScript/ARM), but they're
not as widely used.

**Swix** is an experiment to bring a high-level flow-based engine to these
low-end µC environments. It's written in C, but the application logic is
data-driven, to be defined with a visual diagram editor at a later stage. The
Swix “engine” is uploaded to flash memory once, and then a relatively small
“circuit” can be loaded to make it perform a specific task.

A circuit in turn is made up of “gadgets”, inter-connected by a set of “wires”.
Gadgets are pre-defined modules, implementing core functionality - which
depending on the embedded µC can include I/O pin control, hardware peripherals
such as timers, UARTs, I2C, SPI, as well as primitive functions for logic,
comparison, arithmetic, and more. Circuits can be re-used as gadgets inside
larger circuits, this allows aggregating more complex functionality in a
clearly-defined modular fashion. Building a new application consists of picking
existing gadgets and circuits, adding new gadget types to the core if needed,
and then “wiring them up” to implement the desired behaviour.

Swix borrows heavily from the concepts of [Pure Data](http://puredata.info/) -
the Pd introduction says it all:

> Pure Data (aka Pd) is an open source visual programming language. Pd enables
> musicians, visual artists, performers, researchers, and developers to create
> software graphically, without writing lines of code. Pd is used to process
> and generate sound, video, 2D/3D graphics, and interface sensors, input
> devices, and MIDI. Pd can easily work over local and remote networks to
> integrate wearable technology, motor systems, lighting rigs, and other
> equipment. Pd is suitable for learning basic multimedia processing and visual
> programming methods as well as for realizing complex systems for large-scale
> projects.

Swix is also influenced by the design and implementation of the [Lua][lu] and
[Factor][fa] programming languages, among others.

Swix is currently (early 2016) 99% fantasy and 1% code. This is an attempt to
liquify some of that vapour.

See `tests/*.cpp` for some early examples.

   [af]: http://amforth.sourceforge.net/
   [bc]: http://www.mcselec.com/index.php?option=com_content&task=view&id=14&Itemid=41
   [mm]: http://www.geoffg.net/maximite.html
   [es]: http://www.espruino.com/
   [lu]: http://www.lua.org/
   [fa]: http://factorcode.org/
