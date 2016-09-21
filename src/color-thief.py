from colorthief import ColorThief
import sys
color_thief = ColorThief(sys.argv[1])
palette = color_thief.get_palette(color_count=6)
print palette
