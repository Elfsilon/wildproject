export CONFIG_PATH="./configs/dev.env"

echo "config path is: $CONFIG_PATH"
echo "Running up dev server"

go run cmd/api-server/main.go