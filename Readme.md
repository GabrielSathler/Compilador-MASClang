# Compilador em Go

## Gramática

Para definir a gramática fizemos uma análise nas principais declarações da linguagem, que resultou em:

### Tipos presentes na linguagem
- "int" | "float" | "char" | "bool" | "string"

### Declaração de variáveis
- "var" IDENT ":" tipo ("=" expressao)? ";"

### Entrada e saída
- "input" "(" IDENT ")" ";"
- "print" "(" expressao ")" ";"

### Estruturas Condicionais
- "if" "(" expressao ")" bloco ("else" bloco)?

### Laços de Repetição
- "for" "(" declaracao | atribuicao ";" expressao ";" atribuicao ")" bloco
- "while" "(" expressao ")" bloco

### Funções
- "func" IDENT "(" parametros? ")" ":" tipo bloco

### Bloco de código
- "{" (comando)* "}"

### Retorno
- "return" expressao ";"

### Parâmetros de função
- IDENT ":" tipo

## Arquitetura

Para a arquitetura do projeto decidimos seguir como um "orientado por pacotes", onde cada pacote contém structs principais do projeto, como: AST (Árvore de Sintaxe Abstrata), analisador léxico, os tokens da linguagem, analisador sintático (parser) e analisador semântico.
Cada pacote é responsável por realizar apenas as tarefas designadas a sua respecitva estrutura no compilador. 

## Passo a passo para uso

Para rodar corretamente o programa é necessário adicionar o Golang na máquina. É possível baixar seguindo os passos da documentação oficial em:
- https://go.dev/doc/install

*Obs: a versão mínima para rodar corretamente o projeto é a* `1.24.1`.

Com o Golang instalado corretamente, adicione o código para compilação no arquivo `input.test` e rode o comando `go run main.go` na raiz do projeto.
