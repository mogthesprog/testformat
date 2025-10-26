# TCF Proposal 001: Universal JSON Test Coverage Format

## Metadata

- **Title**: Universal JSON Test Coverage Format
- **Author**: Claude (AI Assistant)
- **Status**: Draft
- **Type**: Standards Track
- **Created**: 2025-10-26
- **Format Version**: 1.0.0

## Abstract

This proposal defines a universal JSON format for representing code coverage data that can accommodate the diverse coverage metrics used across different programming languages, testing frameworks, and coverage tools. The format is designed to be a superset of capabilities found in existing formats including JaCoCo XML, Cobertura XML, LCOV, Coverage.py JSON, Istanbul JSON, and the Codecov custom format.

## Motivation

The test coverage ecosystem lacks a standardized, de facto open format for representing coverage data in JSON. While numerous XML-based formats exist (JaCoCo, Cobertura, Clover), JSON formats are either proprietary, tool-specific, or overly simplified. This creates several problems:

1. **Interoperability**: Tools must implement parsers for multiple XML formats
2. **Modern tooling**: JSON is more widely supported in modern CI/CD pipelines and APIs
3. **Extensibility**: XML schemas can be rigid and difficult to extend
4. **Developer experience**: JSON is easier to generate, parse, and debug
5. **Coverage service compatibility**: Services like Codecov support simplified JSON but lack comprehensive field definitions

By creating a universal JSON format that supports the full range of coverage metrics from existing formats, we enable:

- Easy conversion from existing XML/text formats to JSON
- Unified coverage processing pipelines
- Better integration with modern development tools
- Comprehensive coverage data representation across all languages

## Rationale

### Why Not Use Existing Formats?

**Codecov Custom Format**: Too simplified; only supports line-level coverage with basic hit counts, no function/method coverage, no branch details, no metadata.

**Istanbul Coverage Format**: JavaScript-centric; uses statement IDs rather than line numbers as primary metric, complex mapping structures.

**Coverage.py JSON**: Python-specific conventions; limited branch coverage representation, no method-level granularity.

**JaCoCo/Cobertura XML**: Comprehensive but XML-based; not ideal for modern JSON-based toolchains.

### Design Goals

1. **Comprehensive**: Support all coverage types (line, branch, function, method, statement, condition, decision, MC/DC)
2. **Language-agnostic**: Work for Java, Python, JavaScript, Go, C++, etc.
3. **Hierarchical**: Support package/namespace/module organization
4. **Extensible**: Allow custom metadata and future coverage types
5. **Simple to generate**: Easy for tools to produce valid coverage reports
6. **Simple to consume**: Easy for tools to parse and aggregate
7. **Human-readable**: JSON structure that developers can understand
8. **Backwards compatible**: Support migration from existing formats

## Specification

### Format Version

All coverage reports MUST include a format version to enable future evolution:

```json
{
  "version": "1.0.0"
}
```

### Root Structure

```json
{
  "version": "1.0.0",
  "meta": { /* metadata object */ },
  "coverage": { /* coverage data object */ },
  "totals": { /* summary statistics object */ }
}
```

### Metadata Object (`meta`)

The metadata object contains information about the coverage report generation:

```json
{
  "meta": {
    "format": "testformat-json-coverage",
    "formatVersion": "1.0.0",
    "generatedAt": "2025-10-26T12:34:56Z",
    "generatedBy": {
      "name": "jacoco",
      "version": "0.8.11"
    },
    "source": {
      "type": "conversion",
      "originalFormat": "jacoco-xml",
      "originalFile": "coverage.xml"
    },
    "repository": {
      "root": "/home/user/project",
      "vcs": "git",
      "commit": "abc123def456",
      "branch": "main"
    },
    "settings": {
      "branchCoverage": true,
      "conditionCoverage": false,
      "mcdcCoverage": false,
      "functionCoverage": true,
      "lineCoverage": true
    }
  }
}
```

**Field Descriptions**:

- `format` (string, required): Format identifier, MUST be `"testformat-json-coverage"`
- `formatVersion` (string, required): Semantic version of this format specification
- `generatedAt` (string, optional): ISO 8601 timestamp
- `generatedBy` (object, optional): Tool that generated the report
  - `name` (string): Tool name
  - `version` (string): Tool version
- `source` (object, optional): Information about source of coverage data
  - `type` (string): One of `"native"`, `"conversion"`, `"aggregation"`
  - `originalFormat` (string): Original format name if converted
  - `originalFile` (string): Path to original file if converted
- `repository` (object, optional): Repository information
  - `root` (string): Repository root path
  - `vcs` (string): Version control system (e.g., "git", "svn")
  - `commit` (string): Commit hash/revision
  - `branch` (string): Branch name
- `settings` (object, required): Coverage types measured
  - `branchCoverage` (boolean): Whether branch coverage is measured
  - `conditionCoverage` (boolean): Whether condition coverage is measured
  - `mcdcCoverage` (boolean): Whether MC/DC coverage is measured
  - `functionCoverage` (boolean): Whether function coverage is measured
  - `lineCoverage` (boolean): Whether line coverage is measured

### Coverage Data Object (`coverage`)

The coverage object maps file paths (relative to repository root) to file coverage data:

```json
{
  "coverage": {
    "src/main/java/com/example/Calculator.java": {
      /* file coverage object */
    },
    "src/main/python/utils.py": {
      /* file coverage object */
    }
  }
}
```

**Key Requirements**:

- File paths MUST be relative to repository root
- File paths MUST use forward slashes `/` regardless of OS
- File paths MUST NOT start with `/` (i.e., not absolute paths)

### File Coverage Object

Each file coverage object contains line-level, branch-level, and function-level coverage:

```json
{
  "lines": {
    "1": { "hits": 0, "branches": null },
    "5": { "hits": 12, "branches": null },
    "10": {
      "hits": 8,
      "branches": {
        "total": 2,
        "covered": 1,
        "missed": 1,
        "details": [
          { "id": 0, "type": "if", "hits": 8 },
          { "id": 1, "type": "else", "hits": 0 }
        ]
      }
    }
  },
  "functions": [
    {
      "name": "add",
      "signature": "(int, int)",
      "lineStart": 10,
      "lineEnd": 15,
      "hits": 8,
      "coverage": {
        "instructions": { "covered": 12, "missed": 0, "total": 12 },
        "branches": { "covered": 1, "missed": 1, "total": 2 },
        "lines": { "covered": 5, "missed": 0, "total": 5 }
      }
    }
  ],
  "classes": [
    {
      "name": "Calculator",
      "package": "com.example",
      "lineStart": 5,
      "lineEnd": 50,
      "methods": [
        /* references to functions */
      ],
      "coverage": {
        "instructions": { "covered": 45, "missed": 5, "total": 50 },
        "branches": { "covered": 8, "missed": 2, "total": 10 },
        "lines": { "covered": 20, "missed": 2, "total": 22 },
        "methods": { "covered": 4, "missed": 0, "total": 4 }
      }
    }
  ],
  "summary": {
    "lines": { "covered": 20, "missed": 2, "total": 22, "percent": 90.91 },
    "branches": { "covered": 8, "missed": 2, "total": 10, "percent": 80.0 },
    "functions": { "covered": 4, "missed": 0, "total": 4, "percent": 100.0 },
    "instructions": { "covered": 45, "missed": 5, "total": 50, "percent": 90.0 }
  }
}
```

#### Line Coverage Object (`lines`)

Maps line numbers (as strings) to line coverage data:

```json
{
  "hits": 12,
  "branches": null | {
    "total": 2,
    "covered": 1,
    "missed": 1,
    "details": [
      {
        "id": 0,
        "type": "if" | "else" | "case" | "switch" | "catch" | "conditional" | "jump",
        "hits": 8,
        "conditions": null | {
          "total": 2,
          "covered": 1,
          "coverage": "50% (1/2)"
        }
      }
    ]
  }
}
```

**Field Descriptions**:

- `hits` (integer, required): Number of times line was executed (0 = not covered)
- `branches` (object or null, optional): Branch information if line contains branches
  - `total` (integer): Total number of branches
  - `covered` (integer): Number of branches covered
  - `missed` (integer): Number of branches not covered
  - `details` (array): Per-branch information
    - `id` (integer): Branch identifier (unique within line)
    - `type` (string): Branch type
    - `hits` (integer): Hit count for this branch
    - `conditions` (object, optional): Condition/decision coverage for this branch
      - `total` (integer): Total conditions
      - `covered` (integer): Covered conditions
      - `coverage` (string): Human-readable coverage (e.g., "50% (1/2)")

**Special Values**:

- Line numbers with no executable code SHOULD NOT be included in the `lines` object
- A line with `hits: 0` indicates executable code that was not covered
- `branches: null` or omitted means the line has no branches

#### Function Coverage Object (`functions`)

Array of function/method coverage data:

```json
{
  "name": "calculateTotal",
  "signature": "(items: Item[]): number",
  "descriptor": null,
  "lineStart": 25,
  "lineEnd": 35,
  "columnStart": 0,
  "columnEnd": 1,
  "hits": 42,
  "coverage": {
    "instructions": { "covered": 20, "missed": 0, "total": 20 },
    "branches": { "covered": 4, "missed": 0, "total": 4 },
    "lines": { "covered": 8, "missed": 0, "total": 8 },
    "statements": { "covered": 15, "missed": 0, "total": 15 }
  }
}
```

**Field Descriptions**:

- `name` (string, required): Function/method name
- `signature` (string, optional): Function signature (language-specific)
- `descriptor` (string, optional): JVM descriptor for Java methods (e.g., `"(II)I"`)
- `lineStart` (integer, required): Starting line number
- `lineEnd` (integer, optional): Ending line number
- `columnStart` (integer, optional): Starting column (0-indexed)
- `columnEnd` (integer, optional): Ending column (0-indexed)
- `hits` (integer, required): Number of times function was called
- `coverage` (object, optional): Coverage metrics for this function
  - Each metric follows the Counter Pattern (see below)

#### Class Coverage Object (`classes`)

Array of class/type coverage data (for object-oriented languages):

```json
{
  "name": "UserService",
  "package": "com.example.services",
  "namespace": "App\\Services",
  "sourceFile": "UserService.java",
  "lineStart": 10,
  "lineEnd": 150,
  "methods": ["getUserById", "createUser", "updateUser"],
  "coverage": {
    "instructions": { "covered": 200, "missed": 50, "total": 250 },
    "branches": { "covered": 30, "missed": 10, "total": 40 },
    "lines": { "covered": 80, "missed": 20, "total": 100 },
    "methods": { "covered": 8, "missed": 2, "total": 10 },
    "complexity": { "value": 42 }
  }
}
```

**Field Descriptions**:

- `name` (string, required): Class/type name
- `package` (string, optional): Java-style package name
- `namespace` (string, optional): Namespace (for languages using namespaces)
- `sourceFile` (string, optional): Source filename (without path)
- `lineStart` (integer, required): Starting line number
- `lineEnd` (integer, optional): Ending line number
- `methods` (array of strings, optional): Method names in this class
- `coverage` (object, optional): Coverage metrics for this class

### Counter Pattern

Coverage counters follow a consistent pattern across all coverage types:

```json
{
  "covered": 45,
  "missed": 5,
  "total": 50,
  "percent": 90.0,
  "value": null
}
```

**Field Descriptions**:

- `covered` (integer, optional): Number of covered items
- `missed` (integer, optional): Number of missed items
- `total` (integer, required): Total number of items (`covered + missed`)
- `percent` (number, optional): Coverage percentage (0-100, rounded to 2 decimals)
- `value` (integer, optional): Used for metrics like complexity that don't have covered/missed

**Coverage Types**:

- `instructions`: Individual bytecode instructions or machine instructions
- `branches`: Branch decision points (if/else, switch cases, etc.)
- `lines`: Source code lines
- `statements`: Logical statements (may differ from lines)
- `functions` / `methods`: Functions or methods
- `classes`: Classes or types
- `complexity`: Cyclomatic or other complexity metrics
- `conditions`: Individual boolean sub-expressions
- `decisions`: Decision outcomes (for MC/DC)

### Summary Object (`summary` and `totals`)

The summary provides aggregated coverage statistics. It appears at:

1. File level (`file.summary`) - aggregates all coverage in a file
2. Root level (`totals`) - aggregates all coverage across all files

```json
{
  "lines": { "covered": 1250, "missed": 150, "total": 1400, "percent": 89.29 },
  "branches": { "covered": 450, "missed": 50, "total": 500, "percent": 90.0 },
  "functions": { "covered": 95, "missed": 5, "total": 100, "percent": 95.0 },
  "instructions": { "covered": 5000, "missed": 500, "total": 5500, "percent": 90.91 },
  "statements": { "covered": 1200, "missed": 100, "total": 1300, "percent": 92.31 },
  "classes": { "covered": 45, "missed": 5, "total": 50, "percent": 90.0 },
  "methods": { "covered": 95, "missed": 5, "total": 100, "percent": 95.0 },
  "complexity": { "value": 456 }
}
```

All fields are optional; only include metrics that were measured.

## Coverage Type Support Matrix

This format is designed to support coverage data from all major coverage tools:

| Format | Lines | Branches | Functions | Instructions | Statements | Conditions | Complexity |
|--------|-------|----------|-----------|--------------|------------|------------|------------|
| **JaCoCo** | ✓ | ✓ | ✓ (methods) | ✓ | - | - | ✓ |
| **Cobertura** | ✓ | ✓ | - | - | - | ✓ | ✓ |
| **LCOV** | ✓ | ✓ | ✓ | - | - | - | - |
| **Coverage.py** | ✓ | ✓ (partial) | - | - | - | - | - |
| **Istanbul** | ✓ | ✓ | ✓ | - | ✓ | - | - |
| **Codecov** | ✓ | ✓ (simple) | - | - | - | - | - |
| **Go Coverage** | ✓ | - | - | - | ✓ | - | - |
| **Clover** | ✓ | ✓ | ✓ | - | ✓ | ✓ | ✓ |
| **This Format** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

## Mapping from Existing Formats

### From JaCoCo XML

JaCoCo provides the most comprehensive coverage metrics for JVM languages.

**Mapping**:

- `<report>` → root object
- `<sessioninfo>` → `meta.generatedAt`, `meta.source`
- `<package>` → organize files by package prefix in file paths
- `<class>` → `coverage[file].classes[]`
- `<method>` → `coverage[file].functions[]`
- `<sourcefile>` → `coverage[file]`
- `<line>` attributes:
  - `nr` → line number key
  - `mi`, `ci` → `lines[n].coverage.instructions`
  - `mb`, `cb` → `lines[n].branches`
- `<counter>` → `summary` and `totals` with counter type mapping

**Counter Type Mapping**:
- `INSTRUCTION` → `instructions`
- `LINE` → `lines`
- `BRANCH` → `branches`
- `COMPLEXITY` → `complexity`
- `METHOD` → `methods`
- `CLASS` → `classes`

### From Cobertura XML

Cobertura is widely used for Java, Python, and other languages.

**Mapping**:

- `<coverage>` attributes (`line-rate`, `branch-rate`) → `totals`
- `<package>` → organize files by package
- `<class>` → `coverage[file].classes[]`
- `<line>` attributes:
  - `number` → line number key
  - `hits` → `lines[n].hits`
  - `branch="true"` → indicates line has branches
  - `condition-coverage` → parse `"50% (1/2)"` format for branch coverage
- `<conditions>/<condition>` → `lines[n].branches.details[]`

### From LCOV (`.info` files)

LCOV is used for C/C++ coverage (gcov) and other languages.

**Mapping**:

- `SF:` (source file) → file path
- `FN:` (function) → `functions[].name` and `lineStart`
- `FNDA:` (function data) → `functions[].hits`
- `FNF:`, `FNH:` → `summary.functions.total` and `.covered`
- `DA:` (line data) → `lines[n].hits`
- `LF:`, `LH:` → `summary.lines.total` and `.covered`
- `BRDA:` (branch data) → `lines[n].branches.details[]`
  - Format: `BRDA:<line>,<block>,<branch>,<taken>`
  - `<taken>` = hit count or `-` for zero
- `BRF:`, `BRH:` → `summary.branches.total` and `.covered`

### From Coverage.py JSON

Python coverage tool's native JSON format.

**Mapping**:

- `meta` → `meta` (with format conversion)
- `files[path]`:
  - `executed_lines` → create `lines[n].hits = 1` for each line
  - `missing_lines` → create `lines[n].hits = 0` for each line
  - `excluded_lines` → omit from coverage object
  - `summary` → `summary` with field name mapping
- Branch coverage (when enabled):
  - Coverage.py provides limited branch info; map to `branches.total/covered/missed`

### From Istanbul Coverage JSON

JavaScript coverage format used by Istanbul/NYC.

**Mapping**:

- File path keys → file paths (convert absolute to relative)
- `path` → verify against key
- `s` (statement map) → `statements` counter
- `b` (branch map) → `branches` with details
- `f` (function map) → `functions` with hits
- `l` (line map) → `lines[n].hits`
- `statementMap` → statement location metadata
- `branchMap` → `lines[n].branches.details[]`
  - `locations` array → map to branch details with `id` = index
- `fnMap` → `functions[]` with name, lineStart from `loc`

### From Codecov Custom Format

Simplest format; line-level coverage only.

**Mapping**:

- `coverage[file]` → file path
- Array indices → line numbers
- Array values:
  - `null` → omit line (no executable code)
  - `0` → `lines[n].hits = 0`
  - `1` or `true` → `lines[n].hits = 1`
  - `"1/3"` → parse as partial branch coverage: `branches.covered = 1, branches.total = 3`
- `messages[file][line]` → could map to annotations/comments (future extension)

## Example: Complete Coverage Report

```json
{
  "version": "1.0.0",
  "meta": {
    "format": "testformat-json-coverage",
    "formatVersion": "1.0.0",
    "generatedAt": "2025-10-26T14:23:45Z",
    "generatedBy": {
      "name": "testformat",
      "version": "0.1.0"
    },
    "source": {
      "type": "conversion",
      "originalFormat": "jacoco-xml",
      "originalFile": "target/site/jacoco/jacoco.xml"
    },
    "repository": {
      "root": "/home/user/myproject",
      "vcs": "git",
      "commit": "a1b2c3d4e5f6",
      "branch": "main"
    },
    "settings": {
      "branchCoverage": true,
      "conditionCoverage": false,
      "mcdcCoverage": false,
      "functionCoverage": true,
      "lineCoverage": true
    }
  },
  "coverage": {
    "src/main/java/com/example/Calculator.java": {
      "lines": {
        "5": { "hits": 1, "branches": null },
        "10": { "hits": 42, "branches": null },
        "15": {
          "hits": 42,
          "branches": {
            "total": 2,
            "covered": 2,
            "missed": 0,
            "details": [
              { "id": 0, "type": "if", "hits": 30 },
              { "id": 1, "type": "else", "hits": 12 }
            ]
          }
        },
        "16": { "hits": 30, "branches": null },
        "18": { "hits": 12, "branches": null },
        "20": { "hits": 42, "branches": null }
      },
      "functions": [
        {
          "name": "add",
          "signature": "(int, int)",
          "descriptor": "(II)I",
          "lineStart": 10,
          "lineEnd": 20,
          "hits": 42,
          "coverage": {
            "instructions": { "covered": 12, "missed": 0, "total": 12, "percent": 100.0 },
            "branches": { "covered": 2, "missed": 0, "total": 2, "percent": 100.0 },
            "lines": { "covered": 5, "missed": 0, "total": 5, "percent": 100.0 }
          }
        }
      ],
      "classes": [
        {
          "name": "Calculator",
          "package": "com.example",
          "sourceFile": "Calculator.java",
          "lineStart": 5,
          "lineEnd": 50,
          "methods": ["add", "subtract", "multiply", "divide"],
          "coverage": {
            "instructions": { "covered": 45, "missed": 5, "total": 50, "percent": 90.0 },
            "branches": { "covered": 8, "missed": 2, "total": 10, "percent": 80.0 },
            "lines": { "covered": 20, "missed": 2, "total": 22, "percent": 90.91 },
            "methods": { "covered": 4, "missed": 0, "total": 4, "percent": 100.0 },
            "complexity": { "value": 12 }
          }
        }
      ],
      "summary": {
        "lines": { "covered": 20, "missed": 2, "total": 22, "percent": 90.91 },
        "branches": { "covered": 8, "missed": 2, "total": 10, "percent": 80.0 },
        "functions": { "covered": 4, "missed": 0, "total": 4, "percent": 100.0 },
        "instructions": { "covered": 45, "missed": 5, "total": 50, "percent": 90.0 }
      }
    },
    "src/main/python/utils.py": {
      "lines": {
        "1": { "hits": 1, "branches": null },
        "5": { "hits": 20, "branches": null },
        "10": {
          "hits": 15,
          "branches": {
            "total": 2,
            "covered": 1,
            "missed": 1,
            "details": [
              { "id": 0, "type": "if", "hits": 15 },
              { "id": 1, "type": "else", "hits": 0 }
            ]
          }
        },
        "11": { "hits": 15, "branches": null },
        "13": { "hits": 0, "branches": null }
      },
      "functions": [
        {
          "name": "format_number",
          "signature": "(value: int) -> str",
          "lineStart": 5,
          "lineEnd": 13,
          "hits": 20,
          "coverage": {
            "lines": { "covered": 3, "missed": 1, "total": 4, "percent": 75.0 },
            "branches": { "covered": 1, "missed": 1, "total": 2, "percent": 50.0 }
          }
        }
      ],
      "summary": {
        "lines": { "covered": 4, "missed": 1, "total": 5, "percent": 80.0 },
        "branches": { "covered": 1, "missed": 1, "total": 2, "percent": 50.0 },
        "functions": { "covered": 1, "missed": 0, "total": 1, "percent": 100.0 }
      }
    }
  },
  "totals": {
    "lines": { "covered": 24, "missed": 3, "total": 27, "percent": 88.89 },
    "branches": { "covered": 9, "missed": 3, "total": 12, "percent": 75.0 },
    "functions": { "covered": 5, "missed": 0, "total": 5, "percent": 100.0 },
    "instructions": { "covered": 45, "missed": 5, "total": 50, "percent": 90.0 }
  }
}
```

## Simplified Format (Minimal)

For simple use cases, a minimal format is supported:

```json
{
  "version": "1.0.0",
  "meta": {
    "format": "testformat-json-coverage",
    "formatVersion": "1.0.0",
    "settings": {
      "lineCoverage": true,
      "branchCoverage": false,
      "functionCoverage": false
    }
  },
  "coverage": {
    "src/example.py": {
      "lines": {
        "1": { "hits": 1 },
        "2": { "hits": 10 },
        "5": { "hits": 0 }
      },
      "summary": {
        "lines": { "covered": 2, "missed": 1, "total": 3, "percent": 66.67 }
      }
    }
  },
  "totals": {
    "lines": { "covered": 2, "missed": 1, "total": 3, "percent": 66.67 }
  }
}
```

## Validation Rules

### Required Fields

- Root: `version`, `meta`, `coverage`, `totals`
- Meta: `format`, `formatVersion`, `settings`
- Settings: At least one coverage type MUST be `true`
- File coverage: `lines` OR `functions` MUST be present
- Summary: MUST include metrics for all enabled coverage types

### Consistency Rules

1. Line numbers MUST be positive integers (as strings)
2. Line numbers MUST be in ascending order (for readability, not required for parsing)
3. `hits` MUST be non-negative integers
4. For counters: `covered + missed = total`
5. `percent` MUST be in range [0, 100]
6. File paths in `coverage` MUST be relative, forward-slash separated
7. `formatVersion` MUST match `version` at root level

### Optional Fields

Most fields are optional to support varying levels of coverage detail:

- Branches: Only include if `settings.branchCoverage = true`
- Functions: Only include if function-level data is available
- Classes: Only include for object-oriented languages
- Instructions/statements: Only include if measured by tool
- Complexity: Only include if calculated

## Backwards Compatibility

This format is designed to be forward-compatible:

1. **Version field**: Allows future format evolution
2. **Optional fields**: New coverage types can be added without breaking existing parsers
3. **Extensibility**: Custom fields can be added under a `"x-"` prefix (e.g., `"x-custom-metric"`)

Parsers SHOULD:

- Ignore unknown fields
- Use the `version` field to handle format changes
- Check `settings` to determine which coverage types to expect

## Reference Implementation

A reference implementation will be provided in the `testformat` Go package:

```go
package coverage

type Report struct {
    Version  string              `json:"version"`
    Meta     Metadata            `json:"meta"`
    Coverage map[string]FileCoverage `json:"coverage"`
    Totals   Summary             `json:"totals"`
}

type Metadata struct {
    Format        string     `json:"format"`
    FormatVersion string     `json:"formatVersion"`
    GeneratedAt   string     `json:"generatedAt,omitempty"`
    GeneratedBy   *Tool      `json:"generatedBy,omitempty"`
    Source        *Source    `json:"source,omitempty"`
    Repository    *Repository `json:"repository,omitempty"`
    Settings      Settings   `json:"settings"`
}

type Settings struct {
    BranchCoverage    bool `json:"branchCoverage"`
    ConditionCoverage bool `json:"conditionCoverage"`
    McdcCoverage      bool `json:"mcdcCoverage"`
    FunctionCoverage  bool `json:"functionCoverage"`
    LineCoverage      bool `json:"lineCoverage"`
}

type FileCoverage struct {
    Lines    map[string]LineCoverage `json:"lines,omitempty"`
    Functions []FunctionCoverage      `json:"functions,omitempty"`
    Classes   []ClassCoverage         `json:"classes,omitempty"`
    Summary   Summary                 `json:"summary"`
}

type LineCoverage struct {
    Hits     int             `json:"hits"`
    Branches *BranchCoverage `json:"branches,omitempty"`
}

type BranchCoverage struct {
    Total   int            `json:"total"`
    Covered int            `json:"covered"`
    Missed  int            `json:"missed"`
    Details []BranchDetail `json:"details,omitempty"`
}

type Counter struct {
    Covered int     `json:"covered,omitempty"`
    Missed  int     `json:"missed,omitempty"`
    Total   int     `json:"total"`
    Percent float64 `json:"percent,omitempty"`
    Value   *int    `json:"value,omitempty"`
}

type Summary struct {
    Lines        *Counter `json:"lines,omitempty"`
    Branches     *Counter `json:"branches,omitempty"`
    Functions    *Counter `json:"functions,omitempty"`
    Methods      *Counter `json:"methods,omitempty"`
    Instructions *Counter `json:"instructions,omitempty"`
    Statements   *Counter `json:"statements,omitempty"`
    Classes      *Counter `json:"classes,omitempty"`
    Complexity   *Counter `json:"complexity,omitempty"`
}
```

## Future Considerations

### Potential Extensions

1. **Source code snippets**: Include source code context for uncovered lines
2. **Annotations**: Support for custom annotations/messages per line (from Codecov)
3. **Test mapping**: Link coverage to specific tests that executed each line
4. **Time-series**: Support for tracking coverage over time
5. **Differential coverage**: Support for showing coverage changes vs. base branch
6. **MC/DC coverage**: Full support for Modified Condition/Decision Coverage
7. **Path coverage**: Support for path coverage metrics
8. **Mutation testing**: Integration with mutation testing results

### Open Questions

1. Should we support compressed/gzipped JSON for large reports?
2. Should we define a JSON Schema file for validation?
3. Should we support streaming/chunked parsing for very large reports?
4. Should line numbers be integers or strings? (Currently strings for JSON object keys)

## References

### Coverage Formats Analyzed

- **JaCoCo XML**: https://www.jacoco.org/jacoco/trunk/doc/
- **Cobertura XML**: https://cobertura.github.io/cobertura/
- **LCOV**: https://ltp.sourceforge.net/coverage/lcov.php
- **Coverage.py JSON**: https://coverage.readthedocs.io/
- **Istanbul Coverage**: https://istanbul.js.org/
- **Codecov Custom Format**: https://docs.codecov.com/docs/codecov-custom-coverage-format
- **Clover XML**: https://openclover.org/
- **Go Coverage**: https://go.dev/blog/cover

### Related Standards

- **JSON Schema**: https://json-schema.org/
- **ISO/IEC/IEEE 29119 Software Testing**: Coverage criteria standards
- **Semantic Versioning**: https://semver.org/

## Copyright

This document is placed in the public domain.

## Changelog

- **2025-10-26**: Initial draft (v1.0.0)
