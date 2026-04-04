package com.example.cqdemo.model;

import java.time.LocalDate;

// VIOLATION: Missing Javadoc on public class
public class Student {

    private Long id;
    private String firstName;
    private String lastName;
    private String email;
    private double gpa;
    private int totalCredits;
    private boolean financialHold;
    private boolean academicProbation;
    private String department;
    private LocalDate enrollmentDate;
    private String status; // ACTIVE, SUSPENDED, GRADUATED, WITHDRAWN

    public Student() {
    }

    // VIOLATION: Missing Javadoc on constructor
    public Student(Long id, String firstName, String lastName, String email) {
        this.id = id;
        this.firstName = firstName;
        this.lastName = lastName;
        this.email = email;
        this.gpa = 0.0;
        this.totalCredits = 0;
        this.financialHold = false;
        this.academicProbation = false;
        this.status = "ACTIVE";
        this.enrollmentDate = LocalDate.now();
    }

    // Getters and setters — no Javadoc (intentional violation)
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }

    public String getFirstName() { return firstName; }
    public void setFirstName(String firstName) { this.firstName = firstName; }

    public String getLastName() { return lastName; }
    public void setLastName(String lastName) { this.lastName = lastName; }

    public String getEmail() { return email; }
    public void setEmail(String email) { this.email = email; }

    public double getGpa() { return gpa; }
    public void setGpa(double gpa) { this.gpa = gpa; }

    public int getTotalCredits() { return totalCredits; }
    public void setTotalCredits(int totalCredits) { this.totalCredits = totalCredits; }

    public boolean isFinancialHold() { return financialHold; }
    public void setFinancialHold(boolean financialHold) { this.financialHold = financialHold; }

    public boolean isAcademicProbation() { return academicProbation; }
    public void setAcademicProbation(boolean academicProbation) { this.academicProbation = academicProbation; }

    public String getDepartment() { return department; }
    public void setDepartment(String department) { this.department = department; }

    public LocalDate getEnrollmentDate() { return enrollmentDate; }
    public void setEnrollmentDate(LocalDate enrollmentDate) { this.enrollmentDate = enrollmentDate; }

    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }
}
