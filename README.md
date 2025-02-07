# **snyf**

*A lightweight experimental playground for scanning `go.mod` files in a monorepo.*

## **Usage**

Run Snyf in the root of your monorepo to analyze module relationships:

```
Usage of snyf [flags]:
  -extended
    	Include `require` and `replace` information for each module. !!WARNING!! This can be very verbose.
  -path string
      Path to a specific `go.mod` file or directory containing `go.mod` file to scan.
  -source string
    	Source directory to scan for `go.mod` files (default ".").
```

## **Output**

Snyf returns a structured JSON summary of all `go.mod` files and their interdependencies.

```json
[
  {
      "path": "services/foo",
      "go_version": "1.23.5",
      "module": "github.com/nonsense/foo",
      "require": [
        {
          "Path": "github.com/golang-migrate/migrate/v4",
          "Version": "v4.18.1"
        },
        {
          "Path": "github.com/gorilla/mux",
          "Version": "v1.8.0"
        }
      ],
      "replace": [
        {
          "Old": {
            "Path": "github.com/gorilla/mux"
          },
          "New": {
            "Path": "github.com/gorilla/mux",
            "Version": "v1.0.0"
          }
        }
      ],
      "depends_on": [
        "github.com/nonsense/bar",
        "github.com/nonsense/baz"
      ]
    },
    {
      "path": "libs/bar",
      "go_version": "1.22.1",
      "module": "github.com/nonsense/bar",
      "used_by": [
        "github.com/nonsense/foo"
      ]
    },
    {
      "path": "libs/baz",
      "go_version": "1.22.3",
      "module": "github.com/nonsense/baz",
      "depends_on": [
        "github.com/nonsense/bar"
      ],
      "used_by": [
        "github.com/nonsense/foo"
      ]
    }
  ]
```

Or if the `--path` flag is provided, returns a single `go.mod` file summary:

```json
{
  "path": "libs/baz",
  "go_version": "1.22.3",
  "module": "github.com/nonsense/baz",
  "depends_on": [
    "github.com/nonsense/bar"
  ],
  "used_by": [
    "github.com/nonsense/foo"
  ]
}
```

## **Status**

🚧 **Experimental** – For exploration and internal use only.
