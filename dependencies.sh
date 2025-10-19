#!/bin/bash
wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -O yt-dlp
chmod +x yt-dlp
current_dir= $PWD
# yt-dlp --help
export PATH="$PATH:$current_dir"
# need to add a line to access the command from everywhere
