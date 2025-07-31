#!/usr/bin/env bash
# lib/validate.sh – Validação da versão do Go e dependências

validate_versions() {
    local _GO_SETUP='https://raw.githubusercontent.com/rafa-mori/gosetup/refs/heads/main/go.sh'
    local go_version
    go_version=$(go version | awk '{print $3}' | tr -d 'go' || echo "")
    if [[ -z "$go_version" ]]; then
        log error "Go is not installed or not found in PATH."
        return 1
    fi
    local version_target=""
    version_target="$(grep '^go ' go.mod | awk '{print $2}')"
    if [[ -z "$version_target" ]]; then
        log error "Could not determine the target Go version from go.mod."
        return 1
    fi
    if [[ "$go_version" != "$version_target" ]]; then
      local _go_installation_output
      if [[ -t 0 ]]; then
        _go_installation_output="$(bash -c "$(curl -sSfL "${_GO_SETUP}")" -s --version "$version_target" >/dev/tty)"
      else
        _go_installation_output="$(export NON_INTERACTIVE=true; bash -c "$(curl -sSfL "${_GO_SETUP}")" -s --version "$version_target")"
      fi
      if [[ $? -ne 0 ]]; then
          log error "Failed to install Go version ${version_target}. Output: ${_go_installation_output}"
          return 1
      fi
    fi
    # Check for required dependencies
    check_dependencies "git" "curl" "jq" "upx" "go" || return 1
    return 0
}

check_dependencies() {
  for dep in "$@"; do
    if ! command -v "$dep" > /dev/null; then
      log error "$dep is not installed." true
      return 1
    fi
  done
}

export -f validate_versions
export -f check_dependencies
