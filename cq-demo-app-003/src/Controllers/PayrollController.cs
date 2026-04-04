using Microsoft.AspNetCore.Mvc;
using CqDemoApp003.Models;
using CqDemoApp003.Services;

namespace CqDemoApp003.Controllers;

/// <summary>
/// Controller for payroll operations.
/// INTENTIONAL VIOLATIONS:
/// - Duplicated response formatting patterns from EmployeesController
/// - Duplicated error handling patterns
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class PayrollController : ControllerBase
{
    private readonly PayrollService _payrollService;

    public PayrollController(PayrollService payrollService)
    {
        _payrollService = payrollService;
    }

    [HttpGet]
    public IActionResult GetAll()
    {
        var records = _payrollService.GetAll();
        // VIOLATION: Duplicated response wrapping pattern — same as EmployeesController
        var response = new
        {
            success = true,
            count = records.Count,
            data = records,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "PayrollController"
        };
        return Ok(response);
    }

    [HttpGet("{employeeId}")]
    public IActionResult GetByEmployeeId(int employeeId)
    {
        var record = _payrollService.GetByEmployeeId(employeeId);
        if (record == null)
        {
            // VIOLATION: Duplicated error response pattern — same as EmployeesController
            var errorResponse = new
            {
                success = false,
                error = "Not Found",
                message = $"Payroll record for employee ID {employeeId} was not found",
                timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
                source = "PayrollController"
            };
            return NotFound(errorResponse);
        }
        // VIOLATION: Duplicated response wrapping pattern — same as EmployeesController
        var response = new
        {
            success = true,
            data = record,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "PayrollController"
        };
        return Ok(response);
    }

    [HttpPost("calculate")]
    public IActionResult Calculate([FromBody] PayrollCalculationRequest request)
    {
        if (request.EmployeeId <= 0)
        {
            return BadRequest(new { success = false, error = "Invalid employee ID" });
        }

        var record = _payrollService.CalculatePayroll(
            request.EmployeeId,
            request.PayPeriod,
            request.BonusOverride,
            request.DeductionOverride);

        if (record == null)
        {
            // VIOLATION: Duplicated error response pattern — same as EmployeesController
            var errorResponse = new
            {
                success = false,
                error = "Not Found",
                message = $"Employee with ID {request.EmployeeId} not found or not active",
                timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
                source = "PayrollController"
            };
            return NotFound(errorResponse);
        }

        // VIOLATION: Duplicated response wrapping pattern — same as EmployeesController
        var response = new
        {
            success = true,
            data = record,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "PayrollController"
        };
        return Ok(response);
    }

    [HttpGet("summary/{payPeriod}")]
    public IActionResult GetSummary(string payPeriod)
    {
        var summary = _payrollService.GeneratePayrollSummary(payPeriod);
        // VIOLATION: Duplicated response wrapping pattern — same as EmployeesController
        var response = new
        {
            success = true,
            data = summary,
            timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
            source = "PayrollController"
        };
        return Ok(response);
    }
}
