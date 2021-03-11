/*
Package wire contains types that represent OpenRGB's wire protocol.
There isn't exactly a 1:1 mapping with the records on the wire (they don't "byte-match", ie we don't directly decode one into the other).
But on the write side they contain every field needed to build a message to send to OpenRGB.
And on the read side they store everything that will be necessary in the future for building write messages, like the internal opaque fields that clients are required to preserve.

For normal interaction with an OpenRGB server, you should use package model instead.

There are no official docs for the wire protocol; you're told to read an existing SDK.
I've done my best to document it here:
*/
package wire
