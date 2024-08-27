package scripts


const TmuxSessionizer = `#!/usr/bin/env bash
if [[ $# -eq 1 ]]; then
    selected=$1
else
    selected=$(find ~/dev -maxdepth 1 -type d -exec find {} -maxdepth 1 -type d \; | sed '1d' | fzf)
fi

if [[ -z $selected ]]; then
    exit 0
fi

selected_name=$(basename "$selected" | tr . _)
tmux_running=$(pgrep tmux)

if [[ -z $TMUX ]] && [[ -z $tmux_running ]]; then
    tmux new-session -s $selected_name -c $selected
    exit 0
fi

if ! tmux has-session -t=$selected_name 2> /dev/null; then
    tmux new-session -ds $selected_name -c $selected
fi

tmux switch-client -t $selected_name
`

const TmuxConfig = `set -ga terminal-overrides ",screen-256color*:Tc"
set-option -g default-terminal "screen-256color"
set -g status-style 'bg=#333333 fg=#5eacd3'

# Start windows and panes at 1, not 0
set -g base-index 1
setw -g pane-base-index 1

bind-key -r f run-shell "tmux neww ~/.local/bin/tmux-sessionizer"

set-window-option -g mode-keys vi
bind -T copy-mode-vi v send-keys -X begin-selection
bind -T copy-mode-vi y send-keys -X copy-pipe-and-cancel 'xclip -in -selection clipboard'

# vim-like pane switching
bind -r j select-pane -U
bind -r k select-pane -D
bind -r l select-pane -L
bind -r ';' select-pane -R

bind-key -r j run-shell "~/.local/bin/tmux-sessionizer ~/dev/databyte/api"
bind-key -r k run-shell "~/.local/bin/tmux-sessionizer ~/dev/databyte/scrapers"
bind-key -r l run-shell "~/.local/bin/tmux-sessionizer ~/dev/databyte/ui"
bind-key -r ';' run-shell "~/.local/bin/tmux-sessionizer ~/dev/lazydocker"
`
