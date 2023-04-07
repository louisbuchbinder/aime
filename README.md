# Aime
A simple tool for working with LLMs. Currently only OpenAI api.

### Install
`go install github.com/louisbuchbinder/aime/cmd/aime@latest`

### Basic usage
`aime -key $OPENAI_API_KEY -output inline -prompt "1 + 1 = "

### System prompts
Directly specify the system prompt on the command line with `-systemPrompt`.

Or reference a system prompt using the `systemPromptKey`.

Some system prompts are builtin, see `prompt.go`

You can configure custom system prompts using the config, described below.

### Config
The aime config can be loaded from one of the following file locations in order
1. `.aimerc.json`
2. `.aimerc.yaml`
3. `.aimerc.yml`
4. `~/.aimerc.json`
5. `~/.aimerc.yaml`
6. `~/.aimerc.yml`

An example config might look like (.aimerc.yml)
```
system:
  baby: talk to a baby
  scientist: talk to a scientist
```

Then:
1. `aime -key $OPENAI_API_KEY -prompt "the earth is " -systemPromptKey baby -output inline`
2. `aime -key $OPENAI_API_KEY -prompt "the earth is " -systemPromptKey scientist -output inline`


### Vimrc
Add a shortcut to your .vimrc to call aime with a systemPromptKey of the current open file extension
```.vimrc
function GetCurrentFileExtension()
  let filename = expand('%')
  let extension = substitute(filename, '.*\.\(\w\+\)$', '\1', '')
  return extension
endfunction

function GetVisualSelection()
  let [line1, col1] = getpos("'<")[1:2]
  let [line2, col2] = getpos("'>")[1:2]
  let lines = getline(line1, line2)
  let lines[-1] = lines[-1][: col2 - (&selection == 'inclusive' ? 1 : 2)]
  let lines[0] = lines[0][col1 - 1:]
  return join(lines, "\n")
endfunction

function! Aime()
  " UPDATE THE KEYFILE PATH
  let keyfile = "/path/to/your/openai.key"
  let input = GetVisualSelection()
  let ext = GetCurrentFileExtension()
  let output = system('aime -output inline -keyfile '.keyfile.' -systemPromptKey '.ext, input)
  put=output
endfunction

noremap <C-w>a :call Aime()<CR>
```
