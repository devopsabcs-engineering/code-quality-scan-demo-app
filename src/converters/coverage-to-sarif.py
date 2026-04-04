#!/usr/bin/env python3
"""
coverage-to-sarif.py — Converts coverage reports to SARIF v2.1.0

Usage:
    python coverage-to-sarif.py --input <coverage-file> --format <format> --output <sarif-file> [--threshold 80]

Supported formats:
    cobertura    - Cobertura XML (pytest-cov, Coverlet)
    json-summary - Istanbul JSON summary (jest, vitest)
    lcov         - LCOV info format
    jacoco       - JaCoCo XML
    gocover      - Go coverage profile (go test -coverprofile)
"""

import argparse
import hashlib
import json
import sys
import xml.etree.ElementTree as ET
from pathlib import Path


def parse_args():
    parser = argparse.ArgumentParser(
        description="Convert coverage reports to SARIF v2.1.0"
    )
    parser.add_argument("--input", required=True, help="Path to coverage report file")
    parser.add_argument(
        "--format",
        required=True,
        choices=["cobertura", "json-summary", "lcov", "jacoco", "gocover"],
        help="Coverage report format",
    )
    parser.add_argument("--output", required=True, help="Path for SARIF output file")
    parser.add_argument(
        "--threshold",
        type=float,
        default=80.0,
        help="Coverage threshold percentage (default: 80)",
    )
    return parser.parse_args()


def compute_fingerprint(rule_id: str, file_path: str) -> str:
    """Compute a partial fingerprint for deduplication."""
    raw = f"{rule_id}:{file_path}"
    return hashlib.sha256(raw.encode("utf-8")).hexdigest()[:32]


def get_severity_level(coverage_pct: float) -> str:
    """Map coverage percentage to SARIF severity level."""
    if coverage_pct < 50:
        return "error"  # CRITICAL
    elif coverage_pct < 80:
        return "warning"  # HIGH
    elif coverage_pct < 90:
        return "note"  # MEDIUM
    return "none"  # Pass — no finding


def parse_cobertura(input_path: str) -> list:
    """Parse Cobertura XML coverage report."""
    tree = ET.parse(input_path)
    root = tree.getroot()
    files = []

    for package in root.findall(".//package"):
        for cls in package.findall(".//class"):
            filename = cls.get("filename", "unknown")
            line_rate = float(cls.get("line-rate", "0")) * 100
            branch_rate = float(cls.get("branch-rate", "0")) * 100
            files.append(
                {
                    "file": filename,
                    "line_coverage": line_rate,
                    "branch_coverage": branch_rate,
                }
            )

    return files


def parse_json_summary(input_path: str) -> list:
    """Parse Istanbul/Jest JSON summary coverage report."""
    with open(input_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    files = []
    for file_path, metrics in data.items():
        if file_path == "total":
            continue
        lines = metrics.get("lines", {})
        branches = metrics.get("branches", {})
        functions = metrics.get("functions", {})
        files.append(
            {
                "file": file_path,
                "line_coverage": lines.get("pct", 0),
                "branch_coverage": branches.get("pct", 0),
                "function_coverage": functions.get("pct", 0),
            }
        )

    return files


def parse_lcov(input_path: str) -> list:
    """Parse LCOV info format coverage report."""
    files = []
    current_file = None
    lines_found = 0
    lines_hit = 0

    with open(input_path, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if line.startswith("SF:"):
                current_file = line[3:]
                lines_found = 0
                lines_hit = 0
            elif line.startswith("LF:"):
                lines_found = int(line[3:])
            elif line.startswith("LH:"):
                lines_hit = int(line[3:])
            elif line == "end_of_record" and current_file:
                pct = (lines_hit / lines_found * 100) if lines_found > 0 else 100
                files.append(
                    {
                        "file": current_file,
                        "line_coverage": pct,
                        "branch_coverage": 0,
                    }
                )
                current_file = None

    return files


def parse_jacoco(input_path: str) -> list:
    """Parse JaCoCo XML coverage report."""
    tree = ET.parse(input_path)
    root = tree.getroot()
    files = []

    for package in root.findall(".//package"):
        pkg_name = package.get("name", "").replace("/", ".")
        for source in package.findall("sourcefile"):
            filename = source.get("name", "unknown")
            full_path = f"{pkg_name}/{filename}" if pkg_name else filename

            line_counter = source.find("counter[@type='LINE']")
            branch_counter = source.find("counter[@type='BRANCH']")

            if line_counter is not None:
                missed = int(line_counter.get("missed", "0"))
                covered = int(line_counter.get("covered", "0"))
                total = missed + covered
                line_pct = (covered / total * 100) if total > 0 else 100
            else:
                line_pct = 100

            if branch_counter is not None:
                missed = int(branch_counter.get("missed", "0"))
                covered = int(branch_counter.get("covered", "0"))
                total = missed + covered
                branch_pct = (covered / total * 100) if total > 0 else 100
            else:
                branch_pct = 100

            files.append(
                {
                    "file": full_path,
                    "line_coverage": line_pct,
                    "branch_coverage": branch_pct,
                }
            )

    return files


def parse_gocover(input_path: str) -> list:
    """Parse Go coverage profile."""
    file_stats = {}

    with open(input_path, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if line.startswith("mode:"):
                continue
            # Format: file:startLine.startCol,endLine.endCol numStatements count
            parts = line.split()
            if len(parts) < 3:
                continue
            location = parts[0]
            num_statements = int(parts[1])
            count = int(parts[2])

            file_path = location.split(":")[0]
            if file_path not in file_stats:
                file_stats[file_path] = {"total": 0, "covered": 0}
            file_stats[file_path]["total"] += num_statements
            if count > 0:
                file_stats[file_path]["covered"] += num_statements

    files = []
    for file_path, stats in file_stats.items():
        pct = (stats["covered"] / stats["total"] * 100) if stats["total"] > 0 else 100
        files.append({"file": file_path, "line_coverage": pct, "branch_coverage": 0})

    return files


def convert(input_path: str, fmt: str, output_path: str, threshold: float):
    """Convert coverage report to SARIF v2.1.0."""
    parsers = {
        "cobertura": parse_cobertura,
        "json-summary": parse_json_summary,
        "lcov": parse_lcov,
        "jacoco": parse_jacoco,
        "gocover": parse_gocover,
    }

    parser_func = parsers[fmt]
    file_coverages = parser_func(input_path)

    results = []
    for entry in file_coverages:
        file_path = entry["file"].replace("\\", "/")
        line_cov = entry.get("line_coverage", 100)
        branch_cov = entry.get("branch_coverage", 100)

        level = get_severity_level(line_cov)
        if level == "none":
            continue  # Above threshold, skip

        result = {
            "ruleId": "coverage-below-threshold",
            "level": level,
            "message": {
                "text": f"File '{file_path}' has {line_cov:.1f}% line coverage "
                f"(threshold: {threshold}%). "
                f"Branch coverage: {branch_cov:.1f}%. "
                f"Add unit tests to improve coverage."
            },
            "locations": [
                {
                    "physicalLocation": {
                        "artifactLocation": {"uri": file_path},
                        "region": {"startLine": 1},
                    }
                }
            ],
            "partialFingerprints": {
                "primaryLocationLineHash": compute_fingerprint(
                    "coverage-below-threshold", file_path
                )
            },
            "properties": {
                "precision": "very-high",
                "tags": ["code-quality", "coverage"],
                "metrics": {
                    "lineCoverage": round(line_cov, 2),
                    "branchCoverage": round(branch_cov, 2),
                    "threshold": threshold,
                },
            },
        }
        results.append(result)

    sarif = {
        "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
        "version": "2.1.0",
        "runs": [
            {
                "tool": {
                    "driver": {
                        "name": "coverage-to-sarif",
                        "version": "1.0.0",
                        "informationUri": "https://github.com/devopsabcs-engineering/code-quality-scan-demo-app",
                        "rules": [
                            {
                                "id": "coverage-below-threshold",
                                "name": "CoverageBelowThreshold",
                                "shortDescription": {
                                    "text": "File test coverage is below the required threshold"
                                },
                                "fullDescription": {
                                    "text": "The file's test coverage percentage is below the configured threshold. "
                                    "Low coverage increases the risk of undetected bugs and regressions."
                                },
                                "help": {
                                    "text": "Add unit tests for uncovered functions and branches.",
                                    "markdown": "## Coverage Below Threshold\n\n"
                                    "The file's test coverage is below the required minimum.\n\n"
                                    "### Remediation\n"
                                    "- Add unit tests for all public functions\n"
                                    "- Test error paths and boundary conditions\n"
                                    "- Use code coverage reports to identify uncovered lines\n\n"
                                    "### References\n"
                                    "- [CWE-754: Improper Check for Unusual Conditions](https://cwe.mitre.org/data/definitions/754.html)",
                                },
                                "properties": {
                                    "tags": ["code-quality", "coverage", "testing"],
                                    "cwe": ["CWE-754"],
                                },
                            }
                        ],
                    }
                },
                "automationDetails": {"id": "code-quality/coverage/"},
                "results": results,
            }
        ],
    }

    output = Path(output_path)
    output.parent.mkdir(parents=True, exist_ok=True)
    with open(output, "w", encoding="utf-8") as f:
        json.dump(sarif, f, indent=2)

    print(f"Converted {len(results)} findings to SARIF: {output_path}")
    return len(results)


def main():
    args = parse_args()

    if not Path(args.input).exists():
        print(f"Error: Input file not found: {args.input}", file=sys.stderr)
        sys.exit(1)

    count = convert(args.input, args.format, args.output, args.threshold)
    print(f"  Format: {args.format}")
    print(f"  Threshold: {args.threshold}%")
    print(f"  Files below threshold: {count}")


if __name__ == "__main__":
    main()
