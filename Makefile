BIN_DIR=bin
PRCHECKER_BIN=$(BIN_DIR)/prchecker
RELEASE_BIN=$(BIN_DIR)/release

PRCHECKER_SRC=cmd/prchecker/main.go
RELEASE_SRC=cmd/release/main.go

# Rodar cada serviço individualmente
run-prchecker:
	go run $(PRCHECKER_SRC)

run-release:
	go run $(RELEASE_SRC)

# Compila os binários na pasta bin/
build-all:
	mkdir -p $(BIN_DIR)
	go build -o $(PRCHECKER_BIN) $(PRCHECKER_SRC)
	go build -o $(RELEASE_BIN) $(RELEASE_SRC)

# Remove os binários compilados
clean:
	rm -rf $(BIN_DIR)

# Formata o código conforme padrão do Go
fmt:
	go fmt ./...

# Faz lint do código
lint:
	go vet ./...

help:
	@echo ""
	@echo "💡 Comandos disponíveis:"
	@echo "-------------------------------------------"
	@echo "  make run-prchecker  -> Roda o script PR Checker - Alerta via Google Chat"
	@echo "  make run-release    -> Roda o script Release - Alerta via Google Chat"
	@echo "  make build          -> Compila ambos os binários"
	@echo "  make clean          -> Remove binários compilados"
	@echo "  make fmt            -> Formata o código conforme padrão do Go"
	@echo "  make lint           -> Executa análise estática (lint)"
	@echo "  make help           -> Exibe essa mensagem"
	@echo "-------------------------------------------"
