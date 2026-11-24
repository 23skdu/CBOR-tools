package cli

import "fmt"

// CompletionScript returns a shell completion script for the given shell.
// Supported shells are bash, zsh, and fish.
func CompletionScript(shell string) (string, error) {
	switch shell {
	case "bash":
		return bashCompletion, nil
	case "zsh":
		return zshCompletion, nil
	case "fish":
		return fishCompletion, nil
	}
	return "", fmt.Errorf("%w: %q (want bash, zsh, or fish)", ErrUnknownShell, shell)
}

// bashCompletion is installed as cbor-tool.bash.
const bashCompletion = `# bash completion for cbor-tool(1)

_cbortool() {
    local cur prev flags
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    flags="-a -f -s -t -o -format -root -completion -man -q -quiet -version -v"

    case "${prev}" in
        -a) COMPREPLY=( $(compgen -W "Encode Decode Tag" -- "${cur}") ); return 0 ;;
        -format) COMPREPLY=( $(compgen -W "pretty json" -- "${cur}") ); return 0 ;;
        -completion) COMPREPLY=( $(compgen -W "bash zsh fish" -- "${cur}") ); return 0 ;;
        -f|-o|-root) COMPREPLY=( $(compgen -f -- "${cur}") ); return 0 ;;
    esac

    if [[ "${cur}" == -* ]]; then
        COMPREPLY=( $(compgen -W "${flags}" -- "${cur}") )
    else
        COMPREPLY=( $(compgen -f -- "${cur}") )
    fi
}

complete -F _cbortool cbor-tool
`

// zshCompletion is installed as _cbor-tool.
const zshCompletion = `#compdef cbor-tool

_cbortool() {
    _arguments \
        '(-a)'-a'[action]:action:(Encode Decode Tag)' \
        '(-f)'-f'[input file]:input file:_files' \
        '(-s)'-s'[input string]:text:' \
        '(-t)'-t'[tag number]:tag number:' \
        '(-o)'-o'[output file]:output file:_files' \
        '--format[output format]:format:(pretty json)' \
        '--root[confine file access]:directory:_files -/' \
        '--completion[print a completion script]:shell:(bash zsh fish)' \
        '--man[print the manual page]' \
        '(-q)'-q'[suppress diagnostics]' \
        '(-q)'-quiet'[suppress diagnostics]' \
        '(-v)'-v'[print version]' \
        '(-v)'-version'[print version]'
}

_cbortool "$@"
`

// fishCompletion is installed as cbor-tool.fish.
const fishCompletion = `# fish completion for cbor-tool(1)

complete -c cbor-tool -f
complete -c cbor-tool -s a -r -a 'Encode Decode Tag' -d 'action to perform'
complete -c cbor-tool -s f -r -F -d 'input file path'
complete -c cbor-tool -s s -r -d 'input string'
complete -c cbor-tool -s t -r -d 'CBOR tag number'
complete -c cbor-tool -s o -r -F -d 'output file path'
complete -c cbor-tool -l format -r -a 'pretty json' -d 'output format for Decode'
complete -c cbor-tool -l root -r -F -d 'confine file access to this directory'
complete -c cbor-tool -l completion -r -a 'bash zsh fish' -d 'print a shell completion script'
complete -c cbor-tool -l man -d 'print the manual page and exit'
complete -c cbor-tool -s q -l quiet -d 'suppress diagnostics on stderr'
complete -c cbor-tool -s v -l version -d 'print version and exit'
`
