<?php

define("ENV_FILE", __DIR__ . "/.env");
define("DB_FILE", __DIR__ . "/db.sql");

$config = [];
print("Importing user service database..." . PHP_EOL);
loadEnvVariables();
importDatabase();
print("Completed\n" . PHP_EOL);

function importDatabase()
{
    $short_options = "m::";
    $long_options = ["mode::"];
    $options = getopt($short_options, $long_options);
 
    if (isset($options["mode"]) && $options["mode"] == "testing") {
        importTestDatabase();
    }
    else {
        importAppDatabase();
    }
}

function importAppDatabase()
{
    global $config;
    $dsn = "mysql:host={$config['MYSQL_HOST']};dbname={$config['MYSQL_DB']}";
    $db = new \PDO($dsn, $config['MYSQL_USER'], $config['MYSQL_PASSWORD']);
    $sql = file_get_contents(DB_FILE);
    $db->exec($sql);
}

function importTestDatabase()
{
    global $config;
    $dsn = "mysql:host={$config['MYSQL_HOST_TEST']};dbname={$config['MYSQL_DB_TEST']}";
    $db = new \PDO($dsn, $config['MYSQL_USER_TEST'], $config['MYSQL_PASSWORD_TEST']);
    $sql = file_get_contents(DB_FILE);
    $db->exec($sql);
}

function loadEnvVariables()
{
    global $config;

    if (!is_readable(ENV_FILE)) {
        print('.env file is not readable' . PHP_EOL);
        exit;
    }

    $lines = file(ENV_FILE, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
    foreach ($lines as $line) {
        if (strpos(trim($line), '#') === 0) {
            continue;
        }

        list($name, $value) = explode('=', $line, 2);
        $name = trim($name);
        $value = trim($value);
        $config[$name] = $value;
    }
}