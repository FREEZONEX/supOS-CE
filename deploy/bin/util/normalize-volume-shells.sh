#!/bin/bash

normalize_volume_shell_scripts() {
  local root="${1:-}"
  if [ -z "$root" ] || [ ! -d "$root" ]; then
    return 0
  fi

  local count=0
  local failed=0
  while IFS= read -r -d '' file; do
    if sed -i 's/\r$//' "$file" 2>/dev/null; then
      count=$((count + 1))
    else
      failed=$((failed + 1))
      warn "Skipping line-ending normalization for $file"
    fi
  done < <(
    find "$root" \
      \( -path '*/node_modules/*' -o -path '*/.npm/*' \) -prune -o \
      -type f -name "*.sh" -print0 2>/dev/null
  )

  if [ "$count" -gt 0 ]; then
    info "Normalized shell line endings under $root ($count files)."
  fi
  if [ "$failed" -gt 0 ]; then
    warn "Skipped line-ending normalization for $failed files under $root."
  fi
}
