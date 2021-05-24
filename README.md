# dmolisher

A Sauerbraten demo file parser/debugger.

## Usage

    $ gunzip --suffix=dmo --stdout path/to/demo/file.dmo | dmolisher -versions
    # file version: 1, protocol version: 260
    gamemillis, channel, data length, data (bytes in decimal)
    0, 1, 99, [2 22 114 101 105 115 115 101 110 0 5 1 33 128 88 2 61 2 0 3 255 37 0 1 0 0 0 0 64 100 100 100 1 2 20 20 10 10 20 0 1 1 0 0 0 0 66 100 100 100 1 2 20 20 10 10 20 0 255 3 0 65 113 117 97 124 115 112 52 110 107 124 0 103 111 111 100 0 1 3 1 99 114 52 115 104 124 115 112 52 110 107 0 103 111 111 100 0 1]
    666, 0, 16, [4 0 4 0 95 34 217 101 0 85 54 126 90 0 144 126]
    666, 1, 16, [88 0 13 18 64 100 100 100 1 2 20 20 10 10 20 0]
    [...]