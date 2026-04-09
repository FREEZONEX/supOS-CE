#!/bin/bash

normalize_volume_shell_scripts() {
  local root="${1:-}"
  if [ -z "$root" ] || [ ! -d "$root" ]; then
    return 0
  fi

  local count=0
  while IFS= read -r -d '' file; do
    sed -i 's/\r$//' "$file"
    count=$((count + 1))
  done < <(find "$root" -type f -name "*.sh" -print0 2>/dev/null)

  if [ "$count" -gt 0 ]; then
    info "Normalized shell line endings under $root ($count files)."
  fi
}
