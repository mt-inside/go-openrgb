## Types
bstring
```
len : 2
str : len
```
Note that, despite the presence of a length header, str contains a \0 terminator

## Read Device
Command: TODO

### Request

### Response
```
??? : 4
type : 4
name : bstring
description
version
serial
location
mode_count : 2
active_mode_index : 4
    Name : bstring
    Value : 4
    Flags : 4
    MinSpeed : 4
    MaxSpeed : 4
    MinColors : 4
    MaxColors : 4
    Speed : 4
    Direction : 4
    ColorMode : 4
    ColorCount : 2
        Color : 4
zone_count : 2
    Name : bstring
    Type : 4
    MinLEDs : 4
    MaxLEDs : 4
    TotalLEDs : 4
    MatrixSize : 2
    Matrix : MatrixSize
led_count : 2
    Name : bstring
    Color : 4
ColorCount : 2
    Color : 4
```

Semantics
* Colors
  * Modes may have colors, specifically
    * None, if the colorMode for the Mode is Random or None - no user color input possible
    * One, if the colorMode for the Mode is Mode-Specfic - user gets to specify one color
    * None, if the colorMode for the Mode is Per-LED - colors are instead stored in the Device's color array
    * Colors are remembered by each mode - in the color array of the Mode if it only has one, and in the color array of the Device if the mode is Direct
  * Devices have LEDs
    * These seem to give names to LEDs, not that you can rename them
    * Devices have Zones, which have a size (either device-defined or user-defined).
      * Zones don't contain LEDs in the object model, but it's just a linear carve-up
  * LEDs have Colors
    * I have no idea what these do; they don't seem set to anything sensible
    * The actual color of the LED is in the Device Color array at the corresponding index
  * Devices have Colors
    * These are the actual colors of the LEDs
    * They seem to match to LEDs (and thus pack into Zones) just based on index correspondance

TODO object model
* Put device.Colors in device.LEDs as "actualColor" or summat
* Put device.LEDs in their Zones (just linear pack)
* Put Zones under Mode Direct
  * For other Modes, fake up a "Mode Specific" Zone, like the UI does
* ability to change color on: all devices / Device / Zone / LED (like in the UI) - should get all LEDs in it. SetColor() func
* Get everything by name (with func, or expose a map)

rename to lightload. example: find all the LEDs and turn them on one at a time, 1s each
