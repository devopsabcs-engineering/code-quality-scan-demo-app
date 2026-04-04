using Microsoft.AspNetCore.Mvc;
using CqDemoApp003.Models;
using CqDemoApp003.Services;

namespace CqDemoApp003.Controllers;

/// <summary>
/// Controller for employee management operations.
/// INTENTIONAL VIOLATIONS:
/// - SearchEmployees() has high cyclomatic complexity (CCN > 10)
/// - Duplicated response formatting patterns with PayrollController
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class EmployeesController : ControllerBase
{
    private readonly EmployeeService _employeeService;

    public EmployeesController(EmployeeService employeeService)
    {
        _employeeService = employeeService;
    }

    [HttpGet]
    public IActionResult GetAll()
    {
        var employees = _employeeService.GetAll();
        // VIOLATION: Duplicated response wrapping pattern — same as PayrollController
        var response = new
        {
            success = true,
            count = employees.Count,
            data = employees,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "EmployeesController"
        };
        return Ok(response);
    }

    [HttpGet("{id}")]
    public IActionResult GetById(int id)
    {
        var employee = _employeeService.GetById(id);
        if (employee == null)
        {
            // VIOLATION: Duplicated error response pattern — same as PayrollController
            var errorResponse = new
            {
                success = false,
                error = "Not Found",
                message = $"Employee with ID {id} was not found",
                timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
                source = "EmployeesController"
            };
            return NotFound(errorResponse);
        }
        // VIOLATION: Duplicated response wrapping pattern — same as PayrollController
        var response = new
        {
            success = true,
            data = employee,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "EmployeesController"
        };
        return Ok(response);
    }

    [HttpPost]
    public IActionResult Create([FromBody] Employee employee)
    {
        if (!_employeeService.ValidateEmployeeData(
            employee.FirstName, employee.LastName, employee.Email,
            employee.Department, employee.Salary))
        {
            return BadRequest(new { success = false, error = "Invalid employee data" });
        }

        var created = _employeeService.Create(employee);
        // VIOLATION: Duplicated response wrapping pattern
        var response = new
        {
            success = true,
            data = created,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "EmployeesController"
        };
        return CreatedAtAction(nameof(GetById), new { id = created.Id }, response);
    }

    [HttpPost("{id}/action")]
    public IActionResult ProcessAction(int id, [FromBody] EmployeeActionRequest request)
    {
        var result = _employeeService.ProcessEmployeeAction(
            id, request.ActionType, request.TargetDepartment,
            request.NewSalary, request.Reason);

        if (result.Contains("not found", StringComparison.OrdinalIgnoreCase))
            return NotFound(new { success = false, message = result });

        if (result.Contains("cannot", StringComparison.OrdinalIgnoreCase) ||
            result.Contains("required", StringComparison.OrdinalIgnoreCase) ||
            result.Contains("already", StringComparison.OrdinalIgnoreCase) ||
            result.Contains("exceed", StringComparison.OrdinalIgnoreCase))
            return BadRequest(new { success = false, message = result });

        return Ok(new { success = true, message = result });
    }

    // VIOLATION: High Cyclomatic Complexity (CCN > 10) — many filter branches
    [HttpGet("search")]
    public IActionResult SearchEmployees(
        [FromQuery] string? department,
        [FromQuery] string? status,
        [FromQuery] int? minSeniority,
        [FromQuery] int? maxSeniority,
        [FromQuery] decimal? minSalary,
        [FromQuery] decimal? maxSalary,
        [FromQuery] string? sortBy,
        [FromQuery] string? sortOrder)
    {
        var employees = _employeeService.GetAll().AsEnumerable();

        if (!string.IsNullOrEmpty(department))
        {
            employees = employees.Where(e => e.Department.Equals(department, StringComparison.OrdinalIgnoreCase));
        }

        if (!string.IsNullOrEmpty(status))
        {
            employees = employees.Where(e => e.Status.Equals(status, StringComparison.OrdinalIgnoreCase));
        }

        if (minSeniority.HasValue)
        {
            employees = employees.Where(e => e.SeniorityLevel >= minSeniority.Value);
        }

        if (maxSeniority.HasValue)
        {
            employees = employees.Where(e => e.SeniorityLevel <= maxSeniority.Value);
        }

        if (minSalary.HasValue)
        {
            employees = employees.Where(e => e.Salary >= minSalary.Value);
        }

        if (maxSalary.HasValue)
        {
            employees = employees.Where(e => e.Salary <= maxSalary.Value);
        }

        // Sorting logic adds complexity branches
        if (!string.IsNullOrEmpty(sortBy))
        {
            bool descending = sortOrder?.ToLower() == "desc";

            if (sortBy.ToLower() == "name")
            {
                employees = descending
                    ? employees.OrderByDescending(e => e.LastName)
                    : employees.OrderBy(e => e.LastName);
            }
            else if (sortBy.ToLower() == "salary")
            {
                employees = descending
                    ? employees.OrderByDescending(e => e.Salary)
                    : employees.OrderBy(e => e.Salary);
            }
            else if (sortBy.ToLower() == "seniority")
            {
                employees = descending
                    ? employees.OrderByDescending(e => e.SeniorityLevel)
                    : employees.OrderBy(e => e.SeniorityLevel);
            }
            else if (sortBy.ToLower() == "hiredate")
            {
                employees = descending
                    ? employees.OrderByDescending(e => e.HireDate)
                    : employees.OrderBy(e => e.HireDate);
            }
            else if (sortBy.ToLower() == "department")
            {
                employees = descending
                    ? employees.OrderByDescending(e => e.Department)
                    : employees.OrderBy(e => e.Department);
            }
        }

        var result = employees.ToList();

        // VIOLATION: Duplicated response wrapping pattern — same as PayrollController
        var response = new
        {
            success = true,
            count = result.Count,
            data = result,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "EmployeesController"
        };
        return Ok(response);
    }
}
