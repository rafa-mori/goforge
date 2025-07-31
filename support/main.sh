#!/usr/bin/env bash
# shellcheck disable=SC2015

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$( date +%Y-%m-%d )"
__secure_logic_author="Rafael Mori"
__secure_logic_use_type="exec"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ "${MYNAME_VERBOSE:-false}" == "true" ]]; then
  set -x  # Enable debugging
fi

IFS=$'\n\t'

__secure_logic_sourced_name() {
  local _self="${BASH_SOURCE-}"
  _self="${_self//${_kbx_root:-$()}/}"
  _self="${_self//\.sh/}"
  _self="${_self//\-/_}"
  _self="${_self//\//_}"
  echo "_was_sourced_${_self//__/_}"
  return 0
}

__first(){
  if [ "$EUID" -eq 0 ] || [ "$UID" -eq 0 ]; then
    echo "Please do not run as root." 1>&2 > /dev/tty
    exit 1
  elif [ -n "${SUDO_USER:-}" ]; then
    echo "Please do not run as root, but with sudo privileges." 1>&2 > /dev/tty
    exit 1
  else
    # shellcheck disable=SC2155
    local _ws_name="$(__secure_logic_sourced_name)"

    if test "${BASH_SOURCE-}" != "${0}"; then
      if test $__secure_logic_use_type != "lib"; then
        echo "This script is not intended to be sourced." 1>&2 > /dev/tty
        echo "Please run it directly." 1>&2 > /dev/tty
        exit 1
      fi
      # If the script is sourced, we set the variable to true
      # and export it to the environment without changing
      # the shell options.
      export "${_ws_name}"="true"
    else
      if test $__secure_logic_use_type != "exec"; then
        echo "This script is not intended to be executed directly." 1>&2 > /dev/tty
        echo "Please source it instead." 1>&2 > /dev/tty
        exit 1
      fi
      # If the script is executed directly, we set the variable to false
      # and export it to the environment. We also set the shell options
      # to ensure a safe execution.
      export "${_ws_name}"="false"
      set -o errexit # Exit immediately if a command exits with a non-zero status
      set -o nounset # Treat unset variables as an error when substituting
      set -o pipefail # Return the exit status of the last command in the pipeline that failed
      set -o errtrace # If a command fails, the shell will exit immediately
      set -o functrace # If a function fails, the shell will exit immediately
      shopt -s inherit_errexit # Inherit the errexit option in functions
    fi
  fi
}

_DEBUG=${DEBUG:-false}
_HIDE_ABOUT=${HIDE_ABOUT:-false}
_SCRIPT_DIR="$(dirname "${0}")"

__first "$@" >/dev/tty || exit 1

# Load library files
_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
#shellcheck source=/dev/null disable=SC2065
test -z "${_BANNER:-}" && source "${_SCRIPT_DIR}/config.sh" || true
#shellcheck source=/dev/null disable=SC2065
test -z "$(declare -f log)" >/dev/null && source "${_SCRIPT_DIR}/utils.sh" || true
#shellcheck source=/dev/null disable=SC2065
test -z "$(declare -f what_platform)" >/dev/null && source "${_SCRIPT_DIR}/platform.sh" || true
#shellcheck source=/dev/null disable=SC2065
test -z "$(declare -f check_dependencies)" >/dev/null && source "${_SCRIPT_DIR}/validate.sh" || true
#shellcheck source=/dev/null disable=SC2065
test -z "$(declare -f detect_shell_rc)" >/dev/null && source "${_SCRIPT_DIR}/install_funcs.sh" || true
#shellcheck source=/dev/null disable=SC2065
test -z "$(declare -f build_binary)" >/dev/null && source "${_SCRIPT_DIR}/build.sh" || true
#shellcheck source=/dev/null disable=SC2065
test -z "$(declare -f show_banner)" >/dev/null && source "${_SCRIPT_DIR}/info.sh" || true

# Initialize traps
set_trap "$@"

clear_screen

__run_custom_scripts() {
  local _STAGE="${1:-post}"
  if test -d "${_SCRIPT_DIR}/${_STAGE}.d/"; then
    if ls -1A "${_SCRIPT_DIR}/${_STAGE}.d/"*.sh >/dev/null 2>&1; then
      log info "Running custom ${_STAGE} scripts..." true
      local _CUSTOM_SCRIPTS=()
      # shellcheck disable=SC2011
      _CUSTOM_SCRIPTS=( "$(ls -1A "${_SCRIPT_DIR}/${_STAGE}.d/"*.sh | xargs -I{} basename {} || true)" )
      for _CUSTOM_SCRIPT in "${_CUSTOM_SCRIPTS[@]}"; do
        if [[ -f "${_SCRIPT_DIR}/${_STAGE}.d/${_CUSTOM_SCRIPT}" ]]; then
          log info "Executing script: ${_CUSTOM_SCRIPT}" true
          # Ensure the script is executable
          if [[ ! -x "${_SCRIPT_DIR}/${_STAGE}.d/${_CUSTOM_SCRIPT}" ]]; then
            log info "Making script executable: ${_CUSTOM_SCRIPT}"
            chmod +x "${_SCRIPT_DIR}/${_STAGE}.d/${_CUSTOM_SCRIPT}" || {
              log error "Failed to make script executable: ${_CUSTOM_SCRIPT}"
              return 1
            }
          fi
          # Execute the script
          "${_SCRIPT_DIR}/${_STAGE}.d/${_CUSTOM_SCRIPT}" "$@" || {
            log error "Script execution failed: ${_CUSTOM_SCRIPT}" true
            return 1
          }
          log success "Script executed successfully: ${_CUSTOM_SCRIPT}" true
        else
          log warn "Script not found: ${_CUSTOM_SCRIPT}" true
          return 1
        fi
      done
    fi
  fi
  return 0
}

__main() {
  if ! what_platform; then
    log error "Platform could not be determined." true
    return 1
  fi

  _ARGS=( "$@" )
  local default_label='Auto detect'
  local arrArgs=( "${_ARGS[@]:0:$#}" )
  local PLATFORM_ARG
  PLATFORM_ARG=$(_get_os_from_args "${arrArgs[1]:-${_PLATFORM}}")
  local ARCH_ARG
  ARCH_ARG=$(_get_arch_arr_from_args "${arrArgs[2]:-${_ARCH}}")
  log notice "Command: ${arrArgs[0]:-}" true
  log notice "Platform: ${PLATFORM_ARG:-$default_label}" true
  log notice "Architecture: ${ARCH_ARG:-$default_label}" true
  case "${arrArgs[0]:-}" in
    force|FORCE|-f|-F)
      log notice "Force mode enabled." true
      _FORCE="y"
      export FORCE="y"
      ;;
  esac

  case "${arrArgs[0]:-}" in
    help|HELP|-h|-H)
      log info "Help:"
      echo "Usage: make {build|build-dev|install|build-docs|clean|test|help}"
      echo "Commands:"
      echo "  build    - Compiles the binary for the specified platform and architecture."
      echo "  install  - Installs the binary, either by downloading a pre-compiled version or building it locally."
      echo "  build-dev - Builds the binary in development mode (without compression)."
      echo "  build-docs - Builds the documentation for the project."
      echo "  test     - Runs the tests for the project."
      echo "  clean    - Cleans up build artifacts."
      echo "  help     - Displays this help message."
      exit 0
      ;;
    build-dev|BUILD-DEV|-bd|-BD)
      # validate_versions
      log info "Running build command in development mode..."
      build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" false || return 1
      ;;
    build|BUILD|-b|-B)
      # validate_versions
      log info "Running build command..."
      build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" || return 1
      ;;
    install|INSTALL|-i|-I)
      log info "Running install command..."
      log info "How do you want to install the binary? [d/b/c] (10 seconds to respond, default: cancel)" true
      log question "(d)ownload pre-compiled binary, (b)uild locally, (c)ancel" true
      local choice
      read -t 10 -r -n 1 -p "" choice || choice='c'
      echo ""  # Move to the next line after reading input
      choice="${choice,,}"  # Convert to lowercase
      if [[ $choice =~ [dD] ]]; then
          log info "Downloading pre-compiled binary..."
          install_from_release || {
            log error "Failed to download pre-compiled binary." true
            return 1
          }
      elif [[ $choice =~ [bB] ]]; then
          log info "Building locally..."
          validate_versions || return 1
          build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" || return 1
          if [[ -f "${_BINARY}" ]]; then
            install_binary || {
              log error "Failed to install the binary." true
              return 1
            }
          fi
      else
          log info "Installation cancelled." true
          return 0
      fi
      summary "${arrArgs[@]}" || return 1
      ;;
    clear|clean|CLEAN|-c|-C)
      log info "Running clean command..."
      clean_artifacts || return 1
      log success "Clean completed successfully."
      ;;
    *)
      log error "Invalid command: ${arrArgs[0]:-}" true
      echo "Usage: make {build|build-dev|install|build-docs|clean|test|help}"
      ;;
  esac
}

# Função para limpar artefatos de build
clean_artifacts() {
    log info "Cleaning up build artifacts..."
    local platforms=("windows" "darwin" "linux")
    local archs=("amd64" "386" "arm64")
    for platform in "${platforms[@]}"; do
        for arch in "${archs[@]}"; do
            local output_name
            output_name=$(printf '%s_%s_%s' "${_BINARY}" "${platform}" "${arch}")
            if [[ "${platform}" != "windows" ]]; then
                local compress_name="${output_name}.tar.gz"
            else
                output_name="${output_name}.exe"
                local compress_name="${_BINARY}_${platform}_${arch}.zip"
            fi
            rm -f "${output_name}" || true
            rm -f "${compress_name}" || true
        done
    done
    log success "Build artifacts removed."
}

__secure_logic_main() {
  local _ws_name
  _ws_name="$(__secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval "echo \${$_ws_name}")
  if test "${_ws_name_val}" != "true"; then
    __main "$@"
    return $?
  else
    # If the script is sourced, we export the functions
    log error "This script is not intended to be sourced." true
    log error "Please run it directly." true
    return 1
  fi
}

if [[ "${_DEBUG}" != true ]]; then
  show_headers || log fatal "Failed to display headers." true
  if [[ -z "${_HIDE_ABOUT}" ]]; then
    show_about || log fatal "Failed to display about information." true
  fi
else
  log info "Debug mode enabled; banner will be ignored..."
  if [[ -z "${_HIDE_ABOUT}" ]]; then
    show_about || log fatal "Failed to display about information." true
  fi
fi

__run_custom_scripts "pre" "$@" || log fatal "Failed to execute pre-installation scripts."

__secure_logic_main "$@"

__run_custom_scripts "post" "$@" || log fatal "Failed to execute post-installation scripts."

__secure_logic_elapsed_time="$(($(date +%s) - __secure_logic_init_timestamp))"

if [[ "${MYNAME_VERBOSE:-false}" == "true" || "${_DEBUG:-false}" == "true" ]]; then
  log info "Script executed in ${__secure_logic_elapsed_time} seconds."
fi

# End of script logic