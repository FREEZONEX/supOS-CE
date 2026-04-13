#!/bin/bash

# This script handles the selection and validation of the ENTRANCE_DOMAIN (IP Address).
# 1. It provides an interactive UI to select the IP for different platforms.
# 2. It immediately validates if the selected IP is a loopback address.

info "Configuring network IP address (ENTRANCE_DOMAIN)..."
ENV_FILE="$SCRIPT_DIR/../.env.default"
if [ -f "$SCRIPT_DIR/../.env" ]; then
  ENV_FILE="$SCRIPT_DIR/../.env"
fi

# --- Part 1: IP Selection ---
non_interactive=false
if [ ! -t 0 ]; then
  non_interactive=true
fi

if [[ "$platform" == MINGW64* ]]; then
    # IP selection logic for Windows
    sed -i -e "s/^OS_PLATFORM_TYPE=.*/OS_PLATFORM_TYPE=windows/" "$ENV_FILE"
    current_entrance_domain=$(grep '^ENTRANCE_DOMAIN=' "$ENV_FILE" | cut -d '=' -f2- | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    sed -i "s/^ENTRANCE_DOMAIN=.*/ENTRANCE_DOMAIN=$current_entrance_domain/" "$ENV_FILE"

    if [[ "$non_interactive" == true ]]; then
        selected_ip="$current_entrance_domain"
        echo "Non-interactive mode detected. Keeping current ENTRANCE_DOMAIN: $selected_ip"
    elif [[ -n "$current_entrance_domain" ]]; then
        while true; do
            read -p "Choose IP address for ENTRANCE_DOMAIN (Press Enter for default: [$current_entrance_domain]): " selected_ip
            selected_ip=$(echo "$selected_ip" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//'); selected_ip=${selected_ip:-$current_entrance_domain}
            [[ -n "$selected_ip" ]] && break || echo "Input cannot be empty. Please enter a valid IP."
        done
    else
        while true; do
            read -p "Choose IP address for ENTRANCE_DOMAIN: " selected_ip
            selected_ip=$(echo "$selected_ip" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            [[ -n "$selected_ip" ]] && break || echo "Input cannot be empty. Please enter a valid IP."
        done
    fi
else
    # IP selection menu for Linux/macOS
    ips=($(hostname -I | awk '{print $1, $2, $3}'))
    echo -e "\nAvailable options for ENTRANCE_DOMAIN:"
    current_entrance_domain=$(grep '^ENTRANCE_DOMAIN=' "$ENV_FILE" | cut -d '=' -f2- | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    sed -i "s/^ENTRANCE_DOMAIN=.*/ENTRANCE_DOMAIN=$current_entrance_domain/" "$ENV_FILE"
    echo "0). Keep current: $current_entrance_domain (default)"
    for i in "${!ips[@]}"; do echo "$((i+1))). ${ips[$i]}"; done
    echo "$((${#ips[@]}+1))). Custom IP address (enter manually)"

    if [[ "$non_interactive" == true ]]; then
      selected_ip="$current_entrance_domain"
      echo "Non-interactive mode detected. Keeping current ENTRANCE_DOMAIN: $selected_ip"
    else
      while true; do
        read -p "Select option (0-$((${#ips[@]}+1))), or press Enter to keep current: " choice
        if [[ -z "$choice" ]] || [[ "$choice" == "0" ]]; then
            [[ -n "$current_entrance_domain" ]] || { echo "Current ENTRANCE_DOMAIN is empty. Please select a valid option."; continue; }
            selected_ip="$current_entrance_domain"; echo "Keeping current ENTRANCE_DOMAIN: $selected_ip"; break
        elif [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#ips[@]}" ]; then
          selected_ip=${ips[$((choice-1))]}; [[ -n "$selected_ip" ]] || { echo "Selected IP is empty. Please choose another option."; continue; }; echo "Selected IP: $selected_ip"; break
        elif [[ "$choice" == "$((${#ips[@]}+1))" ]]; then
          while true; do read -p "Enter custom IP address or domain: " custom_ip; custom_ip=$(echo "$custom_ip" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//'); [[ -n "$custom_ip" ]] && { selected_ip="$custom_ip"; echo "Selected custom IP: $selected_ip"; break; } || echo "Input cannot be empty."; done; break
        else echo "Invalid input. Please enter a valid option number (0-$((${#ips[@]}+1))) or press Enter."; fi
      done
    fi
fi

if [ "$selected_ip" != "$current_entrance_domain" ]; then
    escaped_selected_ip=$(sed 's/[&]/\\&/g' <<< "$selected_ip")
    sed -i "s|^ENTRANCE_DOMAIN=.*|ENTRANCE_DOMAIN=$escaped_selected_ip|" "$ENV_FILE"
    source "$ENV_FILE"
fi


# --- Part 2: Loopback Address Validation (Merged) ---
info "Validating selected IP address..."

# Use the '$selected_ip' variable that was just set
trimmed_ip=$(echo "$selected_ip" | xargs) # Trim potential whitespace

if [[ "$trimmed_ip" == "127.0.0.1" || "$trimmed_ip" == "localhost" ]]; then
  echo
  warn "Loopback address detected. Local IAM login remains available, but external SSO-style integrations may not work."
else
  info "IP address is valid."
fi
echo
