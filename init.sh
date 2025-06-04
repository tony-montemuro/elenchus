#!/bin/bash

set -e

DEFAULT_PASS="pass"
ENVIRONMENT="local"
MARIADB_USER_PASS=$DEFAULT_PASS
MARIADB_TEST_USER_PASS=$DEFAULT_PASS
DB_NAME="elenchus"
DB_USER="web"

echo "⚠️ Please ensure the following conditions are met:"
echo -e "\t- MariaDB is installed and running."
echo -e "\t- The root password must be provided to setup the database."

read -p "Is this a production environment? (Y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    ENVIRONMENT="production"
fi

echo -e "\n$ENVIRONMENT environment selected!"
echo -e "\n"

echo "🗃️ Setting up database(s)..."

read -s -p "Enter password for new MariaDB user '$DB_USER' (Press enter to use default '$MARIADB_USER_PASS'): "
echo
if [[ -n $REPLY ]]; then
    MARIADB_USER_PASS=$REPLY
fi

sql="
CREATE DATABASE IF NOT EXISTS $DB_NAME;
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$MARIADB_USER_PASS';
GRANT SELECT, INSERT, UPDATE, DELETE ON $DB_NAME.* TO '$DB_USER'@'localhost';
"

if [[ $ENVIRONMENT == "local" ]]; then
    read -s -p "Enter password for new MariaDB user 'test_$DB_USER' (Press enter to use default '$MARIADB_TEST_USER_PASS'): "
    echo
    if [[ -n $REPLY ]]; then
	MARIADB_TEST_USER_PASS=$REPLY
    fi

    sql="$sql
    CREATE DATABASE IF NOT EXISTS test_$DB_NAME;
    CREATE USER IF NOT EXISTS 'test_$DB_USER'@'localhost' IDENTIFIED BY '$MARIADB_TEST_USER_PASS';
    GRANT SELECT, INSERT, UPDATE, DELETE on test_$DB_NAME.* TO 'test_$DB_USER'@'localhost';
    "
fi

read -s -p "Enter MariaDB root password: " ROOT_PASS
echo

mariadb -u root -p"$ROOT_PASS" << EOF
$sql
EOF
echo "✅ MariaDB database(s) and user(s) generated!"
echo -e "\n"

echo "📦 Installing project dependencies..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies. Ensure Go is properly installed on your system."
    exit 1
else
    echo "✅ Dependencies installed!"
fi

if [[ $ENVIRONMENT == "local" ]]; then
    echo -e "\n"
    echo "🔐 Creating local TLS certificate..."

    # Create tls directory if it does not exist
    echo "Creating the ./tls directory..."
    mkdir -p tls
    cd tls

    # Generate TLS certificate
    echo "Generating tls certificate in ./tls directory..."
    go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" --rsa-bits=2048 --host=localhost
    if [ $? -eq 0 ]; then
	echo "✅ TLS certificate generated!"
    else
	echo "❌ Failed to generate TLS certificate. Ensure Go is properly installed on your system."
	exit 1
    fi

    cd ..
fi

echo -e "\n"
echo "========================================"
echo -e "\n"
echo "✅ Project initialization complete! Use:"
echo -e "go run ./cmd/web/\n"
echo -e "ℹ️ To view all configurable options, run the above command with the -help flag.\n"
echo -e "Database credentials:\n"

echo -e "\tDatabase name: $DB_NAME"
echo -e "\tUsername: $DB_USER"
if [[ $MARIADB_USER_PASS != $DEFAULT_PASS ]]; then
    echo -e "\tPassword set by yourself."
else
    echo -e "\tPassword: $DEFAULT_PASS"
fi
if [[ $ENVIRONMENT == "local" ]]; then
    echo
    echo -e "\tTest database name: test_$DB_NAME"
    echo -e "\tTest username: test_$DB_USER" 
    if [[ $MARIADB_TEST_USER_PASS != $DEFAULT_PASS ]]; then
	echo -e "\tTest user password set by yourself."
    else
	echo -e "\tTest user password: $DEFAULT_PASS"
    fi
fi

echo
if [[ $ENVIRONMENT == "local" ]]; then
    echo "⚠️ Your browser will show a security warning for using a self-signed certificate."
    echo "This is normal for development, just accept the warning to proceed."
fi

