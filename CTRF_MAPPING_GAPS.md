# CTRF Mapping Gaps - What We Should Fix

## Executive Summary

The CTRF schema actually supports **many more fields** than our current implementation uses. We're not losing as much data as initially thought - we just need to update our TestResult type to use the standard CTRF fields.

---

## Standard CTRF Fields We're NOT Using (But Should!)

### 1. **stdout / stderr** ✅ Standard CTRF Fields
**CTRF Has:** `stdout` (array of strings), `stderr` (array of strings)
**We Currently:** ❌ Not mapping SystemOut/SystemErr
**Impact:** HIGH

```go
// Should map:
JUnit.TestCase.SystemOut -> CTRF.TestResult.Stdout (split by lines)
JUnit.TestCase.SystemErr -> CTRF.TestResult.Stderr (split by lines)
```

---

### 2. **filePath / line** ✅ Standard CTRF Fields
**CTRF Has:** `filePath` (string), `line` (integer)
**We Currently:** ❌ Not mapping
**Impact:** MEDIUM-HIGH (IDE integration)

```go
// Should map:
JUnit.TestCase.File -> CTRF.TestResult.FilePath
JUnit.TestCase.Line -> CTRF.TestResult.Line
```

---

### 3. **type** ✅ Standard CTRF Field
**CTRF Has:** `type` (string)
**We Currently:** ❌ Not mapping error/failure type
**Impact:** MEDIUM-HIGH

```go
// Should map:
JUnit.Failure.Type -> CTRF.TestResult.Type  // e.g., "AssertionError"
Surefire.Failure.Type -> CTRF.TestResult.Type
```

---

### 4. **tags** ✅ Standard CTRF Field
**CTRF Has:** `tags` (array of strings)
**We Currently:** ❌ Not mapping groups
**Impact:** MEDIUM

```go
// Should map:
Surefire.TestCase.Group -> CTRF.TestResult.Tags (as array)
```

---

### 5. **flaky** ✅ Standard CTRF Field
**CTRF Has:** `flaky` (boolean)
**We Currently:** ⚠️ We detect flaky but don't set this flag
**Impact:** MEDIUM

```go
// Should set:
if len(tc.FlakyFailures) > 0 || len(tc.FlakyErrors) > 0 {
    test.Flaky = true
}
```

---

### 6. **retries / retryAttempts** ✅ Standard CTRF Fields
**CTRF Has:**
- `retries` (integer) - number of retry attempts
- `retryAttempts` (array) - detailed retry information

**We Currently:** ❌ Not mapping Surefire rerun data
**Impact:** MEDIUM

```go
// Should map Surefire rerun data:
retryAttempts: [
  {
    attempt: 1,
    status: "failed",
    duration: 100,
    message: "First attempt failed",
    trace: "...",
  },
  {
    attempt: 2,
    status: "passed",
    duration: 95,
  }
]
```

---

### 7. **parameters** ✅ Standard CTRF Field
**CTRF Has:** `parameters` (object) - test parameters/properties
**We Currently:** ❌ Not mapping Properties
**Impact:** HIGH

```go
// Should map:
JUnit.Properties -> CTRF.TestResult.Parameters {
  "key1": "value1",
  "key2": "value2"
}
```

---

### 8. **start / stop** ✅ Standard CTRF Fields
**CTRF Has:** `start` (integer timestamp), `stop` (integer timestamp)
**We Currently:** ❌ Not setting per-test timing
**Impact:** LOW (JUnit/Surefire don't provide this)

---

### 9. **rawStatus** ✅ Standard CTRF Field
**CTRF Has:** `rawStatus` (string) - original status from test framework
**We Currently:** ❌ Not preserving original status
**Impact:** LOW-MEDIUM

```go
// Could preserve:
if tc.Skipped != nil {
    test.RawStatus = "skipped"
    test.Status = "skipped"
} else if tc.Failure != nil {
    test.RawStatus = "failure"  // vs "error"
    test.Status = "failed"
}
```

---

## Fields JUnit/Surefire Have That CTRF Doesn't

### Assertions Count
**JUnit Has:** `Assertions` (int)
**CTRF Has:** No standard field
**Recommendation:** Use `Extra.assertions`

### Hostname
**JUnit Has:** `Hostname` (string)
**CTRF Has:** Could go in `Environment` but not test-specific
**Recommendation:** Use `Extra.hostname` or populate `Environment`

### Multiple Failures in Single Test
**Surefire Has:** `Failures []Failure` (array)
**CTRF Has:** Only single `message` and `trace`
**Recommendation:**
- Use first failure for message/trace
- Put additional failures in `Extra.additionalFailures`

### Suite Package/ID
**JUnit Has:** `Package`, `ID` at suite level
**CTRF Has:** `suite` (array) at test level, but no suite metadata
**Recommendation:** Encode in suite name or use `Extra`

---

## Updated TestResult Type (Recommended)

```go
type TestResult struct {
    // Required
    Name     string `json:"name"`
    Status   string `json:"status"`  // passed, failed, skipped, pending, other
    Duration int64  `json:"duration"` // milliseconds

    // Optional - standard CTRF fields
    ID            string        `json:"id,omitempty"`             // UUID
    Start         int64         `json:"start,omitempty"`          // timestamp ms
    Stop          int64         `json:"stop,omitempty"`           // timestamp ms
    Suite         string        `json:"suite,omitempty"`          // Changed to string from array for simplicity
    Message       string        `json:"message,omitempty"`
    Trace         string        `json:"trace,omitempty"`
    Snippet       string        `json:"snippet,omitempty"`
    AI            string        `json:"ai,omitempty"`
    Line          int           `json:"line,omitempty"`
    RawStatus     string        `json:"rawStatus,omitempty"`
    Tags          []string      `json:"tags,omitempty"`
    Type          string        `json:"type,omitempty"`           // Error/failure type
    FilePath      string        `json:"filePath,omitempty"`
    Retries       int           `json:"retries,omitempty"`
    RetryAttempts []RetryAttempt `json:"retryAttempts,omitempty"`
    Flaky         bool          `json:"flaky,omitempty"`
    Stdout        []string      `json:"stdout,omitempty"`
    Stderr        []string      `json:"stderr,omitempty"`
    ThreadID      string        `json:"threadId,omitempty"`
    Browser       string        `json:"browser,omitempty"`
    Device        string        `json:"device,omitempty"`
    Screenshot    string        `json:"screenshot,omitempty"`
    Attachments   []Attachment  `json:"attachments,omitempty"`
    Parameters    interface{}   `json:"parameters,omitempty"`    // For JUnit Properties
    Steps         []Step        `json:"steps,omitempty"`
    Insights      *Insights     `json:"insights,omitempty"`
    Extra         interface{}   `json:"extra,omitempty"`
}

type RetryAttempt struct {
    Attempt     int           `json:"attempt"`
    Status      string        `json:"status"`
    Duration    int64         `json:"duration,omitempty"`
    Message     string        `json:"message,omitempty"`
    Trace       string        `json:"trace,omitempty"`
    Stdout      []string      `json:"stdout,omitempty"`
    Stderr      []string      `json:"stderr,omitempty"`
    Attachments []Attachment  `json:"attachments,omitempty"`
    Extra       interface{}   `json:"extra,omitempty"`
}

type Attachment struct {
    Name        string `json:"name"`
    ContentType string `json:"contentType,omitempty"`
    Path        string `json:"path"`
}

type Step struct {
    Name   string      `json:"name"`
    Status string      `json:"status"`
    Extra  interface{} `json:"extra,omitempty"`
}

type Insights struct {
    PassRate  float64 `json:"passRate,omitempty"`
    FailRate  float64 `json:"failRate,omitempty"`
    FlakyRate float64 `json:"flakyRate,omitempty"`
    Extra     interface{} `json:"extra,omitempty"`
}
```

---

## Priority Action Items

### 🔴 High Priority (Should fix immediately)
1. **Add stdout/stderr arrays** - Split SystemOut/SystemErr by newlines
2. **Add filePath and line** - Map JUnit File/Line fields
3. **Add type field** - Map error/failure type
4. **Add parameters** - Map JUnit/Surefire Properties
5. **Add flaky boolean** - Set for Surefire flaky tests

### 🟡 Medium Priority (Nice to have)
6. **Add tags array** - Map Surefire Group
7. **Add retryAttempts** - Map Surefire rerun data
8. **Add rawStatus** - Preserve original status string

### 🟢 Low Priority (Future enhancement)
9. Add assertions to Extra
10. Add hostname to Extra or Environment
11. Handle multiple failures in Extra

---

## Comparison: Before vs After

### Before (Current Implementation)
```json
{
  "name": "test_failure",
  "status": "failed",
  "duration": 200,
  "suite": "com.example.Test",
  "message": "Expected 5 but got 3",
  "trace": "AssertionError: Expected 5 but got 3\n    at Test.java:42"
}
```

### After (With Full CTRF Fields)
```json
{
  "name": "test_failure",
  "status": "failed",
  "duration": 200,
  "suite": "com.example.Test",
  "message": "Expected 5 but got 3",
  "trace": "AssertionError: Expected 5 but got 3\n    at Test.java:42",
  "type": "AssertionError",
  "filePath": "src/test/java/com/example/Test.java",
  "line": 42,
  "stdout": ["Running test...", "Test completed"],
  "stderr": [],
  "parameters": {
    "testParam1": "value1",
    "environment": "staging"
  },
  "extra": {
    "assertions": 5,
    "hostname": "ci-runner-3"
  }
}
```

---

## Conclusion

**Good News:** CTRF already supports most of what we need! We just need to:
1. Extend our `TestResult` type to include standard CTRF fields
2. Update conversion functions to populate these fields
3. Use `Extra` only for truly non-standard data (assertions, hostname)

This will result in much richer, more useful CTRF reports while maintaining compatibility with the standard schema.
