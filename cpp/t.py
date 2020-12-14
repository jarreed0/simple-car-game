#!/usr/bin/env python3
print("static const unsigned char TILE_IMAGE_DATA[] = {{{}}};".format(
        ",".join(str(b) for b in open("tile.png", "rb").read())))
