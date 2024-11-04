# Application Development

## Prerequisites

- Go 1.21 installed. Download [here](https://go.dev/dl/).
- Postgres 14 installed. I used [Postgres.app](https://postgresapp.com/) for MacOS. For Windows, one option is [EnterpriseDB Postgres installer](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads).
- To work with Strimzi, install `kubectl` and `helm`.

## Init database

Using psql to your Postgres instance

```
drop database dev_db;
create database dev_db;
create user dev_db with encrypted password 'dev_db123';
grant all privileges on database dev_db to dev_db;
alter user dev_db with superuser;
```

## Run

Run migration

```bash
make migrate-up
```

Then compile and run dev_db

```bash
make run
```

dev_db will bootstrap create superadmin user and policy when run the first time on empty database.

## Working with Postgres

### Migration

To add or update existing table, add new migration file to the [migration](./dev_db/db/migrations/) folder, then run `migration-up` to update Postgres.

```bash
make migrate-up
```

If something is wrong with migration, you can either:
- force the migration. See [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#usage) for detail.
- drop all tables. `migration-up` will recreate all database, but you will loose any data that you have created.

To clear the database, run `migration-down`.

```bash
make migrate-down
```

### Updating Deployment Migrations

When there are any changes to database migrations SQL, make sure to run following before committing the code. 

```bash
make copy-migrations
```

This will copy migrations to [deployment folder](./deploy/compose) that will be used when application is deployed to update database schema.

### sqlc

To create functions that can be used for CRUD to the tables, add or update [queries](./dev_db/db/queries/). These are Postgres SQL syntax. 

Then run following to generate Go function.

```bash
make generate
```

Then you can call generated Go function in request handler.

For example, following SQL
```sql
-- name: ListOrganizations :many
SELECT * FROM organizations
ORDER BY organization_name;
```

Will become this function
```go
const ListOrganizations = `-- name: ListOrganizations :many
SELECT organization_id, organization_name, display_name, created_at, updated_at FROM organizations
ORDER BY organization_name
`

func (q *Queries) ListOrganizations(ctx context.Context) ([]*organizations, error) {
	rows, err := q.db.Query(ctx, ListOrganizations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*organizations{}
	for rows.Next() {
		var i organizations
		if err := rows.Scan(
			&i.OrganizationID,
			&i.OrganizationName,
			&i.DisplayName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
```

Then you can use it
```go
func (h *OrganizationHandler) listOrganization(c echo.Context) error {
	ctx := context.Background()
	orgs, err := h.q.ListOrganizations(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": orgs,
	})
}
```

## API documentation

API documentation available as [Insomnia YAML](./docs) that you can import into your local Insomnia.
