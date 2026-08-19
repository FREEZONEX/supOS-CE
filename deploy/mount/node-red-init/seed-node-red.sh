#!/bin/sh
set -eu

mode="full"
seed_name=""

usage() {
  cat >&2 <<'USAGE'
usage: seed-node-red.sh [--quick|--packages|--full] [seed-name]

Modes:
  --quick     Prepare only lightweight startup assets.
  --packages  Install offline modules, npm cache, and manifest packages.
  --full      Run quick setup and package installation. This is the default.
USAGE
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --quick)
      mode="quick"
      shift
      ;;
    --packages)
      mode="packages"
      shift
      ;;
    --full)
      mode="full"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --*)
      usage
      exit 2
      ;;
    *)
      if [ -n "$seed_name" ]; then
        usage
        exit 2
      fi
      seed_name="$1"
      shift
      ;;
  esac
done

seed_name="${seed_name:-${NODE_RED_SEED_NAME:-node-red}}"
data_dir="${NODE_RED_USER_DIR:-/data}"
seed_dir="${NODE_RED_SEED_DIR:-/opt/tier0-mount/${seed_name}}"
version="${NODE_RED_SEED_VERSION:-}"

if [ -z "$version" ] && [ -f "${seed_dir}/.seed-version" ]; then
  version="$(cat "${seed_dir}/.seed-version")"
fi
if [ -z "$version" ]; then
  version="dev"
fi

marker="${data_dir}/.seed-version"
offline_marker="${data_dir}/.seed-offline-modules-version"
package_marker="${data_dir}/.seed-packages-version"
package_status="${data_dir}/.seed-packages-status"
package_lock="${data_dir}/.seed-packages.lock"
current=""
if [ -f "$marker" ]; then
  current="$(cat "$marker")"
fi

mkdir -p "$data_dir" "${data_dir}/cache"
export npm_config_cache="${NPM_CONFIG_CACHE:-${data_dir}/.npm}"
mkdir -p "$npm_config_cache"
cache_marker="${npm_config_cache}/.seed-cache-version"
skip_package_install="false"
if [ -f "${data_dir}/.skip-package-install" ]; then
  skip_package_install="true"
fi

package_lock_acquired="false"

log() {
  printf '[seed-node-red:%s] %s\n' "$seed_name" "$*"
}

step() {
  printf '[seed-node-red:%s] [%s/%s] %s\n' "$seed_name" "$1" "$2" "$3"
}

detail() {
  printf '[seed-node-red:%s]   - %s\n' "$seed_name" "$*"
}

write_package_status() {
  printf '%s %s %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$1" "${2:-}" > "$package_status"
}

cleanup_package_lock() {
  code="$?"
  if [ "$package_lock_acquired" = "true" ]; then
    if [ "$code" -eq 0 ]; then
      write_package_status "success" "version=${version}"
    else
      write_package_status "failed" "version=${version} exit=${code}"
    fi
    rm -rf "$package_lock"
  fi
  exit "$code"
}

acquire_package_lock() {
  if mkdir "$package_lock" 2>/dev/null; then
    package_lock_acquired="true"
    trap cleanup_package_lock EXIT
    write_package_status "running" "version=${version}"
    return 0
  fi
  detail "Package installation is already running: ${package_lock}"
  return 1
}

packages_current() {
  installed_version=""
  if [ -f "$package_marker" ]; then
    installed_version="$(cat "$package_marker")"
  fi
  [ "$installed_version" = "$version" ] && [ -d "${data_dir}/node_modules" ]
}

ensure_package_json() {
  step 1 2 "Ensuring Node-RED project metadata"
  if [ -f "${data_dir}/package.json" ]; then
    detail "package.json already exists"
    return
  fi
  cat > "${data_dir}/package.json" <<'JSON'
{
  "name": "node-red-project",
  "description": "A Node-RED Project",
  "version": "0.0.1",
  "private": true
}
JSON
  detail "Created ${data_dir}/package.json"
}

sync_managed_seed_assets() {
  step 2 2 "Syncing managed seed assets"
  if [ "$current" = "$version" ]; then
    detail "Managed seed marker already at version ${version}"
    if [ -d "${seed_dir}/themes" ]; then
      mkdir -p "${data_dir}/themes"
      cp -R "${seed_dir}/themes/." "${data_dir}/themes/"
      detail "Refreshed managed themes"
    fi
    return
  fi

  detail "Refreshing seed assets: ${current:-none} -> ${version}"
  if [ -f "${seed_dir}/settings.js" ]; then
    cp "${seed_dir}/settings.js" "${data_dir}/settings.js"
    detail "Copied settings.js"
  fi

  if [ -d "${seed_dir}/themes" ]; then
    mkdir -p "${data_dir}/themes"
    cp -R "${seed_dir}/themes/." "${data_dir}/themes/"
    detail "Copied themes"
  fi

  if [ ! -f "${data_dir}/flows.json" ]; then
    if [ -f "${seed_dir}/flows.json" ]; then
      cp "${seed_dir}/flows.json" "${data_dir}/flows.json"
      detail "Copied flows.json"
    else
      printf '[]\n' > "${data_dir}/flows.json"
      detail "Created empty flows.json"
    fi
  fi

  printf '%s\n' "$version" > "$marker"
  detail "Wrote seed marker ${marker}"
}

prepare_offline_modules() {
  step 1 3 "Preparing offline module assets"
  offline_version=""
  if [ -f "$offline_marker" ]; then
    offline_version="$(cat "$offline_marker")"
  fi
  if [ "$offline_version" = "$version" ] && [ -d "${data_dir}/offline_modules" ]; then
    detail "Offline module assets already at version ${version}"
    return
  fi

  if [ -f "${seed_dir}/offline_modules.tar.gz" ]; then
    detail "Extracting offline_modules.tar.gz"
    rm -rf "${data_dir}/offline_modules"
    tar -xzf "${seed_dir}/offline_modules.tar.gz" -C "${data_dir}"
  elif [ -f "${seed_dir}/offline_modules.tar.xz" ]; then
    detail "Extracting offline_modules.tar.xz"
    rm -rf "${data_dir}/offline_modules"
    tar -xf "${seed_dir}/offline_modules.tar.xz" -C "${data_dir}"
  elif [ -d "${seed_dir}/offline_modules" ]; then
    detail "Copying offline_modules directory"
    rm -rf "${data_dir}/offline_modules"
    mkdir -p "${data_dir}/offline_modules"
    cp -R "${seed_dir}/offline_modules/." "${data_dir}/offline_modules/"
  else
    detail "No offline module assets found"
    return
  fi

  printf '%s\n' "$version" > "$offline_marker"
  detail "Wrote offline module marker ${offline_marker}"
}

prepare_npm_cache() {
  step 2 3 "Preparing npm cache"
  cache_version=""
  if [ -f "$cache_marker" ]; then
    cache_version="$(cat "$cache_marker")"
  fi
  if [ -f "${seed_dir}/npmCache.tar.xz" ] && { [ "$cache_version" != "$version" ] || [ ! -d "${npm_config_cache}/_cacache" ]; }; then
    detail "Extracting npmCache.tar.xz: ${cache_version:-none} -> ${version}"
    rm -rf "${npm_config_cache}/_cacache" "${npm_config_cache}/_logs" "${npm_config_cache}/_update-notifier-last-checked"
    tar -xf "${seed_dir}/npmCache.tar.xz" -C "${data_dir}"
    printf '%s\n' "$version" > "$cache_marker"
    detail "Wrote npm cache marker ${cache_marker}"
  else
    detail "npm cache already prepared"
  fi
}

remove_obsolete_packages() {
  obsolete_manifest="${seed_dir}/obsolete-packages.txt"
  if [ ! -f "$obsolete_manifest" ]; then
    return
  fi

  obsolete_specs=""
  while IFS= read -r line || [ -n "$line" ]; do
    line="${line%%#*}"
    line="$(printf '%s' "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')"
    if [ -n "$line" ]; then
      obsolete_specs="${obsolete_specs} ${line}"
    fi
  done < "$obsolete_manifest"

  if [ -z "$obsolete_specs" ]; then
    return
  fi

  detail "Removing obsolete package specs"
  # shellcheck disable=SC2086
  (cd "$data_dir" && npm uninstall --no-audit --no-fund ${obsolete_specs})
}

install_manifest_packages() {
  step 3 3 "Installing manifest packages"
  manifest="${seed_dir}/packages.txt"
  if [ ! -f "$manifest" ]; then
    detail "No packages.txt found"
    return
  fi

  specs=""
  local_specs=""
  while IFS= read -r line || [ -n "$line" ]; do
    line="${line%%#*}"
    line="$(printf '%s' "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')"
    if [ -z "$line" ]; then
      continue
    fi
    case "$line" in
      ./*|../*|/*)
        local_specs="${local_specs} ${line}"
        ;;
      *)
        specs="${specs} ${line}"
        ;;
    esac
  done < "$manifest"

  remove_obsolete_packages

  if [ -z "$specs" ] && [ -z "$local_specs" ]; then
    printf '%s\n' "$version" > "$package_marker"
    detail "Package manifest is empty"
    return
  fi

  if [ -n "$local_specs" ]; then
    detail "Installing local package specs"
    # shellcheck disable=SC2086
    (cd "$data_dir" && npm install --no-audit --no-fund --unsafe-perm ${local_specs})
  fi

  if [ -n "$specs" ]; then
    offline_flag=""
    cache_version=""
    if [ -f "$cache_marker" ]; then
      cache_version="$(cat "$cache_marker")"
    fi
    if [ -f "${seed_dir}/npmCache.tar.xz" ] && [ "$cache_version" = "$version" ] && [ -d "${npm_config_cache}/_cacache" ]; then
      offline_flag="--offline"
    fi
    detail "Installing npm package specs${offline_flag:+ with offline cache}"
    # shellcheck disable=SC2086
    (cd "$data_dir" && npm install --no-audit --no-fund ${offline_flag} ${specs})
  fi

  printf '%s\n' "$version" > "$package_marker"
  detail "Wrote package marker ${package_marker}"
}

run_quick_seed() {
  log "Starting quick seed initialization"
  detail "Version: ${version}"
  detail "Data dir: ${data_dir}"
  detail "Seed dir: ${seed_dir}"
  ensure_package_json
  sync_managed_seed_assets
  log "Quick seed initialization complete"
}

run_package_seed() {
  log "Starting package seed initialization"
  detail "Version: ${version}"
  detail "Data dir: ${data_dir}"
  detail "Seed dir: ${seed_dir}"
  if [ "$skip_package_install" = "true" ]; then
    detail "Skipping package installation: ${data_dir}/.skip-package-install exists"
    write_package_status "skipped" "version=${version}"
    return
  fi
  ensure_package_json
  if ! acquire_package_lock; then
    return
  fi
  if packages_current; then
    detail "Manifest packages already installed for version ${version}"
    return
  fi
  prepare_offline_modules
  prepare_npm_cache
  install_manifest_packages
  log "Package seed initialization complete"
}

case "$mode" in
  quick)
    run_quick_seed
    ;;
  packages)
    run_package_seed
    ;;
  full)
    run_quick_seed
    run_package_seed
    log "Seed initialization complete"
    ;;
  *)
    usage
    exit 2
    ;;
esac
