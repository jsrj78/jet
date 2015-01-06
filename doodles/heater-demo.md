Heater Demo
===========

The heater demo was created to try out ideas and see how a concrete example would look given a particular language or
programming model.

Configuration
-------------
The demo consists of:
 - a heater controlled by a relay
 - two temperature sensors
 - a temperature threshold
 - a manual relay override to force the heater on or off

Functionality
-------------
 - the temperature sensors are read every 10 seconds
 - the value of the first temperature sensor is used to turn the heater relay on or off
 - both temperature sensor values are recorded in a database and are plotted on a web page
 - the control decision (based on temperature, prior to any manual override) and the relay
   position are recorded in a database and are plotted on a web page
 - the temperature threshold can be set on a web page
 - a manual override switch can be used on a web page to force the relay on or off
 - note that the second temperature sensor is just informational and not used in the control
 - the control uses hysteresis: if the relay is on, then it is turned off when the temperature rises N degrees
   above the set point; if the relay is off, then it is turned on if the temperature drops N degrees below
   the set-point. The value for N can be configured into the circuit.
