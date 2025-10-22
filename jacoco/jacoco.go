package jacoco

import "encoding/xml"

// Report represents the root element of a JaCoCo coverage report
type Report struct {
	XMLName      xml.Name      `xml:"report" json:"-"`
	Name         string        `xml:"name,attr" json:"name"`
	SessionInfos []SessionInfo `xml:"sessioninfo,omitempty" json:"sessionInfo,omitempty"`
	Groups       []Group       `xml:"group,omitempty" json:"groups,omitempty"`
	Packages     []Package     `xml:"package,omitempty" json:"packages,omitempty"`
	Counters     []Counter     `xml:"counter,omitempty" json:"counters,omitempty"`
}

// SessionInfo represents information about a session which contributed execution data
type SessionInfo struct {
	ID    string `xml:"id,attr" json:"id"`
	Start string `xml:"start,attr" json:"start"`
	Dump  string `xml:"dump,attr" json:"dump"`
}

// Group represents a logical grouping of packages
type Group struct {
	Name     string    `xml:"name,attr" json:"name"`
	Groups   []Group   `xml:"group,omitempty" json:"groups,omitempty"`
	Packages []Package `xml:"package,omitempty" json:"packages,omitempty"`
	Counters []Counter `xml:"counter,omitempty" json:"counters,omitempty"`
}

// Package represents a Java package with classes and source files
type Package struct {
	Name        string       `xml:"name,attr" json:"name"`
	Classes     []Class      `xml:"class,omitempty" json:"classes,omitempty"`
	SourceFiles []SourceFile `xml:"sourcefile,omitempty" json:"sourceFiles,omitempty"`
	Counters    []Counter    `xml:"counter,omitempty" json:"counters,omitempty"`
}

// Class represents a Java class with coverage information
type Class struct {
	Name           string    `xml:"name,attr" json:"name"`
	SourceFileName string    `xml:"sourcefilename,attr,omitempty" json:"sourceFileName,omitempty"`
	Methods        []Method  `xml:"method,omitempty" json:"methods,omitempty"`
	Counters       []Counter `xml:"counter,omitempty" json:"counters,omitempty"`
}

// Method represents a method within a class
type Method struct {
	Name     string    `xml:"name,attr" json:"name"`
	Desc     string    `xml:"desc,attr" json:"desc"`
	Line     int       `xml:"line,attr,omitempty" json:"line,omitempty"`
	Counters []Counter `xml:"counter,omitempty" json:"counters,omitempty"`
}

// SourceFile represents a source file with line-level coverage
type SourceFile struct {
	Name     string    `xml:"name,attr" json:"name"`
	Lines    []Line    `xml:"line,omitempty" json:"lines,omitempty"`
	Counters []Counter `xml:"counter,omitempty" json:"counters,omitempty"`
}

// Line represents coverage information for a single line of source code
type Line struct {
	Nr              int `xml:"nr,attr" json:"nr"`
	MissedInstr     int `xml:"mi,attr,omitempty" json:"mi,omitempty"`
	CoveredInstr    int `xml:"ci,attr,omitempty" json:"ci,omitempty"`
	MissedBranches  int `xml:"mb,attr,omitempty" json:"mb,omitempty"`
	CoveredBranches int `xml:"cb,attr,omitempty" json:"cb,omitempty"`
}

// Counter represents coverage metrics for different types
type Counter struct {
	Type    string `xml:"type,attr" json:"type"`
	Missed  int    `xml:"missed,attr" json:"missed"`
	Covered int    `xml:"covered,attr" json:"covered"`
}
