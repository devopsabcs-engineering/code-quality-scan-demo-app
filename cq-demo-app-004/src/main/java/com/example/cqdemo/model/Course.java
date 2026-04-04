package com.example.cqdemo.model;

import java.util.ArrayList;
import java.util.List;

// VIOLATION: Missing Javadoc on public class
public class Course {

    private Long id;
    private String courseCode;
    private String title;
    private String department;
    private int credits;
    private int maxCapacity;
    private int currentEnrollment;
    private String instructor;
    private String schedule; // e.g., "MWF 10:00-11:00"
    private List<String> prerequisites;
    private String semester;
    private boolean active;

    public Course() {
        this.prerequisites = new ArrayList<>();
        this.active = true;
    }

    public Course(Long id, String courseCode, String title, int credits) {
        this.id = id;
        this.courseCode = courseCode;
        this.title = title;
        this.credits = credits;
        this.maxCapacity = 30; // VIOLATION: magic number
        this.currentEnrollment = 0;
        this.prerequisites = new ArrayList<>();
        this.active = true;
    }

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public String getCourseCode() { return courseCode; }
    public void setCourseCode(String courseCode) { this.courseCode = courseCode; }

    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }

    public String getDepartment() { return department; }
    public void setDepartment(String department) { this.department = department; }

    public int getCredits() { return credits; }
    public void setCredits(int credits) { this.credits = credits; }

    public int getMaxCapacity() { return maxCapacity; }
    public void setMaxCapacity(int maxCapacity) { this.maxCapacity = maxCapacity; }

    public int getCurrentEnrollment() { return currentEnrollment; }
    public void setCurrentEnrollment(int currentEnrollment) { this.currentEnrollment = currentEnrollment; }

    public String getInstructor() { return instructor; }
    public void setInstructor(String instructor) { this.instructor = instructor; }

    public String getSchedule() { return schedule; }
    public void setSchedule(String schedule) { this.schedule = schedule; }

    public List<String> getPrerequisites() { return prerequisites; }
    public void setPrerequisites(List<String> prerequisites) { this.prerequisites = prerequisites; }

    public String getSemester() { return semester; }
    public void setSemester(String semester) { this.semester = semester; }

    public boolean isActive() { return active; }
    public void setActive(boolean active) { this.active = active; }
}
