### A Service to Crawl Metadata in your GCS Datalake and Provide an API for it

Fetch metadata from blobs in GCS structure and store it in Postgres and make it
available for querying and search via an API


#### High-level architecture
##### Indexer (CLI / service in Go)

Walks GCS (lists Parquet objects).

Reads Parquet metadata only (schema, row count, etc.).

Writes a compact catalog (e.g., SQLite or a small Postgres DB).

##### TUI app (Go)

- Connects to the catalog DB.

- Lets you browse datasets → schemas → files/partitions.

- This lets you keep I/O and GCS credentials out of the interactive TUI, and the TUI stays snappy because it only hits the catalog.

##### Stack
###### GCS + Parquet metadata

- GCS client: cloud.google.com/go/storage.

- Parquet metadata: popular options are github.com/apache/arrow/go/vX/parquet or similar; you only need to open the file and read the footer for metadata.

###### Catalog storage

- SQLite via modernc.org/sqlite or github.com/mattn/go-sqlite3 (if CGO is OK).

- Schema: datasets, files, columns tables like described above

###### TUI library

Bubble Tea (github.com/charmbracelet/bubbletea) + lipgloss / bubbles for widgets, or

tview (github.com/rivo/tview) for a more “widget-oriented” experience (tables, trees, forms).

For a data catalog, tview.Table + TreeView works very well: left side tree for datasets/partitions, right side table for schema.

#### Data model (practical version)
##### In Go/SQL terms:

***datasets***

- id INTEGER PRIMARY KEY
- name TEXT UNIQUE
- description TEXT
- owner TEXT
- tags TEXT (comma-separated or JSON)

***files***
- id INTEGER PRIMARY KEY
- dataset_id INTEGER
- gcs_uri TEXT
- size_bytes INTEGER
- row_count INTEGER
- created_at TEXT
- updated_at TEXT
- partition_values TEXT (JSON string)

***columns***

- id INTEGER PRIMARY KEY
- dataset_id INTEGER
- name TEXT
- path TEXT // user.address.city
- physical_type TEXT
- logical_type TEXT
- nullable BOOLEAN

The TUI only ever queries this DB.

TUI UX sketch
Using tview as an example:

Layout

Left: TreeView

Top nodes: datasets.

Children: partitions (derived from partition_values).

Right: Table

Shows schema or file list depending on what’s selected.

***Key bindings***

↑/↓ to move selection, → to expand dataset, ← to collapse.

s to switch right pane between “Schema” and “Files”.

/ to open a search modal: filter datasets or columns by name.

q to quit.

Views

Schema view: columns table with name | type | logical | nullable.

Files view: partition | file count | total rows | total size.

Implementation steps
Indexer MVP

CLI app: catalog-indexer.

Inputs: bucket, prefix, path to SQLite DB.

***Steps:***

- List .parquet objects under prefix.

- For each:

- Read GCS object metadata.

- Open Parquet, read metadata (schema + row count).

- Upsert into datasets/files/columns.

##### TUI MVP

- Create a small Go module catalog-tui.

- Initialize DB connection (read-only) at startup.

- Build main layout (tree + table).

- Implement dataset selection → query schema and render.

***Refinement**

Add search over datasets/columns.

Add a status bar: current dataset, number of columns, last indexed time.

Add a config file (~/.gcs-catalog.yaml) for DB path, default bucket, etc.
