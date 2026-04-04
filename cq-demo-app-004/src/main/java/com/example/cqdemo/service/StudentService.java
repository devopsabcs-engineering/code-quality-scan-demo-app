package com.example.cqdemo.service;

import com.example.cqdemo.model.Course;
import com.example.cqdemo.model.Enrollment;
import com.example.cqdemo.model.Student;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;

// VIOLATION: Missing Javadoc on public class
@Service
public class StudentService {

    private final Map<Long, Student> students = new HashMap<>();
    private final Map<Long, List<Enrollment>> studentEnrollments = new HashMap<>();
    private final AtomicLong idCounter = new AtomicLong(1);

    public StudentService() {
        // Seed sample data
        createStudent("Alice", "Johnson", "alice.johnson@university.edu", "Computer Science", 3.8, 45);
        createStudent("Bob", "Smith", "bob.smith@university.edu", "Mathematics", 2.1, 30);
        createStudent("Carol", "Williams", "carol.williams@university.edu", "Physics", 3.2, 60);
        createStudent("Dave", "Brown", "dave.brown@university.edu", "Engineering", 1.8, 15);
        createStudent("Eve", "Davis", "eve.davis@university.edu", "Computer Science", 3.95, 90);
    }

    // VIOLATION: Missing Javadoc on public method
    public Student createStudent(String firstName, String lastName, String email,
                                  String department, double gpa, int credits) {
        Long id = idCounter.getAndIncrement();
        Student student = new Student(id, firstName, lastName, email);
        student.setDepartment(department);
        student.setGpa(gpa);
        student.setTotalCredits(credits);
        students.put(id, student);
        studentEnrollments.put(id, new ArrayList<>());
        return student;
    }

    public List<Student> getAllStudents() {
        return new ArrayList<>(students.values());
    }

    public Student getStudentById(Long id) {
        return students.get(id);
    }

    // =========================================================================
    // VIOLATION: HIGH CYCLOMATIC COMPLEXITY (CCN > 15)
    // This method has deeply nested conditionals covering prerequisites, credit
    // limits, schedule conflicts, GPA requirements, financial holds, department
    // approval, academic probation, waitlist logic, and semester validation —
    // all packed into a single monolithic method.
    // =========================================================================
    public Map<String, Object> processEnrollment(Long studentId, Course course,
                                                   boolean departmentApproved,
                                                   boolean isOverrideAllowed,
                                                   String semester) {
        Map<String, Object> result = new HashMap<>();
        result.put("success", false);
        result.put("studentId", studentId);
        result.put("courseCode", course.getCourseCode());

        Student student = students.get(studentId);

        // Branch 1: Student existence check
        if (student == null) {
            result.put("error", "Student not found");
            return result;
        }

        // Branch 2: Student status check
        if (student.getStatus() == null || !student.getStatus().equals("ACTIVE")) {
            result.put("error", "Student is not active");
            return result;
        }

        // Branch 3: Financial hold check
        if (student.isFinancialHold()) {
            if (isOverrideAllowed) {
                result.put("warning", "Financial hold overridden by administrator");
                System.out.println("WARNING: Financial hold override for student " + studentId); // VIOLATION: sysout
            } else {
                result.put("error", "Student has a financial hold — cannot enroll");
                return result;
            }
        }

        // Branch 4: Academic probation check
        if (student.isAcademicProbation()) {
            if (student.getGpa() < 2.0) { // VIOLATION: magic number
                if (!departmentApproved) {
                    result.put("error", "Student on academic probation with GPA below 2.0 requires department approval");
                    return result;
                } else {
                    result.put("warning", "Academic probation enrollment approved by department");
                }
            } else {
                result.put("info", "Student on academic probation but GPA is satisfactory");
            }
        }

        // Branch 5: Credit limit check
        int newTotalCredits = student.getTotalCredits() + course.getCredits();
        if (newTotalCredits > 21) { // VIOLATION: magic number
            if (student.getGpa() >= 3.5 && isOverrideAllowed) { // VIOLATION: magic number
                result.put("warning", "Credit limit exceeded but GPA qualifies for overload");
                System.out.println("Credit overload approved for student " + studentId); // VIOLATION: sysout
            } else if (student.getGpa() >= 3.5 && departmentApproved) { // VIOLATION: magic number
                result.put("warning", "Credit overload approved by department for high-GPA student");
            } else {
                result.put("error", "Exceeds maximum credit limit of 21 credits");
                return result;
            }
        }

        // Branch 6: Prerequisite check — DUPLICATED pattern (also in CourseService)
        List<Enrollment> enrollments = studentEnrollments.getOrDefault(studentId, new ArrayList<>());
        if (course.getPrerequisites() != null && !course.getPrerequisites().isEmpty()) {
            for (String prereq : course.getPrerequisites()) {
                boolean completed = false;
                for (Enrollment e : enrollments) {
                    if (e.getStatus() != null && e.getStatus().equals("COMPLETED")) {
                        // Note: simplified check — in reality would match against course code
                        completed = true;
                        break;
                    }
                }
                if (!completed) {
                    if (isOverrideAllowed && departmentApproved) {
                        result.put("warning", "Prerequisite " + prereq + " waived by department");
                    } else {
                        result.put("error", "Missing prerequisite: " + prereq);
                        return result;
                    }
                }
            }
        }

        // Branch 7: Schedule conflict check — DUPLICATED pattern (also in CourseService)
        if (course.getSchedule() != null) {
            for (Enrollment e : enrollments) {
                if (e.getStatus() != null && e.getStatus().equals("ENROLLED")) {
                    // Simplified schedule conflict detection
                    String existingSchedule = "MWF 10:00-11:00"; // placeholder
                    if (existingSchedule.equals(course.getSchedule())) {
                        if (isOverrideAllowed) {
                            result.put("warning", "Schedule conflict overridden");
                        } else {
                            result.put("error", "Schedule conflict detected with existing enrollment");
                            return result;
                        }
                    }
                }
            }
        }

        // Branch 8: Course capacity check
        if (course.getCurrentEnrollment() >= course.getMaxCapacity()) {
            if (isOverrideAllowed) {
                result.put("warning", "Course is full — override enrollment approved");
                System.out.println("Over-capacity enrollment for course " + course.getCourseCode()); // VIOLATION: sysout
            } else {
                // Waitlist logic
                Enrollment waitlistEntry = new Enrollment(
                        (long) enrollments.size() + 1000,
                        studentId,
                        course.getId()
                );
                waitlistEntry.setStatus("WAITLISTED");
                enrollments.add(waitlistEntry);
                result.put("waitlisted", true);
                result.put("error", "Course full — student added to waitlist");
                return result;
            }
        }

        // Branch 9: Semester validation
        if (semester == null || semester.isEmpty()) {
            result.put("error", "Semester must be specified");
            return result;
        } else {
            if (!semester.matches("(Fall|Spring|Summer)\\s\\d{4}")) {
                result.put("error", "Invalid semester format — expected 'Fall 2024', 'Spring 2025', etc.");
                return result;
            }
        }

        // Branch 10: Department matching check
        if (course.getDepartment() != null && student.getDepartment() != null) {
            if (!course.getDepartment().equals(student.getDepartment())) {
                if (!departmentApproved && course.getCredits() > 3) { // VIOLATION: magic number
                    result.put("error", "Cross-department enrollment for courses > 3 credits requires approval");
                    return result;
                }
            }
        }

        // All checks passed — create enrollment
        Enrollment enrollment = new Enrollment(
                (long) enrollments.size() + 1,
                studentId,
                course.getId()
        );
        enrollment.setStatus("ENROLLED");
        enrollments.add(enrollment);

        course.setCurrentEnrollment(course.getCurrentEnrollment() + 1);
        student.setTotalCredits(newTotalCredits);

        result.put("success", true);
        result.put("enrollmentId", enrollment.getId());
        result.put("message", "Successfully enrolled in " + course.getCourseCode());

        // VIOLATION: System.out.println instead of logger
        System.out.println("Enrollment completed: student=" + studentId +
                ", course=" + course.getCourseCode() +
                ", semester=" + semester);

        return result;
    }

    // =========================================================================
    // VIOLATION: DUPLICATED VALIDATION LOGIC
    // The validateStudentForEnrollment method repeats patterns found in
    // processEnrollment and also in CourseService.validateCourseForEnrollment.
    // =========================================================================
    public Map<String, Object> validateStudentForEnrollment(Long studentId, int additionalCredits) {
        Map<String, Object> validationResult = new HashMap<>();
        validationResult.put("valid", false);
        validationResult.put("studentId", studentId);

        Student student = students.get(studentId);
        if (student == null) {
            validationResult.put("error", "Student not found");
            return validationResult;
        }

        if (student.getStatus() == null || !student.getStatus().equals("ACTIVE")) {
            validationResult.put("error", "Student is not active");
            return validationResult;
        }

        if (student.isFinancialHold()) {
            validationResult.put("error", "Student has a financial hold");
            return validationResult;
        }

        if (student.isAcademicProbation() && student.getGpa() < 2.0) {
            validationResult.put("error", "Student on probation with low GPA");
            return validationResult;
        }

        int projectedCredits = student.getTotalCredits() + additionalCredits;
        if (projectedCredits > 21) {
            if (student.getGpa() < 3.5) {
                validationResult.put("error", "Credit limit exceeded");
                return validationResult;
            } else {
                validationResult.put("warning", "Credit overload — high GPA student");
            }
        }

        validationResult.put("valid", true);
        validationResult.put("studentName", student.getFirstName() + " " + student.getLastName());
        validationResult.put("currentCredits", student.getTotalCredits());
        validationResult.put("projectedCredits", projectedCredits);

        return validationResult;
    }

    // VIOLATION: DUPLICATED data lookup pattern (same pattern in CourseService)
    public Map<String, Object> getStudentSummary(Long studentId) {
        Map<String, Object> summary = new HashMap<>();
        Student student = students.get(studentId);

        if (student == null) {
            summary.put("error", "Student not found");
            return summary;
        }

        summary.put("id", student.getId());
        summary.put("name", student.getFirstName() + " " + student.getLastName());
        summary.put("email", student.getEmail());
        summary.put("department", student.getDepartment());
        summary.put("gpa", student.getGpa());
        summary.put("totalCredits", student.getTotalCredits());
        summary.put("status", student.getStatus());
        summary.put("financialHold", student.isFinancialHold());
        summary.put("academicProbation", student.isAcademicProbation());
        summary.put("enrollmentDate", student.getEnrollmentDate());

        List<Enrollment> enrollments = studentEnrollments.getOrDefault(studentId, new ArrayList<>());
        summary.put("activeEnrollments", enrollments.stream()
                .filter(e -> "ENROLLED".equals(e.getStatus()))
                .count());
        summary.put("completedCourses", enrollments.stream()
                .filter(e -> "COMPLETED".equals(e.getStatus()))
                .count());

        return summary;
    }

    public boolean deleteStudent(Long id) {
        return students.remove(id) != null;
    }

    public Student updateStudent(Long id, Student updated) {
        Student existing = students.get(id);
        if (existing == null) {
            return null;
        }
        if (updated.getFirstName() != null) existing.setFirstName(updated.getFirstName());
        if (updated.getLastName() != null) existing.setLastName(updated.getLastName());
        if (updated.getEmail() != null) existing.setEmail(updated.getEmail());
        if (updated.getDepartment() != null) existing.setDepartment(updated.getDepartment());
        return existing;
    }
}
