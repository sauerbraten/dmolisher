# dmolisher

A Sauerbraten demo file parser. Feed it a gunzipped .dmo file and it gives you the recorded network packets as CSV records.

## Usage
    
    $ ./dmolisher -help
    Reads uncompressed demo file on stdin and emits CSV on stdout.
    Usage of ./dmolisher:
      -channel int
            print only packets sent on channel (0/1/2) (default -1)
      -hex
            print data bytes in hexadecimal instead of decimal
      -versions
            print file and protocol versions



## Example

    $ gunzip --suffix=dmo --stdout path/to/demo/file.dmo | dmolisher -versions
    # file version: 1, protocol version: 260
    gamemillis, channel, data length, data
         0, 1, 18, 2 22 99 111 109 112 108 101 120 0 0 1 33 128 88 2 37 255
         0, 1, 12, 3 0 112 49 120 0 103 111 111 100 0 1
       823, 1, 18, 88 0 15 18 1 100 100 25 0 6 0 0 0 0 1 40 32 5
       854, 0, 16, 4 0 12 0 193 21 254 35 128 36 153 127 90 0 144 126
       854, 1,  5, 88 0 2 32 8
       888, 0, 16, 4 0 12 0 193 21 254 35 128 36 153 127 90 0 144 126
       922, 0, 16, 4 0 12 0 193 21 254 35 128 36 153 127 90 0 144 126
    [...]