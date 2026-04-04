package com.example.cqdemo.util;

import com.example.cqdemo.model.Student;
import com.example.cqdemo.model.Course;
import com.example.cqdemo.model.Enrollment;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.Map;
import java.util.ArrayList;
import java.util.HashMap;

// =========================================================================
// VIOLATION: Missing Javadoc on class
// VIOLATION: String concatenation in loops instead of StringBuilder
// VIOLATION: Unused private methods
// VIOLATION: Raw types (List instead of List<Student>)
// VIOLATION: System.out.println instead of proper logging framework
// =========================================================================
public class ReportFormatter {

    // VIOLATION: unused private method
    private static String padRight(String text, int length) {
        if (text == null) text = "";
        if (text.length() >= length) return text.substring(0, length);
        StringBuilder sb = new StringBuilder(text);
        while (sb.length() < length) {
            sb.append(' ');
        }
        return sb.toString();
    }

    // VIOLATION: unused private method
    private static String padLeft(String text, int length) {
        if (text == null) text = "";
        if (text.length() >= length) return text.substring(0, length);
        StringBuilder sb = new StringBuilder();
        while (sb.length() < length - text.length()) {
            sb.append(' ');
        }
        sb.append(text);
        return sb.toString();
    }

    // VIOLATION: unused private method
    private static String repeatChar(char c, int count) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < count; i++) {
            sb.append(c);
        }
        return sb.toString();
    }

    // VIOLATION: Missing Javadoc
    // VIOLATION: String concatenation in loop instead of StringBuilder
    // VIOLATION: Raw type List (should be List<Student>)
    public static String formatStudentReport(List students) {
        // VIOLATION: String concatenation in loop
        String report = "";
        report = report + "=== Student Enrollment Report ===\n";
        report = report + "Generated: " + LocalDate.now().format(DateTimeFormatter.ISO_DATE) + "\n";
        report = report + "Total Students: " + students.size() + "\n";
        report = report + "================================\n\n";

        for (Object obj : students) {
            Student s = (Student) obj;
            // VIOLATION: String concatenation in loop (should use StringBuilder)
            report = report + "ID: " + s.getId() + "\n";
            report = report + "Name: " + s.getFirstName() + " " + s.getLastName() + "\n";
            report = report + "Email: " + s.getEmail() + "\n";
            report = report + "Department: " + s.getDepartment() + "\n";
            report = report + "GPA: " + s.getGpa() + "\n";
            report = report + "Credits: " + s.getTotalCredits() + "\n";
            report = report + "Status: " + s.getStatus() + "\n";
            if (s.isFinancialHold()) {
                report = report + "*** FINANCIAL HOLD ***\n";
            }
            if (s.isAcademicProbation()) {
                report = report + "*** ACADEMIC PROBATION ***\n";
            }
            report = report + "--------------------------------\n";
        }

        System.out.println("Student report generated for " + students.size() + " students");

        return report;
    }

    // VIOLATION: Missing Javadoc
    // VIOLATION: String concatenation in loop instead of StringBuilder
    // VIOLATION: Raw type List (should be List<Course>)
    public static String formatCourseReport(List courses) {
        // VIOLATION: String concatenation in loop
        String report = "";
        report = report + "=== Course Catalog Report ===\n";
        report = report + "Generated: " + LocalDate.now().format(DateTimeFormatter.ISO_DATE) + "\n";
        report = report + "Total Courses: " + courses.size() + "\n";
        report = report + "=============================\n\n";

        for (Object obj : courses) {
            Course c = (Course) obj;
            // VIOLATION: String concatenation in loop (should use StringBuilder)
            report = report + "Code: " + c.getCourseCode() + "\n";
            report = report + "Title: " + c.getTitle() + "\n";
            report = report + "Department: " + c.getDepartment() + "\n";
            report = report + "Credits: " + c.getCredits() + "\n";
            report = report + "Instructor: " + c.getInstructor() + "\n";
            report = report + "Schedule: " + c.getSchedule() + "\n";
            report = report + "Enrollment: " + c.getCurrentEnrollment() + "/" + c.getMaxCapacity() + "\n";
            if (!c.getPrerequisites().isEmpty()) {
                report = report + "Prerequisites: " + String.join(", ", c.getPrerequisites()) + "\n";
            }
            report = report + "Active: " + (c.isActive() ? "Yes" : "No") + "\n";
            report = report + "------------------------------\n";
        }

        System.out.println("Course report generated for " + courses.size() + " courses");

        return report;
    }

    // VIOLATION: Missing Javadoc
    // VIOLATION: String concatenation in loop instead of StringBuilder
    public static String formatEnrollmentSummary(List<Map<String, Object>> enrollmentData) {
        // VIOLATION: String concatenation in loop
        String report = "";
        report = report + "=== Enrollment Summary ===\n";
        report = report + "Generated: " + LocalDate.now().format(DateTimeFormatter.ISO_DATE) + "\n";
        report = report + "==========================\n\n";

        int totalEnrolled = 0;
        int totalWaitlisted = 0;
        int totalDropped = 0;

        for (Map<String, Object> entry : enrollmentData) {
            String studentName = (String) entry.getOrDefault("studentName", "Unknown");
            String courseCode = (String) entry.getOrDefault("courseCode", "N/A");
            String status = (String) entry.getOrDefault("status", "N/A");

            // VIOLATION: String concatenation in loop (should use StringBuilder)
            report = report + studentName + " | " + courseCode + " | " + status + "\n";

            // VIOLATION: magic number-style string comparison (should use enum)
            if ("ENROLLED".equals(status)) {
                totalEnrolled++;
            } else if ("WAITLISTED".equals(status)) {
                totalWaitlisted++;
            } else if ("DROPPED".equals(status)) {
                totalDropped++;
            }
        }

        report = report + "\n--- Totals ---\n";
        report = report + "Enrolled: " + totalEnrolled + "\n";
        report = report + "Waitlisted: " + totalWaitlisted + "\n";
        report = report + "Dropped: " + totalDropped + "\n";
        report = report + "Total: " + enrollmentData.size() + "\n";

        System.out.println("Enrollment summary generated: " + enrollmentData.size() + " entries");

        return report;
    }

    // VIOLATION: Missing Javadoc
    // VIOLATION: String concatenation in loop instead of StringBuilder
    public static String formatGradeDistribution(Map<String, Integer> distribution) {
        // VIOLATION: String concatenation
        String report = "";
        report = report + "=== Grade Distribution ===\n\n";

        int total = 0;
        for (int count : distribution.values()) {
            total += count;
        }

        for (Map.Entry<String, Integer> entry : distribution.entrySet()) {
            String grade = entry.getKey();
            int count = entry.getValue();
            double percentage = total > 0 ? (double) count / total * 100 : 0; // VIOLATION: magic number 100

            // VIOLATION: String concatenation
            report = report + grade + ": " + count + " students";
            report = report + " (" + String.format("%.1f", percentage) + "%)\n";

            // VIOLATION: String concatenation for bar chart
            String bar = "";
            int barLength = (int) (percentage / 2); // VIOLATION: magic number 2
            for (int i = 0; i < barLength; i++) {
                bar = bar + "#";
            }
            report = report + "  " + bar + "\n";
        }

        report = report + "\nTotal: " + total + " students\n";

        System.out.println("Grade distribution report generated");

        return report;
    }

    // VIOLATION: Missing Javadoc
    // VIOLATION: Multiple return types via Object — poor API design
    public static Object formatQuickStats(int studentCount, int courseCount, double avgGpa) {
        if (studentCount <= 0) {
            return "No data available";
        }

        Map<String, Object> stats = new HashMap<>();
        stats.put("students", studentCount);
        stats.put("courses", courseCount);
        stats.put("averageGpa", Math.round(avgGpa * 100.0) / 100.0); // VIOLATION: magic number
        stats.put("timestamp", LocalDate.now().toString());

        System.out.println("Quick stats: " + studentCount + " students, " + courseCount + " courses");

        return stats;
    }
}
