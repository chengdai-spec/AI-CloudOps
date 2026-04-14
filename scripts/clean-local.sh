#!/bin/bash

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

clean_dir() {
    local dir="$1"
    if [ -d "$dir" ]; then
        shopt -s dotglob nullglob
        local items=("$dir"/*)
        shopt -u dotglob nullglob
        if [ "${#items[@]}" -gt 0 ]; then
            rm -rf "${items[@]}"
        fi
    fi
}

clean_dir "bin"
clean_dir "tmp"
clean_dir "logs"

if [ "${1:-}" = "--data" ]; then
    clean_dir "data/mysql/data"
    clean_dir "data/redis/data"
    clean_dir "data/prometheus/data"
    clean_dir "data/nginx/log"
fi

echo "Cleaned: bin, tmp, logs"
if [ "${1:-}" = "--data" ]; then
    echo "Cleaned: data/mysql/data, data/redis/data, data/prometheus/data, data/nginx/log"
fi
