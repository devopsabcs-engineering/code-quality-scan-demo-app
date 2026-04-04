package com.example.cqdemo.controller;

import com.example.cqdemo.model.Course;
import com.example.cqdemo.model.Student;
import com.example.cqdemo.service.CourseService;
import com.example.cqdemo.service.StudentService;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

// VIOLATION: Missing Javadoc on public class
@RestController
@RequestMapping("/api/students")
public class StudentController {

    private final StudentService studentService;
    private final CourseService courseService;

    public StudentController(StudentService studentService, CourseService courseService) {
        this.studentService = studentService;
        this.courseService = courseService;
    }

    @GetMapping
    public List<Student> getAllStudents() {
        return studentService.getAllStudents();
    }

    @GetMapping("/{id}")
    public Object getStudent(@PathVariable Long id) {
        Student student = studentService.getStudentById(id);
        if (student == null) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Student not found");
            return error;
        }
        return student;
    }

    @GetMapping("/{id}/summary")
    public Map<String, Object> getStudentSummary(@PathVariable Long id) {
        return studentService.getStudentSummary(id);
    }

    @PostMapping
    public Student createStudent(@RequestBody Map<String, Object> body) {
        return studentService.createStudent(
                (String) body.get("firstName"),
                (String) body.get("lastName"),
                (String) body.get("email"),
                (String) body.getOrDefault("department", "Undeclared"),
                body.containsKey("gpa") ? ((Number) body.get("gpa")).doubleValue() : 0.0,
                body.containsKey("credits") ? ((Number) body.get("credits")).intValue() : 0
        );
    }

    @PutMapping("/{id}")
    public Object updateStudent(@PathVariable Long id, @RequestBody Student student) {
        Student updated = studentService.updateStudent(id, student);
        if (updated == null) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Student not found");
            return error;
        }
        return updated;
    }

    @DeleteMapping("/{id}")
    public Map<String, Object> deleteStudent(@PathVariable Long id) {
        Map<String, Object> result = new HashMap<>();
        boolean deleted = studentService.deleteStudent(id);
        result.put("deleted", deleted);
        result.put("studentId", id);
        return result;
    }

    // =========================================================================
    // VIOLATION: HIGH CYCLOMATIC COMPLEXITY
    // The enrollment endpoint performs multiple validation branches inline in the
    // controller rather than delegating cleanly to the service layer. This
    // increases the CCN of the method beyond acceptable thresholds.
    // =========================================================================
    @PostMapping("/{id}/enroll")
    public Map<String, Object> enrollStudent(@PathVariable Long id,
                                              @RequestBody Map<String, Object> body) {
        Map<String, Object> response = new HashMap<>();

        // Extract and validate request parameters with complex branching
        String courseCode = (String) body.get("courseCode");
        String semester = (String) body.get("semester");
        boolean departmentApproved = body.containsKey("departmentApproved")
                && (Boolean) body.get("departmentApproved");
        boolean override = body.containsKey("override")
                && (Boolean) body.get("override");

        if (courseCode == null || courseCode.isEmpty()) {
            response.put("error", "courseCode is required");
            response.put("success", false);
            return response;
        }

        if (semester == null || semester.isEmpty()) {
            response.put("error", "semester is required");
            response.put("success", false);
            return response;
        }

        // Validate student first
        Student student = studentService.getStudentById(id);
        if (student == null) {
            response.put("error", "Student not found with id: " + id);
            response.put("success", false);
            return response;
        }

        // Look up course
        Course course = courseService.getCourseByCode(courseCode);
        if (course == null) {
            response.put("error", "Course not found: " + courseCode);
            response.put("success", false);
            return response;
        }

        // Pre-validate student eligibility
        Map<String, Object> studentValidation = studentService.validateStudentForEnrollment(
                id, course.getCredits());
        if (!(Boolean) studentValidation.get("valid")) {
            if (!override) {
                response.put("error", studentValidation.get("error"));
                response.put("success", false);
                response.put("validation", studentValidation);
                return response;
            } else {
                response.put("warning", "Student validation failed but override is active");
            }
        }

        // Process the enrollment
        Map<String, Object> enrollmentResult = studentService.processEnrollment(
                id, course, departmentApproved, override, semester);

        response.putAll(enrollmentResult);

        // VIOLATION: System.out.println instead of logger
        if ((Boolean) enrollmentResult.getOrDefault("success", false)) {
            System.out.println("REST: Enrollment successful for student " + id +
                    " in course " + courseCode);
        } else {
            System.out.println("REST: Enrollment failed for student " + id +
                    " in course " + courseCode + ": " + enrollmentResult.get("error"));
        }

        return response;
    }

    @GetMapping("/{id}/validate")
    public Map<String, Object> validateStudent(@PathVariable Long id,
                                                @RequestParam(defaultValue = "3") int credits) {
        return studentService.validateStudentForEnrollment(id, credits);
    }

    // VIOLATION: Missing Javadoc, System.out usage
    @GetMapping("/health")
    public Map<String, String> health() {
        System.out.println("Health check requested");
        Map<String, String> status = new HashMap<>();
        status.put("status", "UP");
        status.put("service", "StudentController");
        return status;
    }
}
