FIXME
* 

TODO
* openrgb SetActiveMode PR - Separate PR so we can argue about it: Also if you'd accept a change in API semantics, I'd propose changing the current UPDATEMODE to NOT _set the active mode_, and I'd also rename it to SETMODEDESCRIPTION
* The FooList thing is nasty, just have the collections be maps
* example: find all the LEDs and turn them on one at a time, 1s each
* separate and opensourec logging
* rebase loadsofrgb onto this
* opensource

TODO: publish go-sensors (with get() and observable interfaces) (or pick a name) (to openrgb discord), make lightload example with rxgo
Examples
* X current main
* lightload - in progress
* X redshift

Separate PR so we can argue about it: Also if you'd accept a change in API semantics, I'd propose changing the current UPDATEMODE to NOT _set the active mode_, and I'd also rename it to SETMODEDESCRIPTION

write code for: 2 sticks ram (stacked, swap), 4 sticks ram (all 3 separate + free, scaled to 100%). Needs lib naming enhancements

what am lib? what am point?
* this is mostly about werite. Reading doesn't seem too useful programatically, and I don't expect anyone to build an interactive client using this
* but, to write, we need to know the "schema" - the devices and their configured sizes, so we need to read once - schema and LEDs come over the same request
* we assume we're the only writer, and can just overwrite evertyhing. This massively simplifies things, otherwise we'd need a state object, TF style 3-way diff, and either overwrite or prompt for conflicts
* becuase of this, we can just smash our object model up to the server
* read colors into it, should you care, and modify from there
* as a convenience, we'll provide a *local* diff, against what you last read (overwritten by our writes).
* no schema write support for now

* session: getSchema() [build middle objects, return new, no attempt to merge for now] -> syncColors() [into middle object] -> modify -> diff -> writeCOlors()

lib: impliment object model above.
* How to do? I think:
  * keep current structs, all internal, named to wireFoo, decode into them (not least, binary.read in future)
  * make new structs for desired object model, build from wireFoo
    * drop the opaque Value fields
    * hide all fields, do mut-ness with getters and (fewer) setters
    * do the diff thing too, and a pretty-print of it
    * setters for resizing Zones' LEDs and Modes' color-wheels
  * prolly have to make wireWriteFoo structs as well to update
* assert that mode-per-led never has colours in it on the wire
* throw away led.value, mode.flags etc, at the wireReadFoo level
* doing it this way won't mask any features I don't think - can still set inactive modes' colors, per-led static and dynamic
* my lib is kinda focussed on setting colors *now*, rather than programming up the flash, so it'll warn you if you're coloring inactive modes etc.
  * doc: if you wanna programme up your controller's flash, do this
  * doc: if you wanna be setting your LEDs in real-time, do this
  * doc: basics of API like the orthogonal modes. Link to wire protocol doc.
lib: should be able to set color at
* device - applies to currently-active mode (log.info this) (the one color for a single-col mode, the color for every LED in direct, error if in a zero-color mode). Color is: color, or random
  * mode - same as above, but for a named mode (log.info this if mode is not the active one)
    * zone (applies to direct only, log.info this if direct isn't active)
      * led (ditto)

TODO to GH. Add
* get building on windows...


what am lib? what am point?
* this is mostly about werite. Reading doesn't seem too useful programatically, and I don't expect anyone to build an interactive client using this
* but, to write, we need to know the "schema" - the devices and their configured sizes, so we need to read once - schema and LEDs come over the same request
* we assume we're the only writer, and can just overwrite evertyhing. This massively simplifies things, otherwise we'd need a state object, TF style 3-way diff, and either overwrite or prompt for conflicts
* becuase of this, we can just smash our object model up to the server
* read colors into it, should you care, and modify from there
* as a convenience, we'll provide a *local* diff, against what you last read (overwritten by our writes).
* no schema write support for now

* session: getSchema() [build middle objects, return new, no attempt to merge for now] -> syncColors() [into middle object] -> modify -> diff -> writeCOlors()

TODO to GH. Add
* get building on windows...
* wtf fan curves (reboot to look. custom silent to 60?)
  * ask on reddit/amd - do I need to wait for the kernel to see all my fans? Doesn't seem to be any config for this module.

lib: impliment object model above.
* How to do? I think:
  * keep current structs, all internal, named to wireFoo, decode into them (not least, binary.read in future)
  * make new structs for desired object model, build from wireFoo
    * drop the opaque Value fields
    * hide all fields, do mut-ness with getters and (fewer) setters
    * do the diff thing too, and a pretty-print of it
    * setters for resizing Zones' LEDs and Modes' color-wheels
  * prolly have to make wireWriteFoo structs as well to update
* assert that mode-per-led never has colours in it on the wire
* throw away led.value, mode.flags etc, at the wireReadFoo level
* doing it this way won't mask any features I don't think - can still set inactive modes' colors, per-led static and dynamic
* my lib is kinda focussed on setting colors *now*, rather than programming up the flash, so it'll warn you if you're coloring inactive modes etc.
  * doc: if you wanna programme up your controller's flash, do this
  * doc: if you wanna be setting your LEDs in real-time, do this
  * doc: basics of API like the orthogonal modes. Link to wire protocol doc.
lib: should be able to set color at
* device - applies to currently-active mode (log.info this) (the one color for a single-col mode, the color for every LED in direct, error if in a zero-color mode). Color is: color, or random
  * mode - same as above, but for a named mode (log.info this if mode is not the active one)
    * zone (applies to direct only, log.info this if direct isn't active)
      * led (ditto)
