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
<modes>
    
</modes>
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
    Color : 3
color_count : 2
<colors>

o
</colors>

rename to lightload. example: find all the LEDs and turn them on one at a time, 1s each
