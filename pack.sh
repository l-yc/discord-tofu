#!/bin/sh
zip -r tofu.zip \
  discord-tofu \
  config.toml \
  pics/assets/ \
  answer/autorespond/message-autoresponder/*.{py,pickle}
