#!/bin/bash
wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -O yt-dlp
chmod +x yt-dlp
yt-dlp --help
# need to add a line to access the command from everywhere
