package com.example.cqdemo.service;

import com.example.cqdemo.model.Course;
import com.example.cqdemo.model.Enrollment;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;

// VIOLATION: Missing Javadoc on public class
@Service
public class CourseService {

    private final Map<Long, Course> courses = new HashMap<>();
    private final Map<Long, List<Enrollment>> courseEnrollments = new HashMap<>();
    private final AtomicLong idCounter = new AtomicLong(1);

    public CourseService() {
        // Seed sample courses
        Course cs101 = createCourse("CS101", "Intro to Computer Science", "Computer Science", 3, "Dr. Smith", "MWF 10:00-11:00");
        Course cs201 = createCourse("CS201", "Data Structures", "Computer Science", 4, "Dr. Jones", "TTh 14:00-15:30");
        cs201.getPrerequisites().add("CS101");

        Course math101 = createCourse("MATH101", "Calculus I", "Mathematics", 4, "Dr. Lee", "MWF 09:00-10:00");
        Course math201 = createCourse("MATH201", "Linear Algebra", "Mathematics", 3, "Dr. Kim", "TTh 10:00-11:30");
        math201.getPrerequisites().add("MATH101");

        createCourse("PHYS101", "Physics I", "Physics", 4, "Dr. Newton", "MWF 11:00-12:00");
        createCourse("ENG101", "English Composition", "English", 3, "Dr. Austen", "TTh 09:00-10:30");
    }

    // VIOLATION: Missing Javadoc on public method
    public Course createCourse(String courseCode, String title, String department,
                                int credits, String instructor, String schedule) {
        Long id = idCounter.getAndIncrement();
        Course course = new Course(id, courseCode, title, credits);
        course.setDepartment(department);
        course.setInstructor(instructor);
        course.setSchedule(schedule);
        course.setSemester("Fall 2024");
        courses.put(id, course);
        courseEnrollments.put(id, new ArrayList<>());
        return course;
    }

    public List<Course> getAllCourses() {
        return new ArrayList<>(courses.values());
    }

    public Course getCourseById(Long id) {
        return courses.get(id);
    }

    public Course getCourseByCode(String courseCode) {
        return courses.values().stream()
                .filter(c -> c.getCourseCode().equals(courseCode))
                .findFirst()
                .orElse(null);
    }

    // =========================================================================
    // VIOLATION: DUPLICATED VALIDATION LOGIC
    // This method contains the same validation pattern as
    // StudentService.validateStudentForEnrollment — checking capacity, schedule,
    // prerequisites, and department rules. The structure and conditional flow
    // are nearly identical, just operating on Course instead of Student.
    // =========================================================================
    public Map<String, Object> validateCourseForEnrollment(Long courseId, Long studentId,
                                                             List<String> completedCourses) {
        Map<String, Object> validationResult = new HashMap<>();
        validationResult.put("valid", false);
        validationResult.put("courseId", courseId);

        Course course = courses.get(courseId);
        if (course == null) {
            validationResult.put("error", "Course not found");
            return validationResult;
        }

        if (!course.isActive()) {
            validationResult.put("error", "Course is not active");
            return validationResult;
        }

        // DUPLICATION: Capacity check — same pattern as processEnrollment
        if (course.getCurrentEnrollment() >= course.getMaxCapacity()) {
            validationResult.put("error", "Course is at full capacity");
            return validationResult;
        }

        // DUPLICATION: Prerequisite check — same pattern as processEnrollment
        if (course.getPrerequisites() != null && !course.getPrerequisites().isEmpty()) {
            for (String prereq : course.getPrerequisites()) {
                boolean completed = false;
                for (String completedCourse : completedCourses) {
                    if (completedCourse != null && completedCourse.equals(prereq)) {
                        completed = true;
                        break;
                    }
                }
                if (!completed) {
                    validationResult.put("error", "Missing prerequisite: " + prereq);
                    return validationResult;
                }
            }
        }

        // DUPLICATION: Schedule conflict check — same pattern as processEnrollment
        List<Enrollment> existingEnrollments = courseEnrollments.getOrDefault(courseId, new ArrayList<>());
        for (Enrollment e : existingEnrollments) {
            if (e.getStudentId().equals(studentId) && "ENROLLED".equals(e.getStatus())) {
                validationResult.put("error", "Student is already enrolled in this course");
                return validationResult;
            }
        }

        validationResult.put("valid", true);
        validationResult.put("courseCode", course.getCourseCode());
        validationResult.put("title", course.getTitle());
        validationResult.put("availableSeats", course.getMaxCapacity() - course.getCurrentEnrollment());
        validationResult.put("credits", course.getCredits());

        return validationResult;
    }

    // VIOLATION: DUPLICATED data lookup pattern (same pattern as StudentService.getStudentSummary)
    public Map<String, Object> getCourseSummary(Long courseId) {
        Map<String, Object> summary = new HashMap<>();
        Course course = courses.get(courseId);

        if (course == null) {
            summary.put("error", "Course not found");
            return summary;
        }

        summary.put("id", course.getId());
        summary.put("courseCode", course.getCourseCode());
        summary.put("title", course.getTitle());
        summary.put("department", course.getDepartment());
        summary.put("credits", course.getCredits());
        summary.put("instructor", course.getInstructor());
        summary.put("schedule", course.getSchedule());
        summary.put("semester", course.getSemester());
        summary.put("maxCapacity", course.getMaxCapacity());
        summary.put("currentEnrollment", course.getCurrentEnrollment());
        summary.put("availableSeats", course.getMaxCapacity() - course.getCurrentEnrollment());
        summary.put("active", course.isActive());
        summary.put("prerequisites", course.getPrerequisites());

        List<Enrollment> enrollments = courseEnrollments.getOrDefault(courseId, new ArrayList<>());
        summary.put("enrolledCount", enrollments.stream()
                .filter(e -> "ENROLLED".equals(e.getStatus()))
                .count());
        summary.put("waitlistedCount", enrollments.stream()
                .filter(e -> "WAITLISTED".equals(e.getStatus()))
                .count());
        summary.put("completedCount", enrollments.stream()
                .filter(e -> "COMPLETED".equals(e.getStatus()))
                .count());

        return summary;
    }

    // =========================================================================
    // VIOLATION: DUPLICATED scheduling logic
    // This method mirrors the schedule-conflict detection in StudentService,
    // checking the same fields and using the same iteration pattern.
    // =========================================================================
    public Map<String, Object> checkScheduleConflicts(String schedule, List<Long> enrolledCourseIds) {
        Map<String, Object> conflicts = new HashMap<>();
        conflicts.put("hasConflicts", false);
        List<Map<String, String>> conflictList = new ArrayList<>();

        if (schedule == null || schedule.isEmpty()) {
            conflicts.put("error", "Schedule must be specified");
            return conflicts;
        }

        for (Long courseId : enrolledCourseIds) {
            Course existingCourse = courses.get(courseId);
            if (existingCourse == null) {
                continue;
            }
            // Simplified conflict detection — exact match
            if (existingCourse.getSchedule() != null && existingCourse.getSchedule().equals(schedule)) {
                Map<String, String> conflict = new HashMap<>();
                conflict.put("courseCode", existingCourse.getCourseCode());
                conflict.put("title", existingCourse.getTitle());
                conflict.put("existingSchedule", existingCourse.getSchedule());
                conflict.put("requestedSchedule", schedule);
                conflictList.add(conflict);
            }
        }

        if (!conflictList.isEmpty()) {
            conflicts.put("hasConflicts", true);
            conflicts.put("conflicts", conflictList);
        }

        return conflicts;
    }

    public boolean deleteCourse(Long id) {
        return courses.remove(id) != null;
    }

    public Course updateCourse(Long id, Course updated) {
        Course existing = courses.get(id);
        if (existing == null) {
            return null;
        }
        if (updated.getTitle() != null) existing.setTitle(updated.getTitle());
        if (updated.getInstructor() != null) existing.setInstructor(updated.getInstructor());
        if (updated.getSchedule() != null) existing.setSchedule(updated.getSchedule());
        if (updated.getMaxCapacity() > 0) existing.setMaxCapacity(updated.getMaxCapacity());
        return existing;
    }
}
