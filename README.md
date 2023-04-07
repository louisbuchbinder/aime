# Aime
A simple tool for working with LLMs. Currently only OpenAI api.

### Install
`go install github.com/louisbuchbinder/aime/aime@latest`

### Basic usage
`aime -key $OPENAI_API_KEY -output inline -prompt "1 + 1 = "

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
