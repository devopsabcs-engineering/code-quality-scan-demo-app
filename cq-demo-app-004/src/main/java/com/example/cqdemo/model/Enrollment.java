package com.example.cqdemo.model;

import java.time.LocalDateTime;

// VIOLATION: Missing Javadoc on public class
public class Enrollment {

    private Long id;
    private Long studentId;
    private Long courseId;
    private String status; // ENROLLED, WAITLISTED, DROPPED, COMPLETED
    private String grade;
    private LocalDateTime enrolledAt;
    private LocalDateTime completedAt;

    public Enrollment() {
    }

    public Enrollment(Long id, Long studentId, Long courseId) {
        this.id = id;
        this.studentId = studentId;
        this.courseId = courseId;
        this.status = "ENROLLED";
        this.enrolledAt = LocalDateTime.now();
    }

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public Long getStudentId() { return studentId; }
    public void setStudentId(Long studentId) { this.studentId = studentId; }

    public Long getCourseId() { return courseId; }
    public void setCourseId(Long courseId) { this.courseId = courseId; }

    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }

    public String getGrade() { return grade; }
    public void setGrade(String grade) { this.grade = grade; }

    public LocalDateTime getEnrolledAt() { return enrolledAt; }
    public void setEnrolledAt(LocalDateTime enrolledAt) { this.enrolledAt = enrolledAt; }

    public LocalDateTime getCompletedAt() { return completedAt; }
    public void setCompletedAt(LocalDateTime completedAt) { this.completedAt = completedAt; }
}
