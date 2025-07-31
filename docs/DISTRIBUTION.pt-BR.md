# Fluxo de Distribuição e Versionamento

![GoForge Distribution](assets/top_banner_m_a.png)

---

[![Build](https://github.com/rafa-mori/goforge/actions/workflows/release.yml/badge.svg)](https://github.com/rafa-mori/goforge/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](../LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.20-blue)](../go.mod)
[![Releases](https://img.shields.io/github/v/release/rafa-mori/goforge?include_prereleases)](https://github.com/rafa-mori/goforge/releases)

---

[🇺🇸 Read this documentation in English](DISTRIBUTION.md) | [← Voltar ao README](README.pt-BR.md)

Este documento descreve o sistema avançado de distribuição e versionamento do GoForge, que fornece gerenciamento centralizado de informações da aplicação e geração automatizada de binários.

---

## 📋 Índice

- [Visão Geral](#-visão-geral)
- [Arquitetura Baseada em Manifest](#-arquitetura-baseada-em-manifest)
- [Package de Informações da Aplicação](#-package-de-informações-da-aplicação)
- [Geração de Binários](#️-geração-de-binários)
- [Sistema de Versionamento](#-sistema-de-versionamento)
- [Artefatos de Build](#-artefatos-de-build)
- [Integração CI/CD](#-integração-cicd)
- [Exemplos de Uso](#-exemplos-de-uso)

---

## 🎯 Visão Geral

O GoForge apresenta um sistema sofisticado de distribuição que centraliza todos os metadados da aplicação em uma única fonte de verdade: o arquivo `info/manifest.json`. Esta abordagem garante consistência entre builds, simplifica o gerenciamento de configuração e habilita poderosas capacidades de automação.

### Características Principais

- **Configuração Centralizada**: Um único arquivo `manifest.json` controla todos os metadados da aplicação
- **Manifest Embarcado**: Informações da aplicação são embarcadas em tempo de compilação para acesso em runtime
- **Versionamento Automatizado**: Versionamento semântico com integração ao GitHub
- **Builds Multi-plataforma**: Suporte para Linux, macOS e Windows
- **Otimização de Binários**: Compressão automática UPX para binários menores
- **Artefatos Organizados**: Todas as saídas de build vão para o diretório dedicado `bin/`

---

## 📄 Arquitetura Baseada em Manifest

### O Manifest Central

O arquivo `info/manifest.json` serve como a única fonte de verdade para todos os metadados da aplicação:

```json
{
  "name": "GoForge",
  "application": "goforge",
  "version": "1.0.1",
  "private": false,
  "published": true,
  "aliases": ["goforge"],
  "repository": "https://github.com/rafa-mori/goforge",
  "homepage": "https://github.com/rafa-mori/goforge",
  "description": "Uma ferramenta de scaffolding Go para criar projetos com estrutura padrão e melhores práticas.",
  "main": "cmd/main.go",
  "bin": "goforge",
  "author": "Rafael Mori <faelmori@gmail.com>",
  "organization": "rafa-mori",
  "license": "MIT",
  "keywords": ["goforge", "scaffolding", "tool", "management"],
  "platforms": ["linux"]
}
```

### Campos do Manifest

| Campo | Descrição | Exemplo |
|-------|-----------|---------|
| `name` | Nome de exibição da aplicação | `"GoForge"` |
| `application` | Nome interno da aplicação | `"goforge"` |
| `version` | Versão atual (semver) | `"1.0.1"` |
| `bin` | Nome do binário quando construído | `"goforge"` |
| `repository` | URL do repositório Git | `"https://github.com/..."` |
| `description` | Descrição curta | `"Uma ferramenta de scaffolding Go..."` |
| `main` | Arquivo ponto de entrada | `"cmd/main.go"` |
| `author` | Informações do autor | `"Nome <email>"` |
| `organization` | Organização/proprietário | `"rafa-mori"` |
| `platforms` | Plataformas suportadas | `["linux", "darwin", "windows"]` |
| `private` | Flag de repositório privado | `false` |
| `published` | Status de publicação | `true` |

---

## 🔧 Package de Informações da Aplicação

### Estrutura do Package

O package `info` fornece acesso type-safe aos dados do manifest:

```go
// info/application.go
package manifest

type Manifest interface {
    GetName() string
    GetVersion() string
    GetRepository() string
    GetDescription() string
    // ... outros getters
}
```

### Manifest Embarcado

O manifest é embarcado em tempo de compilação usando a diretiva `embed` do Go:

```go
//go:embed manifest.json
var manifestJSONData []byte
```

### Acesso em Runtime

Aplicações podem acessar dados do manifest em runtime:

```go
import manifest "github.com/rafa-mori/goforge/info"

func main() {
    info, err := manifest.GetManifest()
    if err != nil {
        log.Fatal("Falha ao obter manifest:", err)
    }
    
    fmt.Printf("Aplicação: %s v%s\n", 
        info.GetName(), 
        info.GetVersion())
}
```

---

## 🏗️ Geração de Binários

### Estrutura do Diretório de Build

Todos os artefatos de build são organizados no diretório `bin/`:

```text
bin/
├── goforge_linux_amd64           # Binário Linux
├── goforge_linux_amd64.tar.gz    # Arquivo comprimido
├── goforge_darwin_amd64          # Binário macOS
├── goforge_darwin_amd64.tar.gz   # Arquivo macOS
├── goforge_windows_amd64.exe     # Binário Windows
└── goforge_windows_amd64.zip     # Arquivo Windows
```

### Processo de Build

1. **Leitura do Manifest**: Scripts de build leem `info/manifest.json` para configuração
2. **Compilação do Binário**: Go build com dados do manifest embarcados
3. **Targeting de Plataforma**: Cross-compilation para múltiplas plataformas
4. **Otimização**: Compressão UPX (quando habilitada)
5. **Empacotamento**: Criação de arquivos com checksums

### Integração com Makefile

O Makefile automaticamente lê dados do manifest:

```makefile
APP_NAME := $(shell jq -r '.name' < $(ROOT_DIR)info/manifest.json)
BINARY_NAME := $(shell jq -r '.bin' < $(ROOT_DIR)info/manifest.json)
ORGANIZATION := $(shell jq -r '.organization' < $(ROOT_DIR)info/manifest.json)
```

---

## 🔄 Sistema de Versionamento

### Gerenciamento de Versão

O sistema de versionamento fornece:

- **Versão Atual**: Lida do `info/manifest.json`
- **Última Versão**: Obtida da API de releases do GitHub
- **Comparação de Versões**: Análise e comparação de versão semântica
- **Verificação de Atualizações**: Verificações automáticas para versões mais novas

### Interface do Serviço de Versão

```go
type Service interface {
    GetLatestVersion() (string, error)
    GetCurrentVersion() string
    IsLatestVersion() (bool, error)
    GetRepository() string
}
```

### Integração CLI

```bash
# Verificar versão atual
goforge version

# Verificar atualizações
goforge version --check

# Obter informações de versão
goforge version --info
```

---

## 📁 Artefatos de Build

### Desenvolvimento Local

```bash
# Build para plataforma atual
make build

# Instalar localmente
make install

# Limpar artefatos de build
make clean
```

### Build Multi-plataforma

```bash
# Build para todas as plataformas
make build-all

# Build para plataforma específica
make build PLATFORM=linux ARCH=amd64
```

### Tipos de Artefatos

- **Binários Raw**: `goforge_platform_arch[.exe]`
- **Arquivos Comprimidos**: `.tar.gz` (Unix) / `.zip` (Windows)
- **Checksums**: Hashes SHA256 para verificação de integridade
- **Símbolos de Debug**: Informações de debug separadas (quando habilitado)

---

## 🚀 Integração CI/CD

### Workflow do GitHub Actions

O workflow automatizado:

1. **Checkout do Código**: Obter código fonte e manifest
2. **Setup do Ambiente**: Configurar Go e ferramentas de build
3. **Leitura do Manifest**: Extrair configuração de build
4. **Build Multi-plataforma**: Gerar binários para todas as plataformas
5. **Otimização**: Aplicar compressão UPX
6. **Criação de Arquivos**: Empacotar binários com documentação
7. **Geração de Checksums**: Criar hashes de integridade
8. **Publicação de Release**: Upload para releases do GitHub

### Variáveis de Ambiente

| Variável | Descrição | Fonte |
|----------|-----------|-------|
| `APP_NAME` | Nome da aplicação | `manifest.json` |
| `VERSION` | Versão atual | `manifest.json` |
| `REPOSITORY` | URL do repositório | `manifest.json` |
| `BINARY_NAME` | Nome do binário de saída | `manifest.json` |

---

## 💡 Exemplos de Uso

### 1. Criando um Novo Release

1. Atualizar versão em `info/manifest.json`:

   ```json
   {
     "version": "1.0.2"
   }
   ```

1. Commit e tag:

   ```bash
   git add info/manifest.json
   git commit -m "bump: version 1.0.2"
   git tag v1.0.2
   git push origin main --tags
   ```

1. CI/CD automaticamente constrói e publica release

### 2. Adicionando Suporte a Plataforma

Atualizar o array de plataformas em `manifest.json`:

```json
{
  "platforms": ["linux", "darwin", "windows"]
}
```

### 3. Personalizando Saída de Build

Modificar o nome do binário:

```json
{
  "bin": "meu-nome-customizado"
}
```

Isso gera `meu-nome-customizado_linux_amd64` ao invés de `goforge_linux_amd64`.

### 4. Acesso ao Manifest em Runtime

```go
package main

import (
    "fmt"
    manifest "github.com/rafa-mori/goforge/info"
)

func main() {
    info, _ := manifest.GetManifest()
    
    fmt.Printf("Bem-vindo ao %s v%s\n", 
        info.GetName(), 
        info.GetVersion())
    fmt.Printf("Repositório: %s\n", 
        info.GetRepository())
    fmt.Printf("Autor: %s\n", 
        info.GetAuthor())
}
```

---

## 🎯 Benefícios

### Para Desenvolvedores

- **Única Fonte de Verdade**: Não mais arquivos de configuração espalhados
- **Type Safety**: Validação em tempo de compilação do acesso ao manifest
- **Workflows Automatizados**: Processo de build e release sem intervenção manual
- **Binários Consistentes**: Mesma configuração em todos os builds

### Para Usuários

- **Downloads Organizados**: Nomenclatura e organização clara de binários
- **Verificação de Integridade**: Checksums para segurança
- **Suporte Multi-plataforma**: Binários nativos para todas as plataformas
- **Tamanho Otimizado**: Compressão UPX reduz tempo de download

### Para Operações

- **Builds Reproduzíveis**: Processo de build determinístico
- **Trilha de Auditoria**: Rastreabilidade completa das configurações de build
- **Deploy Fácil**: Estrutura de artefatos padronizada
- **Gerenciamento de Versão**: Rastreamento e comparação automatizada de versões

---

## 🤝 Contribuindo

Ao contribuir para o sistema de distribuição do GoForge:

1. **Teste Localmente**: Execute `make build` para verificar mudanças
2. **Atualize o Manifest**: Modifique `info/manifest.json` conforme necessário
3. **Documente Mudanças**: Atualize esta documentação
4. **Teste CI/CD**: Verifique se o GitHub Actions ainda funciona

---

## 📄 Licença

Licença MIT - veja o arquivo [LICENSE](../LICENSE) para detalhes.

---

## 👤 Autor

**Rafael Mori** - [@rafa-mori](https://github.com/rafa-mori)

---

> Feito com 💙 para workflows de desenvolvimento Go simplificados
