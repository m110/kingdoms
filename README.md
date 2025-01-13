# Kingdoms

A strategy game written in Go and Ebitengine.

https://github.com/user-attachments/assets/f85316ad-bdf6-45ee-90cd-0a07719b5931

## Sandbox (terrain generation)

![](.github/sandbox-1.png)

![](.github/sandbox-2.png)

## Layers

* Domain — should be Kingdoms-specific data.

## Render system

There are two kinds of render area:

* Board - all elements affected by the camera. They are rendered in chunks.
* UI - all elements that are not affected by the camera. They are rendered in one piece, with absolute position.
