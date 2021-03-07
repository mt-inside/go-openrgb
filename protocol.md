Work in progress!
TODO: to a GH wiki

## Types
bstring
```
len : 2
str : len
nul : 1
```
Note that, despite the presence of a length header, there is a \0 terminator

## Read Device
Command: TODO

### Request

### Response
```
TotalLen : 4 - Not needed since everything else is of a predictable / stated size
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

Object Model, logically
* A system has devices
  * A device has Modes
    * A Mode has color(s)
      * Color Cycle, etc: 0
      * Breathing, etc: 1+
      * Direct: n
        * In Direct Mode, you care about the individual LEDs
        * A device has Zones
          * Zones have (named) LEDs
            * LEDs have a color
The API doesn't expose it like this, instead:
* device
  * modes
    * color - only for zero/one/few color modes (mode_colors_mode_specific)
      * mode_colors_per_led -> use device->colours. CalcProgrammer1 isn't stupid, did this becuase there's >1 per-led mode on some devices and he was worried about multiple n-sized arrays.
    * Value: opaque (driver-level); usually the hw's value for the mode.
  * zones -  Zone is just their sizing, and physical layout if applicable "matrix" (keyboards). (1D and 2D zones are both represented as vectors.)
  * LEDs - all of them for the device. Match them to zones by counting the indexes
    * The int field on these is not color. It's an opaque value (driver-level), often used to stash the register/address of the LED
  * Colors - all of them for the device. Match them to the LEDs by comparing indecies
    * THis is what you, or another sdk / gui, has set for direct mode, so morally it is direct's colors
    * it doesn't let you sample colors, eg if the device is in Breathing, you can't watch the brightness go up and down
BUT! remeber this is the READ object model.
* modes' colrs are the ones you'll switch back to if you activate that mode
* device->colours is NOT the current value, eg set by an automated mode?
* direct mode isn't a thing in the HW controller, because it isn't flashed, and often needs a constant stream of packets etc, so direct mode doesn't have an intrinsic set of colors that can be switched back to
  * recall that this ISN'T the object used for *setting* colors

to have direct control, try in order: direct (per-led, don't flash), custom (per-led, flash), static (one col, flash?)

orthogonal things:
* mode (static, breathing, etc). Call breathing etc dynamic. These have 0 (eg random), 1+ colors (eg flash one, flash between a few)
* color_mode (mode-specfic ~=1, per-led == n)
* confused by a lot of stuff not doing the 4th quadrent: per-led, non-static
Direct: per-led, static, don't flash (constant packet stream, so can do rapid transitions. Reqs on the driver: no fading between colors (instant), no flicker due to bit bang timing)
Custom: per-led, static, do flash (poersistent, not suitable for rapid transitions)
Static: mode-specific, static, do flash

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
    * This just names them
    * The actual color of the LED is in the Device Color array at the corresponding index
  * Devices have Colors
    * These are the actual colors of the LEDs
    * They seem to match to LEDs (and thus pack into Zones) just based on index correspondance
