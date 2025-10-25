# CTRF Mapping Analysis

## Fields Not Currently Mapped to CTRF

This document identifies JUnit and Surefire fields that don't have direct mappings in the CTRF format.

---

## High Priority - Significant Data Loss

### 1. **SystemOut / SystemErr** (Test Output)
**Source:** JUnit.TestCase, Surefire.TestCase
**Current Status:** ❌ Not mapped
**Impact:** High - Debugging information lost

```go
// JUnit/Surefire have:
SystemOut  string  // stdout from test execution
SystemErr  string  // stderr from test execution
```

**Use Case:** Critical for debugging failed tests. Often contains print statements, logs, and diagnostic output.

**Recommendation:** Map to CTRF `Extra` field or request CTRF schema extension.

---

### 2. **Properties** (Test Metadata)
**Source:** JUnit.TestSuite, JUnit.TestCase, Surefire.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** High - Contextual metadata lost

```go
// JUnit/Surefire have:
Properties *Properties  // Key-value pairs with test context

type Properties struct {
    Properties []Property
}

type Property struct {
    Name  string
    Value string
    Data  string
}
```

**Use Case:** Contains important metadata like:
- Environment variables
- Configuration settings
- Test parameters
- Custom annotations

**Recommendation:** Map to CTRF `TestResult.Extra` as a properties object.

---

### 3. **Error/Failure Type**
**Source:** JUnit.Failure, JUnit.Error, Surefire.Failure, Surefire.Error
**Current Status:** ⚠️ Partially mapped (only message and trace)
**Impact:** Medium-High - Exception type information lost

```go
// JUnit/Surefire have:
Type string  // e.g., "AssertionError", "NullPointerException"

// Currently only mapping:
Message string  // to CTRF Message
Data    string  // to CTRF Trace
// But NOT Type!
```

**Use Case:** Knowing the exception type helps categorize failures (assertion vs runtime error vs timeout, etc.).

**Recommendation:** Add to CTRF `TestResult.Extra` or map to a structured error object.

---

### 4. **File and Line Number**
**Source:** JUnit.TestCase, JUnit.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Medium - Source location lost

```go
// JUnit has:
File  string  // Source file path
Line  int     // Line number in file
```

**Use Case:** IDEs use this to jump to test definition. Critical for developer experience.

**Recommendation:** Map to CTRF `TestResult.Extra` or request CTRF schema extension.

---

### 5. **Multiple Failures per Test**
**Source:** Surefire.TestCase
**Current Status:** ⚠️ Only first failure mapped
**Impact:** Medium - Additional failure context lost

```go
// Surefire has:
Failures []Failure  // Multiple failures possible
Errors   []Error    // Multiple errors possible

// We only use Failures[0] and Errors[0]
```

**Use Case:** Some test frameworks (especially assertion libraries) can collect multiple assertion failures in a single test run.

**Recommendation:** Map all failures to CTRF `TestResult.Extra.failures` array.

---

## Medium Priority - Useful Context

### 6. **Assertions Count**
**Source:** JUnit.TestCase, JUnit.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Low-Medium - Metric lost

```go
// JUnit has:
Assertions int  // Number of assertions executed
```

**Use Case:** Quality metric - tests with more assertions may be more thorough.

**Recommendation:** Map to CTRF `TestResult.Extra.assertions` or `Summary.Extra.totalAssertions`.

---

### 7. **Flaky/Rerun Metadata**
**Source:** Surefire.TestCase
**Current Status:** ⚠️ Status detected but details lost
**Impact:** Medium - Rerun history lost

```go
// Surefire has:
FlakyFailures  []FlakyFailure   // Failures that later passed
FlakyErrors    []FlakyError     // Errors that later passed
RerunFailures  []RerunFailure   // Failures from reruns
RerunErrors    []RerunError     // Errors from reruns

// Each with:
Message    string
Type       string
StackTrace string
SystemOut  string
SystemErr  string
```

**Current Mapping:** We detect flaky tests (mark as "failed") but don't preserve the flaky/rerun details.

**Use Case:** Understanding test flakiness patterns, how many reruns were needed, what the flaky failures were.

**Recommendation:** Add structured flaky test metadata to CTRF `TestResult.Extra.flaky` with retry history.

---

### 8. **Hostname**
**Source:** JUnit.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Low-Medium - Execution environment lost

```go
// JUnit has:
Hostname string  // Host that ran the tests
```

**Use Case:** Important in distributed test environments (CI/CD, parallel execution).

**Recommendation:** Map to CTRF `Environment` or `Extra`.

---

### 9. **Suite ID and Package**
**Source:** JUnit.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Low - Organizational context lost

```go
// JUnit has:
ID      string  // Unique suite identifier
Package string  // Java package name
```

**Use Case:** Organizational structure, filtering tests by package.

**Recommendation:** Map to CTRF `TestResult.Extra.package` and suite-level metadata.

---

### 10. **Group/Category**
**Source:** Surefire.TestCase, Surefire.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Medium - Test organization lost

```go
// Surefire has:
Group string  // Test group/category (e.g., "integration", "smoke")
```

**Use Case:** Test categorization, filtering by test type, selective execution.

**Recommendation:** Map to CTRF `TestResult.Extra.group` or tags array.

---

### 11. **Timestamp**
**Source:** JUnit.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Low - Test execution time lost

```go
// JUnit has:
Timestamp string  // When suite started (ISO 8601)
```

**Current Mapping:** We set `Report.Timestamp` to current time, but don't use the suite's original timestamp.

**Use Case:** Accurate timing of when tests ran.

**Recommendation:** Parse and use JUnit timestamp for `Report.Timestamp` and `Summary.Start`.

---

### 12. **Surefire Version**
**Source:** Surefire.TestSuite
**Current Status:** ❌ Not mapped
**Impact:** Low - Tool version lost

```go
// Surefire has:
Version string  // Surefire plugin version
```

**Use Case:** Tracking which version of the test runner was used.

**Recommendation:** Map to CTRF `Tool.Version`.

---

## Low Priority - Edge Cases

### 13. **Nested Test Suites**
**Source:** JUnit.TestSuite
**Current Status:** ⚠️ Flattened
**Impact:** Low - Hierarchy lost

```go
// JUnit has:
Suites []TestSuite  // Nested suites
```

**Current Mapping:** We flatten all tests into a single array, losing suite hierarchy.

**Use Case:** Organizational structure, hierarchical reporting.

**Recommendation:** Accept this tradeoff (CTRF uses flat structure) or encode hierarchy in suite names.

---

## Summary Table

| Field | Source | Priority | Current Status | Recommendation |
|-------|--------|----------|----------------|----------------|
| SystemOut/Err | JUnit, Surefire | High | ❌ Not mapped | Map to `Extra.systemOut/systemErr` |
| Properties | JUnit, Surefire | High | ❌ Not mapped | Map to `Extra.properties` |
| Error Type | JUnit, Surefire | High | ⚠️ Partial | Map to `Extra.errorType` |
| File/Line | JUnit | Medium | ❌ Not mapped | Map to `Extra.file/line` |
| Multiple Failures | Surefire | Medium | ⚠️ First only | Map to `Extra.failures` array |
| Assertions | JUnit | Medium | ❌ Not mapped | Map to `Extra.assertions` |
| Flaky Metadata | Surefire | Medium | ⚠️ Detected | Map to `Extra.flaky` with details |
| Hostname | JUnit | Medium | ❌ Not mapped | Map to `Environment` |
| Group | Surefire | Medium | ❌ Not mapped | Map to `Extra.group` |
| Package | JUnit | Low | ❌ Not mapped | Map to `Extra.package` |
| Suite ID | JUnit | Low | ❌ Not mapped | Map to suite metadata |
| Timestamp | JUnit | Low | ⚠️ Overwritten | Use original timestamp |
| Version | Surefire | Low | ❌ Not mapped | Map to `Tool.Version` |
| Nested Suites | JUnit | Low | ⚠️ Flattened | Accept or encode in names |

---

## Recommendations

### Option 1: Use Extra Fields (Immediate)
Extend the conversion functions to populate `TestResult.Extra` with unmapped fields:

```go
test.Extra = map[string]interface{}{
    "systemOut": tc.SystemOut,
    "systemErr": tc.SystemErr,
    "properties": tc.Properties,
    "errorType": tc.Failure.Type,
    "file": tc.File,
    "line": tc.Line,
    "assertions": tc.Assertions,
}
```

**Pros:**
- Preserves all data
- No CTRF schema changes needed
- Backward compatible

**Cons:**
- Non-standard structure
- Not all CTRF consumers will understand these fields

### Option 2: Extend CTRF Types (Better long-term)
Add fields to the CTRF TestResult type that align with the CTRF schema specification:

```go
type TestResult struct {
    Name     string      `json:"name"`
    Status   string      `json:"status"`
    Duration int64       `json:"duration"`
    Suite    string      `json:"suite,omitempty"`
    Message  string      `json:"message,omitempty"`
    Trace    string      `json:"trace,omitempty"`

    // Extended fields (check CTRF spec for what's already supported)
    File     string      `json:"file,omitempty"`
    Line     int         `json:"line,omitempty"`
    Tags     []string    `json:"tags,omitempty"`  // For groups/categories

    Extra    interface{} `json:"extra,omitempty"`
}
```

Check the actual CTRF schema at https://github.com/ctrf-io/ctrf to see what fields are already defined.

### Option 3: Hybrid Approach
- Map commonly-used fields (systemOut, properties, file/line) to standard CTRF fields if they exist
- Use Extra for truly custom/rare fields
- Document what's in Extra
