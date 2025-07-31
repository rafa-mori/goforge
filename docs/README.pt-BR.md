# 🚀 GoForge: Automação, CLI Moderna e Estrutura Profissional para Módulos Go

[![Build](https://github.com/rafa-mori/goforge/actions/workflows/release.yml/badge.svg)](https://github.com/rafa-mori/goforge/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.20-blue)](go.mod)
[![Releases](https://img.shields.io/github/v/release/rafa-mori/goforge?include_prereleases)](https://github.com/rafa-mori/goforge/releases)

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
- **Configuração centralizada** via `info/manifest.json` com acesso embarcado
- **Artefatos de build organizados** em diretório dedicado `bin/`
- **Checksum automático** para garantir integridade
- **CLI customizada e estilizada** (cobra), pronta para ser estendida
- **Arquitetura flexível**: use como biblioteca ou executável
- **Versionamento automático**: CI/CD preenche e embeda a versão no binário
- **Logger estruturado**: logging contextual, colorido, com níveis e rastreio de linha

Tudo isso sem precisar alterar o código do seu módulo individualmente. O workflow é modular, dinâmico e se adapta ao ambiente!

---

## 🏗️ Estrutura do Projeto

```text
./
├── .github/workflows/      # Workflows de CI/CD (release, checksum)
├── goforge.go              # Interface GoForge para uso como lib
├── cmd/                    # Entrypoint e comandos da CLI
│   ├── cli/                # Utilitários e comandos de exemplo
│   ├── main.go             # Main da aplicação CLI
│   ├── usage.go            # Template de usage customizado
│   └── wrpr.go             # Estrutura e registro de comandos
├── go.mod                  # Dependências Go
├── info/                   # Metadados da aplicação e manifest
│   ├── manifest.json       # Configuração central da aplicação
│   └── application.go      # Interface Go para acesso ao manifest
├── logger/                 # Logger global estruturado
│   └── logger.go           # Logger contextual e colorido
├── Makefile                # Entrypoint para build, test, lint, etc.
├── bin/                    # Diretório de artefatos de build (criado durante build)
├── support/                # Scripts auxiliares para build/install
└── version/                # Versionamento automático
    ├── CLI_VERSION         # Preenchido pelo CI/CD (depreciado)
    └── semantic.go         # Utilitários de versionamento semântico
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

O binário será gerado no diretório `bin/` como `bin/goforge`.

### 3. Rode a CLI

```sh
./bin/goforge --help
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

O arquivo `info/manifest.json` contém a versão da aplicação e metadados. O sistema de versão se integra com o GitHub para verificar atualizações. O comando `goforge version` mostra a versão atual e a última versão disponível no GitHub.

---

## 📦 Sistema de Distribuição e Build

O GoForge possui um sistema sofisticado de distribuição com configuração centralizada através do `info/manifest.json`. Para informações detalhadas sobre o processo de build, versionamento e integração CI/CD, veja:

**[📋 Documentação de Distribuição e Versionamento](DISTRIBUTION.pt-BR.md)**

Isso cobre:

- Arquitetura baseada em manifest
- Processo de build multi-plataforma
- Otimização e organização de binários
- Automação CI/CD
- Sistema de gerenciamento de versões

---

## 🤝 Contribua

Pull requests, issues e sugestões são super bem-vindos. Vamos evoluir juntos!

---

## 📄 Licença

MIT. Veja o arquivo LICENSE.

---

## 👤 Autor

Rafael Mori — [@rafa-mori](https://github.com/rafa-mori)

---

## 🌐 Links

- [Repositório no GitHub](https://github.com/rafa-mori/goforge)
- [Documentação de Distribuição e Versionamento](DISTRIBUTION.pt-BR.md)
- [Exemplo de uso do logger](../logger/logger.go)
- [Workflows de CI/CD](../.github/workflows/)

---

> Feito com 💙 para a comunidade Go. Bora automatizar tudo!
