# Kingdoms

## Layers

* Domain — should be Kingdoms-specific data.

## Render system

There are two kinds of render area:

* Board - all elements affected by the camera. They are rendered in chunks.
* UI - all elements that are not affected by the camera. They are rendered in one piece, with absolute position.
