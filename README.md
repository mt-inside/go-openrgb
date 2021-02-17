FIXME
* go:generate not running

TODO
* Put device.Colors in device.LEDs as "actualColor" or summat
* Put device.LEDs in their Zones (just linear pack)
* Put Zones under Mode Direct
  * For other Modes, fake up a "Mode Specific" Zone, like the UI does
* ability to change color on: all devices / Device / Zone / LED (like in the UI) - should get all LEDs in it. SetColor() func
* Get everything by name (with func, or expose a map)
* example: find all the LEDs and turn them on one at a time, 1s each
* rebase loadsofrgb onto this
* rename to lightload. 
