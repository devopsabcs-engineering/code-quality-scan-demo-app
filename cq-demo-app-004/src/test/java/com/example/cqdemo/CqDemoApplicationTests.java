package com.example.cqdemo;

import org.springframework.boot.test.context.SpringBootTest;
import org.junit.jupiter.api.Test;

// =========================================================================
// VIOLATION: LOW TEST COVERAGE (< 50%)
// Only one test — verifies the Spring context loads. All service logic,
// controller endpoints, GradeCalculator utility methods, and
// ReportFormatter methods are COMPLETELY UNTESTED.
//
// Missing test coverage for:
//   - StudentService.processEnrollment() — 0% coverage
//   - StudentService.validateStudentForEnrollment() — 0% coverage
//   - StudentService.getStudentSummary() — 0% coverage
//   - CourseService.validateCourseForEnrollment() — 0% coverage
//   - CourseService.getCourseSummary() — 0% coverage
//   - CourseService.checkScheduleConflicts() — 0% coverage
//   - StudentController (all endpoints) — 0% coverage
//   - CourseController (all endpoints) — 0% coverage
//   - GradeCalculator.calculateWeightedGrade() — 0% coverage
//   - GradeCalculator.getLetterGrade() — 0% coverage
//   - GradeCalculator.calculateClassStatistics() — 0% coverage
//   - GradeCalculator.calculateGPA() — 0% coverage
//   - ReportFormatter (all methods) — 0% coverage
// =========================================================================
@SpringBootTest
class CqDemoApplicationTests {

    @Test
    void contextLoads() {
        // This is the ONLY test in the entire project.
        // It merely verifies that the Spring application context starts
        // without errors. No business logic is tested.
    }
}
