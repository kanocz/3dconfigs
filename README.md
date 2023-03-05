# 3dconfigs
Configurations for different printers for PrusaSlicer

The first PrusaSlicer config bundle is for Sapphire Plus (with standard hotend).

Also added golang "script" which converts PrusaSlicer gcode thumbnail (prepared for Prusa Mini) to
Robin nano MKS firmware format. Can be used as output filter directly.

One more "script" which converts PrusaSlicer gcode thumbnail (prepared for Prusa Mini) to
Chitu board firmware format and also converts M73 to M2100/M2101 (for Qidi X-Max or X-Plus).
Part of output image format code "converted" from python from https://github.com/alkaes/QidiPrint/blob/master/ChituCodeWriter.py .
Can be used as output filter directly.

For snapmaker J1 tool for auto-upload https://github.com/kanocz/j1upload can be used as well as
https://github.com/macdylan/Snapmaker2Slic3rPostProcessor post processor (be sure to put uploader
as the last post-processing script).
