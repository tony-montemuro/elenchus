#!/bin/bash

set -e

DB_DEFAULT_PASS="pass"
ENVIRONMENT="local"
DB_USER_PASS=$DB_DEFAULT_PASS
DB_MIGRATION_PASS=$DB_DEFAULT_PASS
DB_TEST_USER_PASS=$DB_DEFAULT_PASS
DB_NAME="elenchus"
DB_USER="web"
DB_MIGRATION_USER="migration"

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

read -s -p "Enter password for new MariaDB user '$DB_USER' (Press enter to use default '$DB_DEFAULT_PASS'): "
echo
if [[ -n $REPLY ]]; then
    DB_USER_PASS=$REPLY
fi

read -s -p "Enter password for new MariaDB migration user '$DB_MIGRATION_USER' (Press enter to use default '$DB_DEFAULT_PASS'): "
echo
if [[ -n $REPLY ]]; then
    DB_MIGRATION_PASS=$REPLY
fi

sql="
CREATE DATABASE IF NOT EXISTS $DB_NAME;
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_USER_PASS';
GRANT SELECT, INSERT, UPDATE, DELETE ON $DB_NAME.* TO '$DB_USER'@'localhost';
CREATE USER IF NOT EXISTS '$DB_MIGRATION_USER'@'localhost' IDENTIFIED BY '$DB_MIGRATION_PASS';
GRANT CREATE, DROP, ALTER, INDEX, REFERENCES, SELECT, INSERT, UPDATE, DELETE ON $DB_NAME.* TO '$DB_MIGRATION_USER'@'localhost';
"

dbconfig="development:
    dialect: mysql
    datasource: $DB_MIGRATION_USER:$DB_MIGRATION_PASS@/$DB_NAME?parseTime=true
    dir: migrations
"

if [[ $ENVIRONMENT == "local" ]]; then
    read -s -p "Enter password for new MariaDB user 'test_$DB_USER' (Press enter to use default '$DB_TEST_USER_PASS'): "
    echo
    if [[ -n $REPLY ]]; then
	DB_TEST_USER_PASS=$REPLY
    fi

    sql="$sql
    CREATE DATABASE IF NOT EXISTS test_$DB_NAME;
    CREATE USER IF NOT EXISTS 'test_$DB_USER'@'localhost' IDENTIFIED BY '$DB_TEST_USER_PASS';
    GRANT CREATE, DROP, ALTER, INDEX, REFERENCES, SELECT, INSERT, UPDATE, DELETE on test_$DB_NAME.* TO 'test_$DB_USER'@'localhost';
    "

    dbconfig="$dbconfig
test:
    dialect: mysql
    datasource: test_$DB_USER:$DB_TEST_USER_PASS@/test_$DB_NAME?parseTime=true
    dir: migrations 
    "
fi

read -s -p "Enter MariaDB root password: " ROOT_PASS
echo

mariadb -u root -p"$ROOT_PASS" << EOF
$sql
EOF

echo "$dbconfig" > dbconfig.yml

go install github.com/rubenv/sql-migrate/...@latest

echo "✅ MariaDB database(s), user(s), and migration config (.dbconfig) generated! sql-migrate installed!"
echo -e "\n"

echo "📦 Installing project dependencies..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies. Ensure Go is properly installed on your system."
    exit 1
else
    echo "✅ Dependencies installed!"
fi

echo -e "\n"
echo "⬆️ Running database migrations..."
sql-migrate up
echo "✅ Database fully migrated!"

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
if [[ $DB_USER_PASS != $DB_DEFAULT_PASS ]]; then
    echo -e "\tPassword set by yourself."
else
    echo -e "\tPassword: $DB_DEFAULT_PASS"
fi
echo
echo -e "\tMigration username: $DB_MIGRATION_USER"
if [[ $DB_MIGRATION_PASS != $DB_DEFAULT_PASS ]]; then
    echo -e "\tPassword set by yourself."
else
    echo -e "\tPassword: $DB_DEFAULT_PASS"
fi
if [[ $ENVIRONMENT == "local" ]]; then
    echo
    echo -e "\tTest database name: test_$DB_NAME"
    echo -e "\tTest username: test_$DB_USER" 
    if [[ $DB_TEST_USER_PASS != $DB_DEFAULT_PASS ]]; then
	echo -e "\tTest user password set by yourself."
    else
	echo -e "\tTest user password: $DB_DEFAULT_PASS"
    fi
fi

echo
if [[ $ENVIRONMENT == "local" ]]; then
    echo "⚠️ Your browser will show a security warning for using a self-signed certificate."
    echo "This is normal for development, just accept the warning to proceed."
fi

