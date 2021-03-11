## TODO
link to the struct definitions in the source
check the types - std::string? Check the serialisation code. What does this say about headers and termination?
Also about encoding?


## Specificaion
OpenAPI? Ask

## Types
Tl;dr: I've given the sizes of fields in bytes. Numbers are little-endian (and almost certainly so is your target machine, so you don't care).

> The OpenRGB source is C++ and most things are typed `int` or `short`.
_Most_ C compilers have said `int` will be 32bits forver (even on 64bit machines you need `long` to get the full word size).
However the C standard does allow `int` to be any size up to and including the machine's word length.
`int` etc are meant to be used only when you want as much space as possible, eg for a counter.
`stdint.h`'s types like `uint32_t` should be used for known-fixed-size types.
So, what I show as eg 4 bytes might be 2 or 8 depedning on the compiler and platform of the OpenRGB server.
That said, most compilers let you redefine `int` with flags, or there's always `#define`.

> I can't see any endianness conversion code in OpenRGB.
This means what comes over the wire will be the OpenRGB server's system endianess.
Almost everything that's likely to run an OpenRGB _server_ is little-endian (amd64, most ARMs including Apple Silicon and the Raspberry Pi), but I guess someone could build it for a MIPS / PPC / etc home router or NAS.
It's probably a good idea to have your client decode from little-endian just in case.


### Strings
Strings in the protocol have a length header, then their body.
I assume (but can't confirm) that strings are 8-bit ASCII, which means you can also interpret them as UTF-8.
Note that, despite the presence of a length header, there is also a \0 terminator as well.

> You should of course always specify the length when injesting a string, as this is untrusted input.
This means you'll probably have to manually skip the termination byte.

I call these strings "bstring" after C++.

```
bstring:
  len : 2
  str : $len
  nul : 1
```

> This use of variable-length strings in the middle of records means you can't simply interpret the wire bytes as a record type (eg cast to a struct in C, TODO in go, TODO in python).

> NB: every message _received_ from OpenRGB uses these bstrings.
Some messages _send_ to the server do, some don't...

## Object Model
A natural object model might be:
* System
  * Devices (RGB Controller)
    * Zones
      * LEDs

The OpenRGB protocol doesn't provide this model, nor does it match the hierarchy in the OpenRGB UI or even code.
Read requests can only be made at the level of _Device_.
Write requests can be made at _Device_, _Zone_, and _LED_.

> The protocol is so device-centric that every message has a common header which includes _device ID_ (ignored for a few of the metadata commands like setting client name)

The protocol is explicity desinged to be as bandwidth efficient as possible so is normalised.

## Versioning

The protocol is versioned (as one big whole).
Versions are precise; no forwards or backwards compatibility is afforded.

This version number isn't in normal wire messages.
* You can read the server's version with the separate _GetProtocolVersion_ command.
* Your client doesn't report or negotiate its version with the server; you just have to see which version the server expects and send that. Because request messages don't contain the verion, the server blindly attempts to unpack the bytes you send into its structures. If you've sent the wrong version (or just malformed the request), random stuff will happen, possibly even dangerous things involving writing to memory and controller flash.


## Headers
For both request and response, each message's first 16 bytes are a common header format, as follows:

```
magic = "ORGB" : 4  # Trivia: OpenRGB's port number, 6742, are the digits corresponding to "ORGB" on a phone keypad.
device id : 4
command id : 4
body length : 4
```

In a request, you set the fields obviously.
* For commands with no body, set body length to 0.
* For commands not pertaining to a device (eg SetClientName), set device id to 0

In a response, device id and command id will be echoed back to you; assert that if you want.

Although device and command id are echoed back, notice that there's no sequence number - there isn't enough information to tie a response back to a request.
I _believe_ that the server is single-threaded, so a reply will be sent before the next request is read, guarenteeing the ordering on the wire.
Note that this will cause you problems if you try to mutli-thread your client and reuse the same socket.

## Responses
Commands the request data (like getting a device) elcit a response.
Commands that don't, do not recieve any kind of response.
Ie you don't even get a header and empty body as an ACK.

The OpenRGB protocol works over TCP, so unreliable networks are dealt with (or at least you're notified about issues - be sure to actually read the return codes from `send()`).

As for application-level, uh, good luck.
Cross your fingers there aren't any errors, and if you're paranoid I guess just re-assert the command periodically, as most (all?) are idempotent.


# Commands - Protocol Meta

## Get Protocol Version

### Request
Command: `40`

```
header : 16 # Set device id to 0
<empty body>
```

### Response
```
header : 16
protocol version : 4
```

## Set Client Name
Sets the name your client has in the client list in the UI, like a user agent.

### Request
Command: `50`

```
header : 16 # Set device id to 0
name : string
nul-terminator : 1
```

**NB** The string you send is NOT an OpenRGB headed "bstring".
However it DOES need the '\0' terminator.
Despite a) the protocol's heavy use of headed strings, and b) the presence of the string's length in the header (as _name_ is the only field), memory is blindly read.
If the UI shows your client's name as "foo���", that's a dead giveaway that you forgot the terminator.

### Response
None


# Commands - Devices


## Get Device Count

### Request
Command: `0`

```
header: 16
<empty body>
```

### Response
```
header : 16
devices count : 4
```


## Read Device

### Request
Command: `1`

```
header : 16
<empty body>
```

> Recall that the protocol is very device-centeric. This command needs no body, because every message's header contains a device ID.

### Response
```
header : 16
body lenght : 4 # This duplicates the length field in the header. Neither are actually necessary because everything field is of a known or stated size
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

Object Model
* A system has devices
  * A device has Modes
    * A Mode has color(s)
      * Color Cycle, etc: 0
      * Breathing, etc: 1+
      * Direct: n
        * In Direct Mode, you care about the individual LEDs
        * A device has Zones
          * Zones have (named) LEDS
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

per-led mode (eg breathing) - how represented?
orthogonal
* mode (static, breathing, etc). Call breathing etc dynamic. These have 0 (eg random), 1+ colors (eg flash one, flash between a few)
* color_mode (mode-specfic ~=1, per-led == n)
* confused by a lot of stuff not doing the 4th quadrent: per-led, non-static
Direct: per-led, static, don't flash (constant packet stream, so can do rapid transitions. Reqs on the driver: no fading between colors (instant), no flicker due to bit bang timing)
Custom: per-led, static, do flash (poersistent, not suitable for rapid transitions)
Static: mode-specific, static, do flash

Note ZoneType
* idk what Singular means if it can have >1 LED
* Planar is redundant, if Planan == non-empty matrix

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
