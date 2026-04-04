# cq-demo-app-004 — Student Enrollment System (Java/Spring Boot)

Code Quality demo application #4 — a **Student Enrollment Management System** built with Java 21 and Spring Boot 3.3. This app contains **intentional code quality violations** for scanner demonstration purposes.

## Intentional Violations

| Category | Count | Location |
|----------|-------|----------|
| **High Cyclomatic Complexity** | 3 | `StudentService.processEnrollment()` (CCN > 15), `StudentController.enrollStudent()`, `GradeCalculator.calculateGPA()` |
| **Code Duplication** | 4 | Validation logic duplicated across `StudentService`, `CourseService`, `StudentController`, `CourseController` |
| **Missing Javadoc** | 10+ | All public classes and methods across `controller/`, `service/`, `util/` |
| **Magic Numbers** | 15+ | `GradeCalculator` (0.3, 0.7, 65, 70, 80, 90, 100), `StudentService` (2.0, 3.5, 21) |
| **Long Methods** | 3 | `processEnrollment()` > 100 lines, `calculateClassStatistics()` > 50 lines, `calculateGPA()` > 50 lines |
| **System.out.println** | 12+ | Throughout all classes instead of SLF4J logger |
| **Raw Types** | 2 | `ReportFormatter.formatStudentReport(List)`, `formatCourseReport(List)` |
| **String Concatenation in Loops** | 5 | `ReportFormatter` — all formatting methods |
| **Unused Private Methods** | 3 | `ReportFormatter.padRight()`, `padLeft()`, `repeatChar()` |
| **Low Test Coverage** | — | Only `contextLoads()` test — < 10% coverage |

## API Endpoints

### Students
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/students` | List all students |
| GET | `/api/students/{id}` | Get student by ID |
| GET | `/api/students/{id}/summary` | Get student summary |
| POST | `/api/students` | Create a student |
| PUT | `/api/students/{id}` | Update a student |
| DELETE | `/api/students/{id}` | Delete a student |
| POST | `/api/students/{id}/enroll` | Enroll in a course |
| GET | `/api/students/{id}/validate` | Validate enrollment eligibility |

### Courses
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/courses` | List all courses |
| GET | `/api/courses/{id}` | Get course by ID |
| GET | `/api/courses/{id}/summary` | Get course summary |
| GET | `/api/courses/code/{code}` | Get course by code |
| POST | `/api/courses` | Create a course |
| PUT | `/api/courses/{id}` | Update a course |
| DELETE | `/api/courses/{id}` | Delete a course |
| POST | `/api/courses/{id}/validate-enrollment` | Validate course enrollment |
| GET | `/api/courses/{id}/schedule-conflicts` | Check schedule conflicts |

### Health
| Method | Path | Description |
|--------|------|-------------|
| GET | `/actuator/health` | Spring Actuator health check |

## Tech Stack

- **Java 21** with **Spring Boot 3.3**
- **Maven** build system
- **JaCoCo** for coverage reporting
- **Checkstyle** (Google style) for lint analysis
- **Docker** multi-stage build

## Run Locally

Build and run with Docker (works in GitHub Codespaces):

```bash
docker build -t cq-demo-app-004 .
docker run -p 8080:8080 cq-demo-app-004
```

Then browse to [http://localhost:8080/api/students](http://localhost:8080/api/students).

### Run with Maven (requires Java 21)

```bash
mvn spring-boot:run
```

### Run Tests

```bash
mvn test
```

### Generate Coverage Report

```bash
mvn test jacoco:report
```

Coverage report is generated at `target/site/jacoco/index.html`.

### Run Checkstyle

```bash
mvn checkstyle:check
```

## Azure Deployment

This app deploys as a Docker container to Azure Web App for Containers via the GitHub Actions workflow in `.github/workflows/deploy.yml`.

Infrastructure is defined in `infra/main.bicep` and provisions:
- Azure Container Registry (ACR) with globally unique name
- App Service Plan (Linux, B1)
- Web App for Containers

All resource names use `uniqueString(resourceGroup().id)` for global uniqueness across multiple workshop participants.
