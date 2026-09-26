# shellcheck shell=bash
# shellcheck disable=SC2016 # Variables expand inside the database container.
k create namespace database
k apply -k contents/helm-vs-kustomize/dependencies/mysql
k -n database rollout status deployment/mysql --timeout=360s
# Connect over the Service with the non-root sample credentials.
result=$(k -n database exec deployment/mysql -- sh -ec '
  export MYSQL_PWD="$MYSQL_PASSWORD"
  mysql --connect-timeout=10 -h mysql.database.svc -u "$MYSQL_USER" "$MYSQL_DATABASE" --batch --skip-column-names -e "
    CREATE TABLE e2e (id INT PRIMARY KEY, value INT);
    INSERT INTO e2e VALUES (1, 42);
    SELECT value FROM e2e WHERE id=1;
    DROP TABLE e2e;
  "
')
test "$result" = 42
