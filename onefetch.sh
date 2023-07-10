#!/usr/bin/bash

if [[ "$TERM" == "xterm-kitty" ]]; then
  onefetch --image assets/Tatake.png --image-protocol kitty
elif [[ "$TERM_PROGRAM" == "iTerm.app" ]]; then
  onefetch --image assets/Tatake.png --image-protocol iterm
elif [[ "$TERM" == *sixel* ]]; then
  onefetch --image assets/Tatake.png --image-protocol sixel
else
  onefetch
fi
