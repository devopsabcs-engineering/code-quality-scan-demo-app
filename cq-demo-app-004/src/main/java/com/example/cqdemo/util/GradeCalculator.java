package com.example.cqdemo.util;

import com.example.cqdemo.model.Student;
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import java.util.ArrayList;

// =========================================================================
// VIOLATION: Missing Javadoc on ALL public methods throughout this class.
// VIOLATION: Magic numbers used extensively instead of named constants.
// VIOLATION: Methods exceed 50 lines.
// VIOLATION: System.out.println instead of proper logging framework.
// =========================================================================
public class GradeCalculator {

    // These SHOULD be named constants but are used as magic numbers inline
    // private static final double HOMEWORK_WEIGHT = 0.3;
    // private static final double EXAM_WEIGHT = 0.7;

    public static double calculateWeightedGrade(double homeworkScore, double examScore) {
        // VIOLATION: magic numbers 0.3 and 0.7 instead of named constants
        double weighted = homeworkScore * 0.3 + examScore * 0.7;
        System.out.println("Calculated weighted grade: " + weighted);
        return weighted;
    }

    public static String getLetterGrade(double numericGrade) {
        // VIOLATION: magic numbers 90, 80, 70, 65, 60 instead of named constants
        if (numericGrade >= 90) {
            return "A";
        } else if (numericGrade >= 80) {
            return "B";
        } else if (numericGrade >= 70) {
            return "C";
        } else if (numericGrade >= 65) {
            return "C-";
        } else if (numericGrade >= 60) {
            return "D";
        } else {
            return "F";
        }
    }

    public static double getGradePoints(String letterGrade) {
        // VIOLATION: magic numbers for GPA points
        if (letterGrade == null) return 0.0;
        switch (letterGrade) {
            case "A": return 4.0;
            case "A-": return 3.7;
            case "B+": return 3.3;
            case "B": return 3.0;
            case "B-": return 2.7;
            case "C+": return 2.3;
            case "C": return 2.0;
            case "C-": return 1.7;
            case "D+": return 1.3;
            case "D": return 1.0;
            case "F": return 0.0;
            default: return 0.0;
        }
    }

    // =========================================================================
    // VIOLATION: Method exceeds 50 lines
    // VIOLATION: Magic numbers throughout
    // VIOLATION: Missing Javadoc
    // VIOLATION: System.out.println instead of logger
    // =========================================================================
    public static Map<String, Object> calculateClassStatistics(List<Double> scores) {
        Map<String, Object> stats = new HashMap<>();

        if (scores == null || scores.isEmpty()) {
            stats.put("error", "No scores provided");
            return stats;
        }

        double sum = 0;
        double min = 999999;  // VIOLATION: magic number
        double max = -1;      // VIOLATION: magic number
        int countA = 0;
        int countB = 0;
        int countC = 0;
        int countD = 0;
        int countF = 0;
        int countAbove90 = 0;
        int countAbove80 = 0;
        int countBelow60 = 0;

        for (double score : scores) {
            sum += score;
            if (score < min) min = score;
            if (score > max) max = score;

            // VIOLATION: magic numbers repeated from getLetterGrade
            if (score >= 90) {
                countA++;
                countAbove90++;
                countAbove80++;
            } else if (score >= 80) {
                countB++;
                countAbove80++;
            } else if (score >= 70) {
                countC++;
            } else if (score >= 60) {
                countD++;
            } else {
                countF++;
                countBelow60++;
            }
        }

        double mean = sum / scores.size();

        // Calculate standard deviation
        double varianceSum = 0;
        for (double score : scores) {
            varianceSum += (score - mean) * (score - mean);
        }
        double stdDev = Math.sqrt(varianceSum / scores.size());

        // Calculate median
        List<Double> sorted = new ArrayList<>(scores);
        sorted.sort(Double::compareTo);
        double median;
        int n = sorted.size();
        if (n % 2 == 0) {
            median = (sorted.get(n / 2 - 1) + sorted.get(n / 2)) / 2.0;
        } else {
            median = sorted.get(n / 2);
        }

        stats.put("count", scores.size());
        stats.put("mean", Math.round(mean * 100.0) / 100.0);     // VIOLATION: magic number 100.0
        stats.put("median", Math.round(median * 100.0) / 100.0);  // VIOLATION: magic number 100.0
        stats.put("min", min);
        stats.put("max", max);
        stats.put("stdDev", Math.round(stdDev * 100.0) / 100.0);  // VIOLATION: magic number 100.0
        stats.put("range", max - min);
        stats.put("gradeDistribution", Map.of(
                "A", countA, "B", countB, "C", countC, "D", countD, "F", countF
        ));
        stats.put("passRate", Math.round(((double)(countA + countB + countC + countD) / scores.size()) * 100.0) / 100.0);
        stats.put("above90Percent", countAbove90);
        stats.put("above80Percent", countAbove80);
        stats.put("below60Percent", countBelow60);

        // VIOLATION: System.out.println instead of logger
        System.out.println("Class statistics calculated for " + scores.size() + " students");
        System.out.println("Mean: " + mean + ", Median: " + median + ", StdDev: " + stdDev);

        return stats;
    }

    // =========================================================================
    // VIOLATION: Method exceeds 50 lines
    // VIOLATION: Magic numbers for thresholds
    // VIOLATION: Missing Javadoc
    // VIOLATION: HIGH COMPLEXITY — nested conditionals for GPA calculation
    // =========================================================================
    public static Map<String, Object> calculateGPA(List<Map<String, Object>> courseGrades) {
        Map<String, Object> result = new HashMap<>();

        if (courseGrades == null || courseGrades.isEmpty()) {
            result.put("gpa", 0.0);
            result.put("error", "No course grades provided");
            return result;
        }

        double totalPoints = 0;
        int totalCredits = 0;
        int coursesCompleted = 0;
        int coursesFailed = 0;
        List<Map<String, Object>> breakdown = new ArrayList<>();

        for (Map<String, Object> courseGrade : courseGrades) {
            String grade = (String) courseGrade.get("grade");
            int credits = courseGrade.containsKey("credits")
                    ? ((Number) courseGrade.get("credits")).intValue()
                    : 3; // VIOLATION: magic number default

            if (grade == null || grade.isEmpty()) {
                System.out.println("Skipping course with no grade"); // VIOLATION: sysout
                continue;
            }

            double gradePoints = getGradePoints(grade);
            double qualityPoints = gradePoints * credits;

            totalPoints += qualityPoints;
            totalCredits += credits;
            coursesCompleted++;

            if (gradePoints < 1.0) { // VIOLATION: magic number
                coursesFailed++;
            }

            Map<String, Object> courseBreakdown = new HashMap<>();
            courseBreakdown.put("course", courseGrade.getOrDefault("courseName", "Unknown"));
            courseBreakdown.put("grade", grade);
            courseBreakdown.put("credits", credits);
            courseBreakdown.put("gradePoints", gradePoints);
            courseBreakdown.put("qualityPoints", qualityPoints);
            breakdown.add(courseBreakdown);
        }

        double gpa = totalCredits > 0 ? totalPoints / totalCredits : 0.0;
        gpa = Math.round(gpa * 100.0) / 100.0; // VIOLATION: magic number 100.0

        result.put("gpa", gpa);
        result.put("totalCredits", totalCredits);
        result.put("totalQualityPoints", Math.round(totalPoints * 100.0) / 100.0);
        result.put("coursesCompleted", coursesCompleted);
        result.put("coursesFailed", coursesFailed);
        result.put("breakdown", breakdown);

        // Academic standing based on GPA — all magic numbers
        if (gpa >= 3.5) {
            result.put("standing", "Dean's List");
        } else if (gpa >= 3.0) {
            result.put("standing", "Good Standing");
        } else if (gpa >= 2.0) {
            result.put("standing", "Satisfactory");
        } else if (gpa >= 1.0) {
            result.put("standing", "Academic Probation");
        } else {
            result.put("standing", "Academic Suspension");
        }

        // VIOLATION: System.out.println instead of logger
        System.out.println("GPA calculated: " + gpa + " (" + result.get("standing") + ")");
        System.out.println("Credits: " + totalCredits + ", Failed: " + coursesFailed);

        return result;
    }

    // VIOLATION: magic number for curve adjustment
    public static double applyCurve(double score, double curveAmount) {
        double curved = score + curveAmount;
        if (curved > 100) { // VIOLATION: magic number
            curved = 100;   // VIOLATION: magic number
        }
        System.out.println("Score " + score + " curved to " + curved); // VIOLATION: sysout
        return curved;
    }
}
