#compdef ldss

_ldss()
{
	cur=${COMP_WORDS[COMP_CWORD]}
	echo $cur
	COMPREPLY=( $( compgen -W "$use" -- $cur ) )
}
complete -o default -o nospace -F _ldss ldss
