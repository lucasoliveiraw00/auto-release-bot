# 🚀 Auto-Release-Bot

## 📌 Sobre o Projeto

O **Auto-Release-Bot** é um **bot em Golang** projetado para enviar **notificações automáticas** para uma **sala no Google Chat**, mantendo a equipe informada sobre eventos críticos do **ciclo de releases**.

---

### 📢 Como Funciona?

O **Auto-Release-Bot** monitora diariamente o calendário de versões e envia **notificações programadas** para garantir que a equipe esteja sempre informada.

### ✅ Benefícios

- ✅ **Acompanhamento contínuo** do ciclo de releases
- ✅ **Notificações automáticas** sem necessidade de intervenção manual
- ✅ **Redução de esquecimentos** e aumento da organização

---

### 📅 Eventos Monitorados

O bot acompanha o **calendário de releases** e envia notificações para garantir que a equipe esteja sempre atualizada.

| 🚨 **Evento**                   | 📅 **Quando ocorre?**           | 🔔 **Motivo da Notificação**                                           |
| ------------------------------- | ------------------------------- | ---------------------------------------------------------------------- |
| **🛠 Início da preparação**      | No **1º dia** da preparação     | Início dos testes das demandas alocadas na Release Candidate (RC).     |
| **⏳ Último dia da preparação** | No **último dia** da preparação | Última chance para validar e concluir os testes antes da entrega.      |
| **🚀 Entrega da Release**       | **2 dias antes**                | Aviso para revisão final das tarefas e preparação da entrega.          |
| **🚀 Entrega da Release**       | **1 dia antes**                 | Último alerta para garantir que todas as demandas estejam finalizadas. |
| **🚀 Entrega da Release**       | **No dia da entrega**           | Confirmação da entrega e aviso final para a equipe.                    |

---

### 🛠 **Modos de Execução**

O script pode ser executado de duas formas:

- **🔹 Manualmente** → Executado sob demanda pelo usuário.
- **🔹 Automaticamente** → Configurado via **cron job** para execução diária.

---

### ⚙️ Tecnologias Utilizadas

O **Auto-Release-Bot** é desenvolvido com tecnologias robustas e eficientes:

- **Golang** 🦫 → Linguagem principal do projeto.
- **Google Chat Webhooks** 💬 → Para envio das notificações automáticas.

---

### 📦 Pré-requisitos

Antes de rodar o projeto, certifique-se de ter instalado:

- **[Go](https://golang.org/doc/install) (versão 1.23.4 ou superior)** 🦫
- **[Make](https://www.gnu.org/software/make/)** ⚙️ (para facilitar a execução de comandos)

Para verificar se estão instalados, execute:

```sh
go version  # Deve exibir a versão do Go instalada
make --version  # Deve exibir a versão do Make instalada
```

---

### 🔧 Configuração

Para utilizar o bot, siga os passos abaixo:

#### 1️⃣ **Criar um Webhook no Google Chat**

Siga as instruções da [documentação oficial do Google Chat](https://support.google.com/chat/answer/7655820?hl=en&co=GENIE.Platform%3DDesktop) para criar um **Webhook**.

#### 2️⃣ **Definir as Variáveis de Ambiente**

Configure as seguintes variáveis no seu ambiente de execução:

```ini
# 🛠️ Configuração do Webhook do Google Chat
GOOGLE_CHAT_WEBHOOK=""

# Definir uma data específica para testes (YYYY-MM-DD). Deixe vazio para usar a data atual.
MOCK_DATE=

# 🛠️ Configurações do PR Checker
GITHUB_OWNER=""
GITHUB_REPO=""
GITHUB_TOKEN=""
SONAR_TOKEN=""
SONAR_OWNER=""

# 📅 URLs de Referência
VERSION_CALENDAR_URL=""
QUALITY_CRITERIA_URL=""

# 🛠️ Sonar Thresholds
SONAR_NEW_SECURITY_HOTSPOTS=0
SONAR_NEW_VIOLATIONS=0
SONAR_NEW_ACCEPTED_ISSUES=0
SONAR_NEW_COVERAGE=75
SONAR_NEW_DUPLICATED_LINES_DENSITY=0
```

#### 3️⃣ **Executar o bot**

Execute o script manualmente:

```ini
make run-release   # Executa o monitoramento de releases
make run-prchecker # Executa a verificação de PRs no GitHub
```

#### 4️⃣ **(Opcional) Configurar um Cron Job**

Para automatizar a execução diária, adicione uma entrada no cron:

```ini
0 9 * * * /caminho/para/o/binario/release-notifier
```

_(Esse exemplo executa o bot todos os dias às **09:00 AM**)._

### 📋 Exemplo de Notificação no Google Chat

Abaixo está um exemplo de como o **Auto-Release-Bot** envia notificações automáticas para o **Google Chat**:

![Exemplo de Notificação no Google Chat](https://raw.githubusercontent.com/lucasoliveiraw00/auto-release-bot/main/assets/google-chat-notification.png)
