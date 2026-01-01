#!/bin/bash

# Скрипт для тестування STDIN програм

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║          Тестування STDIN програм                             ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Кольори для виводу
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}1. Тест basic_stdin.go${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "\n${GREEN}→ Через echo:${NC}"
echo -e "Іван\n25" | go run basic_stdin.go

echo -e "\n${GREEN}→ Через printf:${NC}"
printf "Марія\n22\n" | go run basic_stdin.go

echo -e "\n${GREEN}→ З файлу:${NC}"
cat > /tmp/test_input.txt << EOF
Петро
30
EOF
go run basic_stdin.go < /tmp/test_input.txt

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}2. Тест calculator_stdin.go${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "\n${GREEN}→ Додавання (10 + 5):${NC}"
echo -e "10\n+\n5" | go run calculator_stdin.go

echo -e "\n${GREEN}→ Множення (20 * 3):${NC}"
echo -e "20\n*\n3" | go run calculator_stdin.go

echo -e "\n${GREEN}→ Ділення (100 / 4):${NC}"
echo -e "100\n/\n4" | go run calculator_stdin.go

echo -e "\n${GREEN}→ З heredoc:${NC}"
go run calculator_stdin.go << EOF
15
-
7
EOF

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}3. Тест з циклом (множинні дані)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "\n${GREEN}→ Обробка списку імен:${NC}"
for name in Іван Марія Петро Олена; do
    echo "$name" | go run basic_stdin.go 2>/dev/null | grep "Привіт"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}4. Тест pipe з результатом${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "\n${GREEN}→ Збереження результату в файл:${NC}"
echo -e "10\n+\n5" | go run calculator_stdin.go 2>/dev/null > /tmp/calc_result.txt
echo "Результат збережено в /tmp/calc_result.txt:"
cat /tmp/calc_result.txt

echo ""
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║              ✅ Тестування завершено!                         ║"
echo "╚═══════════════════════════════════════════════════════════════╝"

# Очистка
rm -f /tmp/test_input.txt /tmp/calc_result.txt

