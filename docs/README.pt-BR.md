# 🚀 GoForge: Automação, CLI Moderna e Estrutura Profissional para Módulos Go

[![Build](https://github.com/rafa-mori/goforge/actions/workflows/release.yml/badge.svg)](https://github.com/rafa-mori/goforge/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.20-blue)](go.mod)
[![Releases](https://img.shields.io/github/v/release/faelmori/goforge?include_prereleases)](https://github.com/rafa-mori/goforge/releases)

Se você já cansou de builds manuais, deploys complicados, versionamento confuso e quer uma CLI estilosa, fácil de estender e pronta para produção, o **GoForge** é pra você!

---

## 🌟 Exemplos Avançados

### 1. Estendendo a CLI com um novo comando

Crie um novo arquivo em `cmd/cli/hello.go`:

```go
package cli

import (
    "fmt"
    "github.com/spf13/cobra"
)

var HelloCmd = &cobra.Command{
    Use:   "hello",
    Short: "Exemplo de comando customizado",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Olá, mundo! Comando customizado funcionando!")
    },
}
```

No `wrpr.go`, registre o comando:

```go
// ...existing code...
rootCmd.AddCommand(cli.HelloCmd)
// ...existing code...
```

---

### 2. Logger avançado com contexto extra

```go
import gl "github.com/rafa-mori/goforge/logger"

func exemploComContexto() {
    gl.Log("warn", "Atenção! Algo pode estar errado.")
    gl.Log("debug", map[string]interface{}{
        "user": "rafael",
        "action": "login",
        "success": true,
    })
}
```

---

### 3. Usando como biblioteca Go

```go
import "github.com/rafa-mori/goforge"

func main() {
    var myModule goforge.GoForge = &MeuModulo{}
    if myModule.Active() {
        _ = myModule.Execute()
    }
}

// Implemente a interface GoForge no seu módulo
```

---

## ✨ O que é o GoForge?

O GoForge é um template/projeto base para qualquer módulo Go moderno. Ele entrega:

- **Build multi-plataforma** (Linux, macOS, Windows) sem mexer no código
- **Compactação UPX** automática para binários otimizados
- **Publicação automática** no GitHub Releases
- **Gerenciamento de dependências** unificado
- **Checksum automático** para garantir integridade
- **CLI customizada e estilizada** (cobra), pronta para ser estendida
- **Arquitetura flexível**: use como biblioteca ou executável
- **Versionamento automático**: CI/CD preenche e embeda a versão no binário
- **Logger estruturado**: logging contextual, colorido, com níveis e rastreio de linha

Tudo isso sem precisar alterar o código do seu módulo individualmente. O workflow é modular, dinâmico e se adapta ao ambiente!

---

## 🏗️ Estrutura do Projeto

```plain text
./
├── .github/workflows/      # Workflows de CI/CD (release, checksum)
├── article.go              # Interface GoForge para uso como lib
├── cmd/                    # Entrypoint e comandos da CLI
│   ├── cli/                # Utilitários e comandos de exemplo
│   ├── main.go             # Main da aplicação CLI
│   ├── usage.go            # Template de usage customizado
│   └── wrpr.go             # Estrutura e registro de comandos
├── go.mod                  # Dependências Go
├── logger/                 # Logger global estruturado
│   └── logger.go           # Logger contextual e colorido
├── Makefile                # Entrypoint para build, test, lint, etc.
├── support/                # Scripts auxiliares para build/install
├── version/                # Versionamento automático
│   ├── CLI_VERSION         # Preenchido pelo CI/CD
│   └── semantic.go         # Utilitários de versionamento semântico
```

---

## 💡 Por que usar?

- **Zero dor de cabeça** com builds e deploys
- **CLI pronta para produção** e fácil de customizar
- **Logger poderoso**: debug, info, warn, error, success, tudo com contexto
- **Versionamento automático**: nunca mais esqueça de atualizar a versão
- **Fácil de estender**: adicione comandos, use como lib, plugue em outros projetos

---

## 🚀 Como usar

### 1. Instale as dependências

```sh
make install
```

### 2. Build do projeto

```sh
make build
```

### 3. Rode a CLI

```sh
./goforge --help
```

### 4. Adicione comandos customizados

Crie arquivos em `cmd/cli/` e registre no `wrpr.go`.

---

## 🛠️ Exemplo de uso do Logger

```go
import gl "github.com/rafa-mori/goforge/logger"

gl.Log("info", "Mensagem informativa")
gl.Log("error", "Algo deu errado!")
```

O logger já inclui contexto (linha, arquivo, função) automaticamente!

---

## 🔄 Versionamento automático

O arquivo `version/CLI_VERSION` é preenchido pelo CI/CD a cada release/tag. O comando `goforge version` mostra a versão atual e a última disponível no GitHub.

---

## 🤝 Contribua

Pull requests, issues e sugestões são super bem-vindos. Vamos evoluir juntos!

---

## 📄 Licença

MIT. Veja o arquivo LICENSE.

---

## 👤 Autor

Rafael Mori — [@faelmori](https://github.com/rafa-mori)

---

## 🌐 Links

- [Repositório no GitHub](https://github.com/rafa-mori/goforge)
- [Exemplo de uso do logger](logger/logger.go)
- [Workflows de CI/CD](.github/workflows/)

---

> Feito com 💙 para a comunidade Go. Bora automatizar tudo!