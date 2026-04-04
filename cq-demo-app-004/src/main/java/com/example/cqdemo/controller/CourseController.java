package com.example.cqdemo.controller;

import com.example.cqdemo.model.Course;
import com.example.cqdemo.service.CourseService;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

// VIOLATION: Missing Javadoc on public class
@RestController
@RequestMapping("/api/courses")
public class CourseController {

    private final CourseService courseService;

    public CourseController(CourseService courseService) {
        this.courseService = courseService;
    }

    @GetMapping
    public List<Course> getAllCourses() {
        return courseService.getAllCourses();
    }

    @GetMapping("/{id}")
    public Object getCourse(@PathVariable Long id) {
        Course course = courseService.getCourseById(id);
        if (course == null) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Course not found");
            return error;
        }
        return course;
    }

    @GetMapping("/{id}/summary")
    public Map<String, Object> getCourseSummary(@PathVariable Long id) {
        return courseService.getCourseSummary(id);
    }

    @GetMapping("/code/{courseCode}")
    public Object getCourseByCode(@PathVariable String courseCode) {
        Course course = courseService.getCourseByCode(courseCode);
        if (course == null) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Course not found: " + courseCode);
            return error;
        }
        return course;
    }

    @PostMapping
    public Course createCourse(@RequestBody Map<String, Object> body) {
        return courseService.createCourse(
                (String) body.get("courseCode"),
                (String) body.get("title"),
                (String) body.getOrDefault("department", "General"),
                body.containsKey("credits") ? ((Number) body.get("credits")).intValue() : 3,
                (String) body.getOrDefault("instructor", "TBD"),
                (String) body.getOrDefault("schedule", "TBD")
        );
    }

    @PutMapping("/{id}")
    public Object updateCourse(@PathVariable Long id, @RequestBody Course course) {
        Course updated = courseService.updateCourse(id, course);
        if (updated == null) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Course not found");
            return error;
        }
        return updated;
    }

    @DeleteMapping("/{id}")
    public Map<String, Object> deleteCourse(@PathVariable Long id) {
        Map<String, Object> result = new HashMap<>();
        boolean deleted = courseService.deleteCourse(id);
        result.put("deleted", deleted);
        result.put("courseId", id);
        return result;
    }

    // =========================================================================
    // VIOLATION: DUPLICATION — Inline validation in the controller mirrors
    // the same patterns from CourseService.validateCourseForEnrollment and
    // StudentController.enrollStudent. Both controllers repeat the same
    // request-parameter extraction and null-check cascades.
    // =========================================================================
    @PostMapping("/{id}/validate-enrollment")
    public Map<String, Object> validateEnrollment(@PathVariable Long id,
                                                    @RequestBody Map<String, Object> body) {
        Map<String, Object> response = new HashMap<>();

        Long studentId = body.containsKey("studentId")
                ? ((Number) body.get("studentId")).longValue()
                : null;

        if (studentId == null) {
            response.put("error", "studentId is required");
            response.put("valid", false);
            return response;
        }

        // DUPLICATION: same completed-courses extraction pattern as controller enrollStudent
        @SuppressWarnings("unchecked")
        List<String> completedCourses = body.containsKey("completedCourses")
                ? (List<String>) body.get("completedCourses")
                : new ArrayList<>();

        Map<String, Object> validationResult = courseService.validateCourseForEnrollment(
                id, studentId, completedCourses);

        response.putAll(validationResult);

        // VIOLATION: System.out.println instead of logger
        System.out.println("Course enrollment validation: courseId=" + id +
                ", studentId=" + studentId +
                ", valid=" + validationResult.get("valid"));

        return response;
    }

    @GetMapping("/{id}/schedule-conflicts")
    public Map<String, Object> checkScheduleConflicts(@PathVariable Long id,
                                                       @RequestParam(required = false) List<Long> enrolledCourseIds) {
        Course course = courseService.getCourseById(id);
        if (course == null) {
            Map<String, Object> error = new HashMap<>();
            error.put("error", "Course not found");
            return error;
        }

        if (enrolledCourseIds == null) {
            enrolledCourseIds = new ArrayList<>();
        }

        return courseService.checkScheduleConflicts(course.getSchedule(), enrolledCourseIds);
    }

    // VIOLATION: Missing Javadoc, System.out usage
    @GetMapping("/health")
    public Map<String, String> health() {
        System.out.println("Health check requested");
        Map<String, String> status = new HashMap<>();
        status.put("status", "UP");
        status.put("service", "CourseController");
        return status;
    }
}
