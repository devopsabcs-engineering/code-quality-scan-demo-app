#!/usr/bin/env python3
"""
lizard-to-sarif.py — Converts Lizard CSV output to SARIF v2.1.0

Usage:
    python lizard-to-sarif.py --input <lizard-csv> --output <sarif-file> [--ccn-threshold 10] [--length-threshold 50]

Lizard CSV columns:
    NLOC, CCN, token, PARAM, length, location, file, function, long_name, start, end
"""

import argparse
import csv
import hashlib
import json
import sys
from pathlib import Path


def parse_args():
    parser = argparse.ArgumentParser(
        description="Convert Lizard CSV output to SARIF v2.1.0"
    )
    parser.add_argument("--input", required=True, help="Path to Lizard CSV output file")
    parser.add_argument("--output", required=True, help="Path for SARIF output file")
    parser.add_argument(
        "--ccn-threshold",
        type=int,
        default=10,
        help="Cyclomatic complexity threshold (default: 10)",
    )
    parser.add_argument(
        "--length-threshold",
        type=int,
        default=50,
        help="Function length threshold in lines (default: 50)",
    )
    return parser.parse_args()


def compute_fingerprint(rule_id: str, file_path: str, function_name: str) -> str:
    """Compute a partial fingerprint for deduplication across runs."""
    raw = f"{rule_id}:{file_path}:{function_name}"
    return hashlib.sha256(raw.encode("utf-8")).hexdigest()[:32]


def get_severity_level(ccn: int, length: int, ccn_threshold: int, length_threshold: int) -> str:
    """Map complexity/length values to SARIF severity levels."""
    if ccn > ccn_threshold * 2:
        return "error"  # CRITICAL: CCN > 2x threshold
    elif ccn > ccn_threshold:
        return "warning"  # HIGH: CCN exceeds threshold
    elif length > length_threshold * 2:
        return "warning"  # HIGH: Very long function
    elif length > length_threshold:
        return "note"  # MEDIUM: Long function
    return "note"


def create_sarif_result(row: dict, ccn_threshold: int, length_threshold: int) -> list:
    """Create SARIF result entries for a single Lizard CSV row."""
    results = []

    try:
        ccn = int(row.get("CCN", 0))
        nloc = int(row.get("NLOC", 0))
        length = int(row.get("length", 0))
        start_line = int(row.get("start", 1))
        end_line = int(row.get("end", start_line))
    except (ValueError, TypeError):
        return results

    file_path = row.get("file", "unknown")
    function_name = row.get("function", "unknown")
    long_name = row.get("long_name", function_name)

    # Normalize file path separators
    file_path = file_path.replace("\\", "/")

    if ccn > ccn_threshold:
        result = {
            "ruleId": "ccn-exceeded",
            "level": get_severity_level(ccn, length, ccn_threshold, length_threshold),
            "message": {
                "text": f"Function '{long_name}' has cyclomatic complexity {ccn} (threshold: {ccn_threshold}). "
                f"Consider refactoring into smaller functions using early returns and method extraction."
            },
            "locations": [
                {
                    "physicalLocation": {
                        "artifactLocation": {"uri": file_path},
                        "region": {
                            "startLine": start_line,
                            "endLine": end_line,
                        },
                    }
                }
            ],
            "partialFingerprints": {
                "primaryLocationLineHash": compute_fingerprint(
                    "ccn-exceeded", file_path, function_name
                )
            },
            "properties": {
                "precision": "very-high",
                "tags": ["code-quality", "complexity"],
                "metrics": {
                    "ccn": ccn,
                    "nloc": nloc,
                    "length": length,
                    "threshold": ccn_threshold,
                },
            },
        }
        results.append(result)

    if length > length_threshold:
        result = {
            "ruleId": "function-length-exceeded",
            "level": "note" if length <= length_threshold * 2 else "warning",
            "message": {
                "text": f"Function '{long_name}' has {length} lines (threshold: {length_threshold}). "
                f"Consider breaking into smaller, focused functions."
            },
            "locations": [
                {
                    "physicalLocation": {
                        "artifactLocation": {"uri": file_path},
                        "region": {
                            "startLine": start_line,
                            "endLine": end_line,
                        },
                    }
                }
            ],
            "partialFingerprints": {
                "primaryLocationLineHash": compute_fingerprint(
                    "function-length-exceeded", file_path, function_name
                )
            },
            "properties": {
                "precision": "very-high",
                "tags": ["code-quality", "complexity"],
                "metrics": {
                    "length": length,
                    "threshold": length_threshold,
                },
            },
        }
        results.append(result)

    return results


def convert(input_path: str, output_path: str, ccn_threshold: int, length_threshold: int):
    """Convert Lizard CSV to SARIF v2.1.0."""
    results = []

    with open(input_path, "r", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        for row in reader:
            results.extend(
                create_sarif_result(row, ccn_threshold, length_threshold)
            )

    sarif = {
        "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
        "version": "2.1.0",
        "runs": [
            {
                "tool": {
                    "driver": {
                        "name": "Lizard",
                        "version": "1.17.10",
                        "informationUri": "https://github.com/terryyin/lizard",
                        "rules": [
                            {
                                "id": "ccn-exceeded",
                                "name": "CyclomaticComplexityExceeded",
                                "shortDescription": {
                                    "text": "Function cyclomatic complexity exceeds threshold"
                                },
                                "fullDescription": {
                                    "text": "The function's cyclomatic complexity (CCN) exceeds the configured threshold. "
                                    "High complexity indicates code that is difficult to test and maintain."
                                },
                                "help": {
                                    "text": "Reduce complexity by extracting helper functions, using early returns, "
                                    "and simplifying conditional logic.",
                                    "markdown": "## Cyclomatic Complexity Exceeded\n\n"
                                    "The function's cyclomatic complexity (CCN) exceeds the threshold.\n\n"
                                    "### Remediation\n"
                                    "- Extract helper functions for distinct logic branches\n"
                                    "- Use early returns (guard clauses) to reduce nesting\n"
                                    "- Replace switch/case with strategy pattern or lookup table\n\n"
                                    "### References\n"
                                    "- [CWE-1121: Excessive McCabe Cyclomatic Complexity](https://cwe.mitre.org/data/definitions/1121.html)",
                                },
                                "properties": {
                                    "tags": ["code-quality", "complexity", "maintainability"],
                                    "cwe": ["CWE-1121"],
                                },
                            },
                            {
                                "id": "function-length-exceeded",
                                "name": "FunctionLengthExceeded",
                                "shortDescription": {
                                    "text": "Function length exceeds threshold"
                                },
                                "fullDescription": {
                                    "text": "The function body exceeds the recommended line count. "
                                    "Long functions are harder to understand, test, and maintain."
                                },
                                "help": {
                                    "text": "Break the function into smaller, focused functions with descriptive names.",
                                    "markdown": "## Function Length Exceeded\n\n"
                                    "The function exceeds the recommended line count.\n\n"
                                    "### Remediation\n"
                                    "- Extract logical sections into named helper functions\n"
                                    "- Apply Single Responsibility Principle\n"
                                    "- Consider using composition over procedural steps\n",
                                },
                                "properties": {
                                    "tags": ["code-quality", "complexity", "maintainability"],
                                },
                            },
                        ],
                    }
                },
                "automationDetails": {"id": "code-quality/complexity/"},
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

    count = convert(args.input, args.output, args.ccn_threshold, args.length_threshold)
    print(f"  CCN threshold: {args.ccn_threshold}")
    print(f"  Length threshold: {args.length_threshold}")
    print(f"  Total findings: {count}")


if __name__ == "__main__":
    main()
